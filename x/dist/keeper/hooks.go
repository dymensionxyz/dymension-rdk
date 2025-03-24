package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
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

	err = h.convertBalanceIfNeeded(ctx, withdrawAddr)
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

	err = h.convertBalanceIfNeeded(ctx, withdrawAddr)
	if err != nil {
		return err
	}

	return nil
}

// convertBalanceIfNeeded converts the bond denom balance of a given address from the coin to the ERC20 token.
// If the bond denom is not registered as an ERC20 token, or if the balance is zero, it returns nil.
func (h Hooks) convertBalanceIfNeeded(ctx sdk.Context, addr sdk.AccAddress) error {
	// Check if the bond denom is registered as an ERC20 token
	bondDenom := h.DistKeeper.stakingKeeper.GetParams(ctx).BondDenom
	if !h.DistKeeper.erc20k.IsDenomRegistered(ctx, bondDenom) {
		return nil
	}

	// get the balance, and convert
	balance := h.DistKeeper.bankKeeper.GetBalance(ctx, addr, bondDenom)
	if balance.IsZero() {
		return nil
	}

	// Create a MsgConvertCoin message
	msg := erc20types.NewMsgConvertCoin(balance, common.BytesToAddress(addr), addr)

	// Call the ERC20 keeper to convert the coin
	_, err := h.DistKeeper.erc20k.ConvertCoin(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return errorsmod.Wrap(err, "failed to convert coin")
	}

	return nil
}
