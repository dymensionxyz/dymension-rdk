package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Event Hooks
// These can be utilized to communicate between a staking keeper and another
// keeper which must take particular actions when governors/delegators change
// state. The second keeper must implement this interface, which then the
// staking keeper can call.

// StakingHooks event hooks for staking governor object (noalias)
type StakingHooks interface {
	AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error   // Must be called when a validator is created
	BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) error // Must be called when a validator's state changes

	AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error        // Must be called when a validator is deleted
	AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error         // Must be called when a validator is bonded
	AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error // Must be called when a validator begins unbonding

	BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error        // Must be called when a delegation is created
	BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error // Must be called when a delegation's shares are modified
	BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error        // Must be called when a delegation is removed
	AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error
}

// combine multiple staking hooks, all hook functions are run in array sequence
var _ StakingHooks = &MultiStakingHooks{}

// var _ stakingtypes.StakingHooks = &MultiStakingHooks{}

type MultiStakingHooks []StakingHooks

func NewMultiStakingHooks(hooks ...StakingHooks) MultiStakingHooks {
	return hooks
}

func (h MultiStakingHooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].AfterValidatorCreated(ctx, valAddr); err != nil {
			return err
		}
	}

	return nil
}

func (h MultiStakingHooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].BeforeValidatorModified(ctx, valAddr); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiStakingHooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].AfterValidatorRemoved(ctx, consAddr, valAddr); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiStakingHooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].AfterValidatorBonded(ctx, consAddr, valAddr); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiStakingHooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].AfterValidatorBeginUnbonding(ctx, consAddr, valAddr); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiStakingHooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].BeforeDelegationCreated(ctx, delAddr, valAddr); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiStakingHooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].BeforeDelegationSharesModified(ctx, delAddr, valAddr); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiStakingHooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].BeforeDelegationRemoved(ctx, delAddr, valAddr); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiStakingHooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].AfterDelegationModified(ctx, delAddr, valAddr); err != nil {
			return err
		}
	}
	return nil
}
