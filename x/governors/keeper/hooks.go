package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AfterGovernorCreated - call hook if registered
func (k Keeper) AfterGovernorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.AfterValidatorCreated(ctx, valAddr)
	}
	return nil
}

// BeforeGovernorModified - call hook if registered
func (k Keeper) BeforeGovernorModified(ctx sdk.Context, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.BeforeValidatorModified(ctx, valAddr)
	}
	return nil
}

// AfterGovernorRemoved - call hook if registered
func (k Keeper) AfterGovernorRemoved(ctx sdk.Context, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.AfterValidatorRemoved(ctx, nil, valAddr)
	}
	return nil
}

// AfterGovernorBonded - call hook if registered
func (k Keeper) AfterGovernorBonded(ctx sdk.Context, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.AfterValidatorBonded(ctx, nil, valAddr)
	}
	return nil
}

// AfterGovernorBeginUnbonding - call hook if registered
func (k Keeper) AfterGovernorBeginUnbonding(ctx sdk.Context, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.AfterValidatorBeginUnbonding(ctx, nil, valAddr)
	}
	return nil
}

// BeforeDelegationCreated - call hook if registered
func (k Keeper) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.BeforeDelegationCreated(ctx, delAddr, valAddr)
	}
	return nil
}

// BeforeDelegationSharesModified - call hook if registered
func (k Keeper) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	}
	return nil
}

// BeforeDelegationRemoved - call hook if registered
func (k Keeper) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.BeforeDelegationRemoved(ctx, delAddr, valAddr)
	}
	return nil
}

// AfterDelegationModified - call hook if registered
func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	if k.hooks != nil {
		return k.hooks.AfterDelegationModified(ctx, delAddr, valAddr)
	}
	return nil
}
