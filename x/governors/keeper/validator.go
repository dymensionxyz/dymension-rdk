package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

// get a single governor
func (k Keeper) GetGovernor(ctx sdk.Context, addr sdk.ValAddress) (governor types.Governor, found bool) {
	store := ctx.KVStore(k.storeKey)

	value := store.Get(types.GetGovernorKey(addr))
	if value == nil {
		return governor, false
	}

	governor = types.MustUnmarshalGovernor(k.cdc, value)
	return governor, true
}

func (k Keeper) mustGetGovernor(ctx sdk.Context, addr sdk.ValAddress) types.Governor {
	governor, found := k.GetGovernor(ctx, addr)
	if !found {
		panic(fmt.Sprintf("governor record not found for address: %X\n", addr))
	}

	return governor
}

// set the main record holding governor details
func (k Keeper) SetGovernor(ctx sdk.Context, governor types.Governor) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalGovernor(k.cdc, &governor)
	store.Set(types.GetGovernorKey(governor.GetOperator()), bz)
}

// governor index
func (k Keeper) SetGovernorByPowerIndex(ctx sdk.Context, governor types.Governor) {
	// jailed governors are not kept in the power index
	if governor.Jailed {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetGovernorsByPowerIndexKey(governor, k.PowerReduction(ctx)), governor.GetOperator())
}

// governor index
func (k Keeper) DeleteGovernorByPowerIndex(ctx sdk.Context, governor types.Governor) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetGovernorsByPowerIndexKey(governor, k.PowerReduction(ctx)))
}

// governor index
func (k Keeper) SetNewGovernorByPowerIndex(ctx sdk.Context, governor types.Governor) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetGovernorsByPowerIndexKey(governor, k.PowerReduction(ctx)), governor.GetOperator())
}

// Update the tokens of an existing governor, update the governors power index key
func (k Keeper) AddGovernorTokensAndShares(ctx sdk.Context, governor types.Governor,
	tokensToAdd math.Int,
) (valOut types.Governor, addedShares sdk.Dec) {
	k.DeleteGovernorByPowerIndex(ctx, governor)
	governor, addedShares = governor.AddTokensFromDel(tokensToAdd)
	k.SetGovernor(ctx, governor)
	k.SetGovernorByPowerIndex(ctx, governor)

	return governor, addedShares
}

// Update the tokens of an existing governor, update the governors power index key
func (k Keeper) RemoveGovernorTokensAndShares(ctx sdk.Context, governor types.Governor,
	sharesToRemove sdk.Dec,
) (valOut types.Governor, removedTokens math.Int) {
	k.DeleteGovernorByPowerIndex(ctx, governor)
	governor, removedTokens = governor.RemoveDelShares(sharesToRemove)
	k.SetGovernor(ctx, governor)
	k.SetGovernorByPowerIndex(ctx, governor)

	return governor, removedTokens
}

// Update the tokens of an existing governor, update the governors power index key
func (k Keeper) RemoveGovernorTokens(ctx sdk.Context,
	governor types.Governor, tokensToRemove math.Int,
) types.Governor {
	k.DeleteGovernorByPowerIndex(ctx, governor)
	governor = governor.RemoveTokens(tokensToRemove)
	k.SetGovernor(ctx, governor)
	k.SetGovernorByPowerIndex(ctx, governor)

	return governor
}

// UpdateGovernorCommission attempts to update a governor's commission rate.
// An error is returned if the new commission rate is invalid.
func (k Keeper) UpdateGovernorCommission(ctx sdk.Context,
	governor types.Governor, newRate sdk.Dec,
) (types.Commission, error) {
	commission := governor.Commission
	blockTime := ctx.BlockHeader().Time

	if err := commission.ValidateNewRate(newRate, blockTime); err != nil {
		return commission, err
	}

	if newRate.LT(k.MinCommissionRate(ctx)) {
		return commission, fmt.Errorf("cannot set governor commission to less than minimum rate of %s", k.MinCommissionRate(ctx))
	}

	commission.Rate = newRate
	commission.UpdateTime = blockTime

	return commission, nil
}

// remove the governor record and associated indexes
// except for the bonded governor index which is only handled in ApplyAndReturnTendermintUpdates
// TODO, this function panics, and it's not good.
func (k Keeper) RemoveGovernor(ctx sdk.Context, address sdk.ValAddress) {
	// first retrieve the old governor record
	governor, found := k.GetGovernor(ctx, address)
	if !found {
		return
	}

	if !governor.IsUnbonded() {
		panic("cannot call RemoveGovernor on bonded or unbonding governors")
	}

	if governor.Tokens.IsPositive() {
		panic("attempting to remove a governor which still contains tokens")
	}

	// delete the old governor record
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetGovernorKey(address))
	store.Delete(types.GetGovernorsByPowerIndexKey(governor, k.PowerReduction(ctx)))

	// call hooks
	k.AfterGovernorRemoved(ctx, governor.GetOperator())
}

