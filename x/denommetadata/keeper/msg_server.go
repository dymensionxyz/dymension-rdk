package keeper

import (
	"context"

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
	return &types.MsgUpdateDenomMetadataResponse{}, nil
}
