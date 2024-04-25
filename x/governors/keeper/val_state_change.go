package keeper

import (
	"bytes"
	"fmt"
	"sort"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

// BlockGovernorUpdates calculates the GovernorUpdates for the current block
// Called in each EndBlock
func (k Keeper) BlockGovernorUpdates(ctx sdk.Context) {
	// Calculate governor set changes.

	err := k.ApplyGovernorSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	// unbond all mature governors from the unbonding queue
	k.UnbondAllMatureGovernors(ctx)

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		addr, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress := sdk.MustAccAddressFromBech32(dvPair.DelegatorAddress)

		balances, err := k.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyGovernor, dvPair.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	for _, dvvTriplet := range matureRedelegations {
		valSrcAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorSrcAddress)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress := sdk.MustAccAddressFromBech32(dvvTriplet.DelegatorAddress)

		balances, err := k.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeySrcGovernor, dvvTriplet.ValidatorSrcAddress),
				sdk.NewAttribute(types.AttributeKeyDstGovernor, dvvTriplet.ValidatorDstAddress),
			),
		)
	}

	return
}

// ApplyGovernorSetUpdates applies accumulated updates to the bonded governor set. Also,
// * Updates the active valset as keyed by LastGovernorPowerKey.
// * Updates the total power as keyed by LastTotalPowerKey.
// * Updates governor status' according to updated powers.
// * Updates the fee pool bonded vs not-bonded tokens.
// * Updates relevant indices.
// It gets called once after genesis, another time maybe after genesis transactions,
// then once at every EndBlock.
//
// CONTRACT: Only governors with non-zero power or zero-power that were bonded
// at the previous block height or were removed from the governor set entirely
// are returned to Tendermint.
func (k Keeper) ApplyGovernorSetUpdates(ctx sdk.Context) (err error) {

	// FIXME: keep updates
	params := k.GetParams(ctx)
	maxGovernors := params.MaxValidators
	powerReduction := k.PowerReduction(ctx)
	totalPower := sdk.ZeroInt()
	amtFromBondedToNotBonded, amtFromNotBondedToBonded := sdk.ZeroInt(), sdk.ZeroInt()

	// Retrieve the last governor set.
	// The persistent set is updated later in this function.
	// (see LastGovernorPowerKey).
	last, err := k.getLastGovernorsByAddr(ctx)
	if err != nil {
		return err
	}

	// Iterate over governors, highest power to lowest.
	iterator := k.GovernorsPowerStoreIterator(ctx)
	defer iterator.Close()

	for count := 0; iterator.Valid() && count < int(maxGovernors); iterator.Next() {
		// everything that is iterated in this loop is becoming or already a
		// part of the bonded governor set
		valAddr := sdk.ValAddress(iterator.Value())
		governor := k.mustGetGovernor(ctx, valAddr)

		// if we get to a zero-power governor (which we don't bond),
		// there are no more possible bonded governors
		if governor.PotentialConsensusPower(k.PowerReduction(ctx)) == 0 {
			break
		}

		// apply the appropriate state change if necessary
		switch {
		case governor.IsUnbonded():
			governor, err = k.unbondedToBonded(ctx, governor)
			if err != nil {
				return
			}
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(governor.GetTokens())
		case governor.IsUnbonding():
			governor, err = k.unbondingToBonded(ctx, governor)
			if err != nil {
				return
			}
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(governor.GetTokens())
		case governor.IsBonded():
			// no state change
		default:
			panic("unexpected governor status")
		}

		// fetch the old power bytes
		valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
		if err != nil {
			return err
		}
		oldPowerBytes, found := last[valAddrStr]
		newPower := governor.ConsensusPower(powerReduction)
		newPowerBytes := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: newPower})

		// update the governor set if power has changed
		if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
			k.SetLastGovernorPower(ctx, valAddr, newPower)
		}

		delete(last, valAddrStr)
		count++

		totalPower = totalPower.Add(sdk.NewInt(newPower))
	}

	noLongerBonded, err := sortNoLongerBonded(last)
	if err != nil {
		return err
	}

	for _, valAddrBytes := range noLongerBonded {
		governor := k.mustGetGovernor(ctx, sdk.ValAddress(valAddrBytes))
		governor, err = k.bondedToUnbonding(ctx, governor)
		if err != nil {
			return
		}
		amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(governor.GetTokens())
		k.DeleteLastGovernorPower(ctx, governor.GetOperator())
	}

	// Update the pools based on the recent updates in the governor set:
	// - The tokens from the non-bonded candidates that enter the new governor set need to be transferred
	// to the Bonded pool.
	// - The tokens from the bonded governors that are being kicked out from the governor set
	// need to be transferred to the NotBonded pool.
	switch {
	// Compare and subtract the respective amounts to only perform one transfer.
	// This is done in order to avoid doing multiple updates inside each iterator/loop.
	case amtFromNotBondedToBonded.GT(amtFromBondedToNotBonded):
		k.notBondedTokensToBonded(ctx, amtFromNotBondedToBonded.Sub(amtFromBondedToNotBonded))
	case amtFromNotBondedToBonded.LT(amtFromBondedToNotBonded):
		k.bondedTokensToNotBonded(ctx, amtFromBondedToNotBonded.Sub(amtFromNotBondedToBonded))
	default: // equal amounts of tokens; no update required
	}

	// set total power on lookup index if there are any updates
	//FIXME: check only when changed
	// if len(updates) > 0 {
	k.SetLastTotalPower(ctx, totalPower)
	// }

	return nil
}

