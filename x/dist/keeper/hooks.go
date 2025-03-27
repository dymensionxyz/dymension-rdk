package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/utils/erc20"
)

// Wrapper struct
// It holds the vanilla implementation of the hooks and our overrides
type Hooks struct {
	distkeeper.Hooks
	DistKeeper Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Create new distribution hooks
func (k Keeper) Hooks() Hooks {
	return Hooks{
		Hooks:      k.Keeper.Hooks(),
		DistKeeper: k,
	}
}

// AfterValidatorRemoved performs clean up after a validator is removed
// It first calls the base implementation, then it checks if the bond denom is
// registered as an ERC20 token, and if so, converts the balance of the
// delegator's withdraw address from the coin to the ERC20 token.
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	withdrawAddr := h.DistKeeper.GetDelegatorWithdrawAddr(ctx, sdk.AccAddress(valAddr))

	err := h.Hooks.AfterValidatorRemoved(ctx, consAddr, valAddr)
	if err != nil {
		return err
	}

	err = erc20.ConvertAllBalances(ctx, h.DistKeeper.erc20k, h.DistKeeper.bankKeeper, withdrawAddr)
	if err != nil {
		return err
	}

	return nil
}

// BeforeDelegationSharesModified is called before modifying the delegation shares of a delegator.
// It first calls the base implementation, then it checks if the bond denom is
// registered as an ERC20 token, and if so, converts the balance of the
// delegator's withdraw address from the coin to the ERC20 token.
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	withdrawAddr := h.DistKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)

	err := h.Hooks.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	if err != nil {
		return err
	}

	err = erc20.ConvertAllBalances(ctx, h.DistKeeper.erc20k, h.DistKeeper.bankKeeper, withdrawAddr)
	if err != nil {
		return err
	}

	return nil
}