// get groups of governors

// get the set of all governors with no limits, used during genesis dump
func (k Keeper) GetAllGovernors(ctx sdk.Context) (governors []types.Governor) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GovernorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		governor := types.MustUnmarshalGovernor(k.cdc, iterator.Value())
		governors = append(governors, governor)
	}

	return governors
}

// return a given amount of all the governors
func (k Keeper) GetGovernors(ctx sdk.Context, maxRetrieve uint32) (governors []types.Governor) {
	store := ctx.KVStore(k.storeKey)
	governors = make([]types.Governor, maxRetrieve)

	iterator := sdk.KVStorePrefixIterator(store, types.GovernorsKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		governor := types.MustUnmarshalGovernor(k.cdc, iterator.Value())
		governors[i] = governor
		i++
	}

	return governors[:i] // trim if the array length < maxRetrieve
}

// get the current group of bonded governors sorted by power-rank
func (k Keeper) GetBondedGovernorsByPower(ctx sdk.Context) []types.Governor {
	maxGovernors := k.MaxGovernors(ctx)
	governors := make([]types.Governor, maxGovernors)

	iterator := k.GovernorsPowerStoreIterator(ctx)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxGovernors); iterator.Next() {
		address := iterator.Value()
		governor := k.mustGetGovernor(ctx, address)

		if governor.IsBonded() {
			governors[i] = governor
			i++
		}
	}

	return governors[:i] // trim
}

// returns an iterator for the current governor power store
func (k Keeper) GovernorsPowerStoreIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.GovernorsByPowerIndexKey)
}

// Last Governor Index

// Load the last governor power.
// Returns zero if the operator was not a governor last block.
func (k Keeper) GetLastGovernorPower(ctx sdk.Context, operator sdk.ValAddress) (power int64) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetLastGovernorPowerKey(operator))
	if bz == nil {
		return 0
	}

	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshal(bz, &intV)

	return intV.GetValue()
}

// Set the last governor power.
func (k Keeper) SetLastGovernorPower(ctx sdk.Context, operator sdk.ValAddress, power int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: power})
	store.Set(types.GetLastGovernorPowerKey(operator), bz)
}

// Delete the last governor power.
func (k Keeper) DeleteLastGovernorPower(ctx sdk.Context, operator sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLastGovernorPowerKey(operator))
}

// returns an iterator for the consensus governors in the last block
func (k Keeper) LastGovernorsIterator(ctx sdk.Context) (iterator sdk.Iterator) {
	store := ctx.KVStore(k.storeKey)
	iterator = sdk.KVStorePrefixIterator(store, types.LastGovernorPowerKey)

	return iterator
}

// Iterate over last governor powers.
func (k Keeper) IterateLastGovernorPowers(ctx sdk.Context, handler func(operator sdk.ValAddress, power int64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.LastGovernorPowerKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(types.AddressFromLastGovernorPowerKey(iter.Key()))
		intV := &gogotypes.Int64Value{}

		k.cdc.MustUnmarshal(iter.Value(), intV)

		if handler(addr, intV.GetValue()) {
			break
		}
	}
}

// get the group of the bonded governors
func (k Keeper) GetLastGovernors(ctx sdk.Context) (governors []types.Governor) {
	store := ctx.KVStore(k.storeKey)

	// add the actual governor power sorted store
	maxGovernors := k.MaxGovernors(ctx)
	governors = make([]types.Governor, maxGovernors)

	iterator := sdk.KVStorePrefixIterator(store, types.LastGovernorPowerKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {
		// sanity check
		if i >= int(maxGovernors) {
			panic("more governors than maxGovernors found")
		}

		address := types.AddressFromLastGovernorPowerKey(iterator.Key())
		governor := k.mustGetGovernor(ctx, address)

		governors[i] = governor
		i++
	}

	return governors[:i] // trim
}

// GetUnbondingGovernors returns a slice of mature governor addresses that
// complete their unbonding at a given time and height.
func (k Keeper) GetUnbondingGovernors(ctx sdk.Context, endTime time.Time, endHeight int64) []string {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetGovernorQueueKey(endTime, endHeight))
	if bz == nil {
		return []string{}
	}

	addrs := stakingtypes.ValAddresses{}
	k.cdc.MustUnmarshal(bz, &addrs)

	return addrs.Addresses
}

