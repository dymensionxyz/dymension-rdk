package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

// Return all governors that a delegator is bonded to. If maxRetrieve is supplied, the respective amount will be returned.
func (k Keeper) GetDelegatorGovernors(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, maxRetrieve uint32,
) types.Governors {
	governors := make([]types.Governor, maxRetrieve)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delegatorAddr)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())

		governor, found := k.GetGovernor(ctx, delegation.GetGovernorAddr())
		if !found {
			panic(types.ErrNoGovernorFound)
		}

		governors[i] = governor
		i++
	}

	return governors[:i] // trim
}

// return a governor that a delegator is bonded to
func (k Keeper) GetDelegatorGovernor(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, governorAddr sdk.ValAddress,
) (governor types.Governor, err error) {
	delegation, found := k.GetDelegation(ctx, delegatorAddr, governorAddr)
	if !found {
		return governor, types.ErrNoDelegation
	}

	governor, found = k.GetGovernor(ctx, delegation.GetGovernorAddr())
	if !found {
		panic(types.ErrNoGovernorFound)
	}

	return governor, nil
}

// return all delegations for a delegator
func (k Keeper) GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []types.Delegation {
	delegations := make([]types.Delegation, 0)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	i := 0

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
		i++
	}

	return delegations
}

// return all unbonding-delegations for a delegator
func (k Keeper) GetAllUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress) []types.UnbondingDelegation {
	unbondingDelegations := make([]types.UnbondingDelegation, 0)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetUBDsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	for i := 0; iterator.Valid(); iterator.Next() {
		unbondingDelegation := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
		i++
	}

	return unbondingDelegations
}

// return all redelegations for a delegator
func (k Keeper) GetAllRedelegations(
	ctx sdk.Context, delegator sdk.AccAddress, srcValAddress, dstValAddress sdk.ValAddress,
) []types.Redelegation {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetREDsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	srcValFilter := !(srcValAddress.Empty())
	dstValFilter := !(dstValAddress.Empty())

	redelegations := []types.Redelegation{}

	for ; iterator.Valid(); iterator.Next() {
		redelegation := types.MustUnmarshalRED(k.cdc, iterator.Value())
		valSrcAddr, err := sdk.ValAddressFromBech32(redelegation.ValidatorSrcAddress)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := sdk.ValAddressFromBech32(redelegation.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}
		if srcValFilter && !(srcValAddress.Equals(valSrcAddr)) {
			continue
		}

		if dstValFilter && !(dstValAddress.Equals(valDstAddr)) {
			continue
		}

		redelegations = append(redelegations, redelegation)
	}

	return redelegations
}
