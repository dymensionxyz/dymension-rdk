package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
)

// CreateValidatorERC20 implements types.MsgServer.
func (k Keeper) CreateValidatorERC20(goCtx context.Context, msg *types.MsgCreateValidatorERC20) (*stakingtypes.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	stakingMsgServer := stakingkeeper.NewMsgServerImpl(k.Keeper)

	// Convert if needed

	if k.erc20k.IsDenomRegistered(ctx, msg.Value.Value.Denom) {
		delegatorAddress, err := sdk.AccAddressFromBech32(msg.Value.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		msg := erc20types.NewMsgConvertCoin(msg.Value.Value, common.BytesToAddress(delegatorAddress), delegatorAddress)
		if _, err = k.erc20k.ConvertCoin(sdk.WrapSDKContext(ctx), msg); err != nil {
			k.Logger(ctx).Error("Failed to convert coin", "err", err, "delegator", delegatorAddress)
			return nil, err
		}
	}

	// call create validator
	return stakingMsgServer.CreateValidator(sdk.WrapSDKContext(ctx), &msg.Value)
}

// DelegateERC20 implements types.MsgServer.
func (k Keeper) DelegateERC20(goCtx context.Context, msg *types.MsgDelegateERC20) (*stakingtypes.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	stakingMsgServer := stakingkeeper.NewMsgServerImpl(k.Keeper)

	// Convert if needed
	if k.erc20k.IsDenomRegistered(ctx, msg.Value.Amount.Denom) {
		delegatorAddress, err := sdk.AccAddressFromBech32(msg.Value.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		msg := erc20types.NewMsgConvertCoin(msg.Value.Amount, common.BytesToAddress(delegatorAddress), delegatorAddress)
		if _, err = k.erc20k.ConvertCoin(sdk.WrapSDKContext(ctx), msg); err != nil {
			k.Logger(ctx).Error("Failed to convert coin", "err", err, "delegator", delegatorAddress)
			return nil, err
		}
	}

	// call delegate
	return stakingMsgServer.Delegate(sdk.WrapSDKContext(ctx), &msg.Value)
}
