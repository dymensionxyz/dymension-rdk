package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

var _ types.MsgServer = &Keeper{}

// CreateDenomMetadata create the denom metadata in bank module
func (k Keeper) CreateDenomMetadata(
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

	found := k.bankKeeper.HasDenomMetaData(ctx, msg.TokenMetadata.Base)
	if found {
		return nil, types.ErrDenomAlreadyExists
	}

	k.bankKeeper.SetDenomMetaData(ctx, msg.TokenMetadata)
	// set hook after denom metadata creation
	err := k.hooks.AfterDenomMetadataCreation(ctx, msg.TokenMetadata)
	if err != nil {
		return nil, fmt.Errorf("error in after denom metadata creation hook: %w", err)
	}

	return &types.MsgCreateDenomMetadataResponse{}, nil
}

// UpdateDenomMetadata update the denom metadata in bank module
func (k Keeper) UpdateDenomMetadata(
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

	found := k.bankKeeper.HasDenomMetaData(ctx, msg.TokenMetadata.Base)
	if !found {
		return nil, types.ErrDenomDoesNotExist
	}

	k.bankKeeper.SetDenomMetaData(ctx, msg.TokenMetadata)

	// set hook after denom metadata update
	err := k.hooks.AfterDenomMetadataUpdate(ctx, msg.TokenMetadata)
	if err != nil {
		return nil, fmt.Errorf("error in after denom metadata update hook: %w", err)
	}

	return &types.MsgUpdateDenomMetadataResponse{}, nil
}
