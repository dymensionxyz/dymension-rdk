package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	types2 "github.com/dymensionxyz/dymension-rdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
)

// CreateValidatorERC20 implements types.MsgServer.
func (k Keeper) CreateValidatorERC20(goCtx context.Context, msg *types2.MsgCreateValidator) (*types2.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	stakingMsgServer := stakingkeeper.NewMsgServerImpl(k.Keeper)

	// Convert
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.Value.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	if k.erc20k.IsDenomRegistered(ctx, msg.Value.Value.Denom) {
		msg := erc20types.NewMsgConvertCoin(msg.Value.Value, common.BytesToAddress(delegatorAddress), delegatorAddress)
		if _, err = k.erc20k.ConvertCoin(sdk.WrapSDKContext(ctx), msg); err != nil {
			k.Logger(ctx).Error("Failed to convert coin", "err", err, "delegator", delegatorAddress)
			return nil, err
		}
	}

	// call create validator
	res, err := stakingMsgServer.CreateValidator(sdk.WrapSDKContext(ctx), msg.Value)
	if err != nil {
		return nil, err
	}

	return &types2.MsgCreateValidatorResponse{Value: res}, nil
}

// DelegateERC20 implements types.MsgServer.
func (k Keeper) DelegateERC20(goCtx context.Context, msg *types2.MsgDelegate) (*types2.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	stakingMsgServer := stakingkeeper.NewMsgServerImpl(k.Keeper)

	// Convert
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.Value.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	if k.erc20k.IsDenomRegistered(ctx, msg.Value.Amount.Denom) {
		msg := erc20types.NewMsgConvertCoin(msg.Value.Amount, common.BytesToAddress(delegatorAddress), delegatorAddress)
		if _, err = k.erc20k.ConvertCoin(sdk.WrapSDKContext(ctx), msg); err != nil {
			k.Logger(ctx).Error("Failed to convert coin", "err", err, "delegator", delegatorAddress)
			return nil, err
		}
	}

	// call delegate
	res, err := stakingMsgServer.Delegate(sdk.WrapSDKContext(ctx), msg.Value)
	if err != nil {
		return nil, err
	}

	return &types2.MsgDelegateResponse{Value: res}, nil
}