// SetUnbondingGovernorsQueue sets a given slice of governor addresses into
// the unbonding governor queue by a given height and time.
func (k Keeper) SetUnbondingGovernorsQueue(ctx sdk.Context, endTime time.Time, endHeight int64, addrs []string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stakingtypes.ValAddresses{Addresses: addrs})
	store.Set(types.GetGovernorQueueKey(endTime, endHeight), bz)
}

// InsertUnbondingGovernorQueue inserts a given unbonding governor address into
// the unbonding governor queue for a given height and time.
func (k Keeper) InsertUnbondingGovernorQueue(ctx sdk.Context, val types.Governor) {
	addrs := k.GetUnbondingGovernors(ctx, val.UnbondingTime, val.UnbondingHeight)
	addrs = append(addrs, val.OperatorAddress)
	k.SetUnbondingGovernorsQueue(ctx, val.UnbondingTime, val.UnbondingHeight, addrs)
}

// DeleteGovernorQueueTimeSlice deletes all entries in the queue indexed by a
// given height and time.
func (k Keeper) DeleteGovernorQueueTimeSlice(ctx sdk.Context, endTime time.Time, endHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetGovernorQueueKey(endTime, endHeight))
}

// DeleteGovernorQueue removes a governor by address from the unbonding queue
// indexed by a given height and time.
func (k Keeper) DeleteGovernorQueue(ctx sdk.Context, val types.Governor) {
	addrs := k.GetUnbondingGovernors(ctx, val.UnbondingTime, val.UnbondingHeight)
	newAddrs := []string{}

	for _, addr := range addrs {
		if addr != val.OperatorAddress {
			newAddrs = append(newAddrs, addr)
		}
	}

	if len(newAddrs) == 0 {
		k.DeleteGovernorQueueTimeSlice(ctx, val.UnbondingTime, val.UnbondingHeight)
	} else {
		k.SetUnbondingGovernorsQueue(ctx, val.UnbondingTime, val.UnbondingHeight, newAddrs)
	}
}

// GovernorQueueIterator returns an interator ranging over governors that are
// unbonding whose unbonding completion occurs at the given height and time.
func (k Keeper) GovernorQueueIterator(ctx sdk.Context, endTime time.Time, endHeight int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.GovernorQueueKey, sdk.InclusiveEndBytes(types.GetGovernorQueueKey(endTime, endHeight)))
}

// UnbondAllMatureGovernors unbonds all the mature unbonding governors that
// have finished their unbonding period.
func (k Keeper) UnbondAllMatureGovernors(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	blockTime := ctx.BlockTime()
	blockHeight := ctx.BlockHeight()

	// unbondingValIterator will contains all governor addresses indexed under
	// the GovernorQueueKey prefix. Note, the entire index key is composed as
	// GovernorQueueKey | timeBzLen (8-byte big endian) | timeBz | heightBz (8-byte big endian),
	// so it may be possible that certain governor addresses that are iterated
	// over are not ready to unbond, so an explicit check is required.
	unbondingValIterator := k.GovernorQueueIterator(ctx, blockTime, blockHeight)
	defer unbondingValIterator.Close()

	for ; unbondingValIterator.Valid(); unbondingValIterator.Next() {
		key := unbondingValIterator.Key()
		keyTime, keyHeight, err := types.ParseGovernorQueueKey(key)
		if err != nil {
			panic(fmt.Errorf("failed to parse unbonding key: %w", err))
		}

		// All addresses for the given key have the same unbonding height and time.
		// We only unbond if the height and time are less than the current height
		// and time.
		if keyHeight <= blockHeight && (keyTime.Before(blockTime) || keyTime.Equal(blockTime)) {
			addrs := stakingtypes.ValAddresses{}
			k.cdc.MustUnmarshal(unbondingValIterator.Value(), &addrs)

			for _, valAddr := range addrs.Addresses {
				addr, err := sdk.ValAddressFromBech32(valAddr)
				if err != nil {
					panic(err)
				}
				val, found := k.GetGovernor(ctx, addr)
				if !found {
					panic("governor in the unbonding queue was not found")
				}

				if !val.IsUnbonding() {
					panic("unexpected governor in unbonding queue; status was not unbonding")
				}

				val = k.UnbondingToUnbonded(ctx, val)
				if val.GetDelegatorShares().IsZero() {
					k.RemoveGovernor(ctx, val.GetOperator())
				}
			}

			store.Delete(key)
		}
	}
}

// send a governor to jail
func (k Keeper) jailGovernor(ctx sdk.Context, governor types.Governor) {
	if governor.Jailed {
		panic(fmt.Sprintf("cannot jail already jailed governor, governor: %v\n", governor))
	}

	governor.Jailed = true
	k.SetGovernor(ctx, governor)
	k.DeleteGovernorByPowerIndex(ctx, governor)
}