// Governor state transitions

func (k Keeper) bondedToUnbonding(ctx sdk.Context, governor types.Governor) (types.Governor, error) {
	if !governor.IsBonded() {
		panic(fmt.Sprintf("bad state transition bondedToUnbonding, governor: %v\n", governor))
	}

	return k.beginUnbondingGovernor(ctx, governor)
}

func (k Keeper) unbondingToBonded(ctx sdk.Context, governor types.Governor) (types.Governor, error) {
	if !governor.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, governor: %v\n", governor))
	}

	return k.bondGovernor(ctx, governor)
}

func (k Keeper) unbondedToBonded(ctx sdk.Context, governor types.Governor) (types.Governor, error) {
	if !governor.IsUnbonded() {
		panic(fmt.Sprintf("bad state transition unbondedToBonded, governor: %v\n", governor))
	}

	return k.bondGovernor(ctx, governor)
}

// UnbondingToUnbonded switches a governor from unbonding state to unbonded state
func (k Keeper) UnbondingToUnbonded(ctx sdk.Context, governor types.Governor) types.Governor {
	if !governor.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, governor: %v\n", governor))
	}

	return k.completeUnbondingGovernor(ctx, governor)
}

// perform all the store operations for when a governor status becomes bonded
func (k Keeper) bondGovernor(ctx sdk.Context, governor types.Governor) (types.Governor, error) {
	// delete the governor by power index, as the key will change
	k.DeleteGovernorByPowerIndex(ctx, governor)

	governor = governor.UpdateStatus(types.Bonded)

	// save the now bonded governor record to the two referenced stores
	k.SetGovernor(ctx, governor)
	k.SetGovernorByPowerIndex(ctx, governor)

	// delete from queue if present
	k.DeleteGovernorQueue(ctx, governor)

	// trigger hook
	err := k.AfterGovernorBonded(ctx, governor.GetOperator())
	if err != nil {
		return governor, err
	}

	return governor, nil
}

// perform all the store operations for when a governor begins unbonding
func (k Keeper) beginUnbondingGovernor(ctx sdk.Context, governor types.Governor) (types.Governor, error) {
	params := k.GetParams(ctx)

	// delete the governor by power index, as the key will change
	k.DeleteGovernorByPowerIndex(ctx, governor)

	// sanity check
	if governor.Status != types.Bonded {
		panic(fmt.Sprintf("should not already be unbonded or unbonding, governor: %v\n", governor))
	}

	governor = governor.UpdateStatus(types.Unbonding)

	// set the unbonding completion time and completion height appropriately
	governor.UnbondingTime = ctx.BlockHeader().Time.Add(params.UnbondingTime)
	governor.UnbondingHeight = ctx.BlockHeader().Height

	// save the now unbonded governor record and power index
	k.SetGovernor(ctx, governor)
	k.SetGovernorByPowerIndex(ctx, governor)

	// Adds to unbonding governor queue
	k.InsertUnbondingGovernorQueue(ctx, governor)

	// trigger hook
	k.AfterGovernorBeginUnbonding(ctx, governor.GetOperator())

	return governor, nil
}

// perform all the store operations for when a governor status becomes unbonded
func (k Keeper) completeUnbondingGovernor(ctx sdk.Context, governor types.Governor) types.Governor {
	governor = governor.UpdateStatus(types.Unbonded)
	k.SetGovernor(ctx, governor)

	return governor
}

// map of operator bech32-addresses to serialized power
// We use bech32 strings here, because we can't have slices as keys: map[[]byte][]byte
type governorsByAddr map[string][]byte

// get the last governor set
func (k Keeper) getLastGovernorsByAddr(ctx sdk.Context) (governorsByAddr, error) {
	last := make(governorsByAddr)

	iterator := k.LastGovernorsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// extract the governor address from the key (prefix is 1-byte, addrLen is 1-byte)
		valAddr := types.AddressFromLastGovernorPowerKey(iterator.Key())
		valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)
		if err != nil {
			return nil, err
		}

		powerBytes := iterator.Value()
		last[valAddrStr] = make([]byte, len(powerBytes))
		copy(last[valAddrStr], powerBytes)
	}

	return last, nil
}

// given a map of remaining governors to previous bonded power
// returns the list of governors to be unbonded, sorted by operator address
func sortNoLongerBonded(last governorsByAddr) ([][]byte, error) {
	// sort the map keys for determinism
	noLongerBonded := make([][]byte, len(last))
	index := 0

	for valAddrStr := range last {
		valAddrBytes, err := sdk.ValAddressFromBech32(valAddrStr)
		if err != nil {
			return nil, err
		}
		noLongerBonded[index] = valAddrBytes
		index++
	}
	// sorted by address - order doesn't matter
	sort.SliceStable(noLongerBonded, func(i, j int) bool {
		// -1 means strictly less than
		return bytes.Compare(noLongerBonded[i], noLongerBonded[j]) == -1
	})

	return noLongerBonded, nil
}
