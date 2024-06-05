package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateGasTank defines a method to create a new gas tank
func (m msgServer) CreateGasTank(goCtx context.Context, msg *types.MsgCreateGasTank) (*types.MsgCreateGasTankResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.CreateGasTank(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCreateGasTankResponse{}, nil
}

// UpdateGasTankStatus defines a method to update the active status of gas tank
func (m msgServer) UpdateGasTankStatus(goCtx context.Context, msg *types.MsgUpdateGasTankStatus) (*types.MsgUpdateGasTankStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.UpdateGasTankStatus(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgUpdateGasTankStatusResponse{}, nil
}

// UpdateGasTankConfigs defines a method to update a gas tank
func (m msgServer) UpdateGasTankConfigs(goCtx context.Context, msg *types.MsgUpdateGasTankConfig) (*types.MsgUpdateGasTankConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.UpdateGasTankConfig(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgUpdateGasTankConfigResponse{}, nil
}

// BlockConsumer defines a method to block a gas consumer
func (m msgServer) BlockConsumer(goCtx context.Context, msg *types.MsgBlockConsumer) (*types.MsgBlockConsumerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.BlockConsumer(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgBlockConsumerResponse{}, nil
}

// UnblockConsumer defines a method to unblock a consumer
func (m msgServer) UnblockConsumer(goCtx context.Context, msg *types.MsgUnblockConsumer) (*types.MsgUnblockConsumerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.UnblockConsumer(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgUnblockConsumerResponse{}, nil
}

// UpdateGasConsumerLimit defines a method to increase consumption limit for a consumer
func (m msgServer) UpdateGasConsumerLimit(goCtx context.Context, msg *types.MsgUpdateGasConsumerLimit) (*types.MsgUpdateGasConsumerLimitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := m.Keeper.UpdateGasConsumerLimit(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgUpdateGasConsumerLimitResponse{}, nil
}
