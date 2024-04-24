package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

// InitGenesis sets the pool and parameters for the provided keeper.  For each
// governor in data, it sets that governor in the keeper along with manually
// setting the indexes. In addition, it also sets any delegations found in
// data. Finally, it updates the bonded governors.
// Returns final governor set after applying all declaration and delegations
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) []abci.ValidatorUpdate {
	bondedTokens := sdk.ZeroInt()
	notBondedTokens := sdk.ZeroInt()

	k.SetParams(ctx, data.Params)
	k.SetLastTotalPower(ctx, data.LastTotalPower)

	for _, governor := range data.Governors {
		k.SetGovernor(ctx, governor)

		// Manually set indices for the first time
		k.SetGovernorByPowerIndex(ctx, governor)

		// Call the creation hook if not exported
		if !data.Exported {
			if err := k.AfterGovernorCreated(ctx, governor.GetOperator()); err != nil {
				panic(err)
			}
		}

		// update timeslice if necessary
		if governor.IsUnbonding() {
			k.InsertUnbondingGovernorQueue(ctx, governor)
		}

		switch governor.GetStatus() {
		case types.Bonded:
			bondedTokens = bondedTokens.Add(governor.GetTokens())

		case types.Unbonding, types.Unbonded:
			notBondedTokens = notBondedTokens.Add(governor.GetTokens())

		default:
			panic("invalid governor status")
		}
	}

	for _, delegation := range data.Delegations {
		delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)

		// Call the before-creation hook if not exported
		if !data.Exported {
			if err := k.BeforeDelegationCreated(ctx, delegatorAddress, delegation.GetValidatorAddr()); err != nil {
				panic(err)
			}
		}

		k.SetDelegation(ctx, delegation)

		// Call the after-modification hook if not exported
		if !data.Exported {
			if err := k.AfterDelegationModified(ctx, delegatorAddress, delegation.GetValidatorAddr()); err != nil {
				panic(err)
			}
		}
	}

	for _, ubd := range data.UnbondingDelegations {
		k.SetUnbondingDelegation(ctx, ubd)

		for _, entry := range ubd.Entries {
			k.InsertUBDQueue(ctx, ubd, entry.CompletionTime)
			notBondedTokens = notBondedTokens.Add(entry.Balance)
		}
	}

	for _, red := range data.Redelegations {
		k.SetRedelegation(ctx, red)

		for _, entry := range red.Entries {
			k.InsertRedelegationQueue(ctx, red, entry.CompletionTime)
		}
	}

	bondedCoins := sdk.NewCoins(sdk.NewCoin(data.Params.BondDenom, bondedTokens))
	notBondedCoins := sdk.NewCoins(sdk.NewCoin(data.Params.BondDenom, notBondedTokens))

	// check if the unbonded and bonded pools accounts exists
	bondedPool := k.GetBondedPool(ctx)
	if bondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	// TODO: remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	bondedBalance := k.bankKeeper.GetAllBalances(ctx, bondedPool.GetAddress())
	if bondedBalance.IsZero() {
		k.authKeeper.SetModuleAccount(ctx, bondedPool)
	}

	// if balance is different from bonded coins panic because genesis is most likely malformed
	if !bondedBalance.IsEqual(bondedCoins) {
		panic(fmt.Sprintf("bonded pool balance is different from bonded coins: %s <-> %s", bondedBalance, bondedCoins))
	}

	notBondedPool := k.GetNotBondedPool(ctx)
	if notBondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	notBondedBalance := k.bankKeeper.GetAllBalances(ctx, notBondedPool.GetAddress())
	if notBondedBalance.IsZero() {
		k.authKeeper.SetModuleAccount(ctx, notBondedPool)
	}

	// If balance is different from non bonded coins panic because genesis is most
	// likely malformed.
	if !notBondedBalance.IsEqual(notBondedCoins) {
		panic(fmt.Sprintf("not bonded pool balance is different from not bonded coins: %s <-> %s", notBondedBalance, notBondedCoins))
	}

	err := k.ApplyGovernorSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, governors, and bonds found in
// the keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	var unbondingDelegations []stakingtypes.UnbondingDelegation

	k.IterateUnbondingDelegations(ctx, func(_ int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
		unbondingDelegations = append(unbondingDelegations, ubd)
		return false
	})

	var redelegations []stakingtypes.Redelegation

	k.IterateRedelegations(ctx, func(_ int64, red stakingtypes.Redelegation) (stop bool) {
		redelegations = append(redelegations, red)
		return false
	})

	return &types.GenesisState{
		Params:               k.GetParams(ctx),
		LastTotalPower:       k.GetLastTotalPower(ctx),
		Governors:            k.GetAllGovernors(ctx),
		Delegations:          k.GetAllDelegations(ctx),
		UnbondingDelegations: unbondingDelegations,
		Redelegations:        redelegations,
		Exported:             true,
	}
}
