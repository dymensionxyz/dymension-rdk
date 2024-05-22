package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
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

// Deprecated
// CreateDenomMetadata create the denom metadata in bank module
func (k msgServer) CreateDenomMetadata(
	goCtx context.Context,
	msg *types.MsgCreateDenomMetadata,
) (*types.MsgCreateDenomMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	if !k.IsAddressPermissioned(ctx, msg.SenderAddress) {
		return nil, types.ErrNoPermission
	}

	if err := k.Keeper.CreateDenomMetadata(ctx, msg.Metadatas...); err != nil {
		return nil, fmt.Errorf("failed to create denom metadata: %w", err)
	}

	return &types.MsgCreateDenomMetadataResponse{}, nil
}

// Deprecated
// UpdateDenomMetadata update the denom metadata in bank module
func (k msgServer) UpdateDenomMetadata(
	goCtx context.Context,
	msg *types.MsgUpdateDenomMetadata,
) (*types.MsgUpdateDenomMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	if !k.IsAddressPermissioned(ctx, msg.SenderAddress) {
		return nil, types.ErrNoPermission
	}

	if err := k.Keeper.UpdateDenomMetadata(ctx, msg.Metadatas...); err != nil {
		return nil, fmt.Errorf("failed to update denom metadata: %w", err)
	}

	return &types.MsgUpdateDenomMetadataResponse{}, nil
}
