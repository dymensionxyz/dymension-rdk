package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// combine multiple staking hooks, all hook functions are run in array sequence
var _ StakingHooks = &MultiStakingHooks{}

type MultiStakingHooks []StakingHooks

func NewMultiStakingHooks(hooks ...StakingHooks) MultiStakingHooks {
	return hooks
}

func (h MultiStakingHooks) AfterGovernorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].AfterGovernorCreated(ctx, valAddr); err != nil {
			return err
		}
	}

	return nil
}

func (h MultiStakingHooks) BeforeGovernorModified(ctx sdk.Context, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].BeforeGovernorModified(ctx, valAddr); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiStakingHooks) AfterGovernorRemoved(ctx sdk.Context, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].AfterGovernorRemoved(ctx, valAddr); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiStakingHooks) AfterGovernorBonded(ctx sdk.Context, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].AfterGovernorBonded(ctx, valAddr); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiStakingHooks) AfterGovernorBeginUnbonding(ctx sdk.Context, valAddr sdk.ValAddress) error {
	for i := range h {
		if err := h[i].AfterGovernorBeginUnbonding(ctx, valAddr); err != nil {
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
