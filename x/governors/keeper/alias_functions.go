package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

/*

The following functions are aliases for the staking module functions that are required for the governance module to function.
These functions are required to implement the StakingComptability interface.
*/

// IterateValidators implements types.StakingComptability.
func (k Keeper) IterateValidators(ctx sdk.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GovernorsKey)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		governor := types.MustUnmarshalGovernor(k.cdc, iterator.Value())
		validator := governor.ToValidator()
		stop := fn(i, validator)

		if stop {
			break
		}
		i++
	}
}

// IterateBondedValidatorsByPower implements types.StakingKeeper.
func (k Keeper) IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	maxGovernors := k.MaxGovernors(ctx)

	iterator := sdk.KVStoreReversePrefixIterator(store, types.GovernorsByPowerIndexKey)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid() && i < int64(maxGovernors); iterator.Next() {
		address := iterator.Value()
		governor := k.mustGetGovernor(ctx, address)

		if governor.IsBonded() {
			validator := governor.ToValidator()
			stop := fn(i, validator)
			if stop {
				break
			}
			i++
		}
	}
}

// Validator implements types.StakingComptability.
func (k Keeper) Validator(ctx sdk.Context, address sdk.ValAddress) stakingtypes.ValidatorI {
	governor := k.Governor(ctx, address)
	if governor == nil {
		return nil
	}
	return governor.ToValidator()
}

// ValidatorByConsAddr implements types.StakingComptability.
func (k Keeper) ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI {
	return nil
}

// Governor Set

// iterate through the governor set and perform the provided function
func (k Keeper) IterateGovernors(ctx sdk.Context, fn func(index int64, governor types.GovernorI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GovernorsKey)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		governor := types.MustUnmarshalGovernor(k.cdc, iterator.Value())
		stop := fn(i, governor) // XXX is this safe will the governor unexposed fields be able to get written to?

		if stop {
			break
		}
		i++
	}
}

// iterate through the bonded governor set and perform the provided function
func (k Keeper) IterateBondedGovernorsByPower(ctx sdk.Context, fn func(index int64, governor types.GovernorI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	maxGovernors := k.MaxGovernors(ctx)

	iterator := sdk.KVStoreReversePrefixIterator(store, types.GovernorsByPowerIndexKey)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid() && i < int64(maxGovernors); iterator.Next() {
		address := iterator.Value()
		governor := k.mustGetGovernor(ctx, address)

		if governor.IsBonded() {
			stop := fn(i, governor) // XXX is this safe will the governor unexposed fields be able to get written to?
			if stop {
				break
			}
			i++
		}
	}
}

// iterate through the active governor set and perform the provided function
func (k Keeper) IterateLastGovernors(ctx sdk.Context, fn func(index int64, governor types.GovernorI) (stop bool)) {
	iterator := k.LastGovernorsIterator(ctx)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		address := types.AddressFromLastGovernorPowerKey(iterator.Key())

		governor, found := k.GetGovernor(ctx, address)
		if !found {
			panic(fmt.Sprintf("governor record not found for address: %v\n", address))
		}

		stop := fn(i, governor) // XXX is this safe will the governor unexposed fields be able to get written to?
		if stop {
			break
		}
		i++
	}
}

// Governor gets the Governor interface for a particular address
func (k Keeper) Governor(ctx sdk.Context, address sdk.ValAddress) types.GovernorI {
	val, found := k.GetGovernor(ctx, address)
	if !found {
		return nil
	}

	return val
}

// Delegation Set

// Delegation get the delegation interface for a particular set of delegator and governor addresses
func (k Keeper) Delegation(ctx sdk.Context, addrDel sdk.AccAddress, addrVal sdk.ValAddress) stakingtypes.DelegationI {
	bond, ok := k.GetDelegation(ctx, addrDel, addrVal)
	if !ok {
		return nil
	}

	return bond
}

// iterate through all of the delegations from a delegator
func (k Keeper) IterateDelegations(ctx sdk.Context, delAddr sdk.AccAddress,
	fn func(index int64, del stakingtypes.DelegationI) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delAddr)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		del := stakingtypes.MustUnmarshalDelegation(k.cdc, iterator.Value())

		stop := fn(i, del)
		if stop {
			break
		}
		i++
	}
}

// return all delegations used during genesis dump
// TODO: remove this func, change all usage for iterate functionality
func (k Keeper) GetAllSDKDelegations(ctx sdk.Context) (delegations []stakingtypes.Delegation) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.DelegationKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := stakingtypes.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
	}

	return
}
