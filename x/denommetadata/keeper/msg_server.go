package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

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

	for _, metadata := range msg.Metadatas {
		found := k.bankKeeper.HasDenomMetaData(ctx, metadata.TokenMetadata.Base)
		if found {
			return nil, types.ErrDenomAlreadyExists
		}

		k.bankKeeper.SetDenomMetaData(ctx, metadata.TokenMetadata)
		// set hook after denom metadata creation
		err := k.hooks.AfterDenomMetadataCreation(ctx, metadata.TokenMetadata)
		if err != nil {
			return nil, fmt.Errorf("error in after denom metadata creation hook: %w", err)
		}

		// construct the denomination trace from the full raw denomination
		denomTrace := transfertypes.ParseDenomTrace(metadata.DenomTrace)

		traceHash := denomTrace.Hash()
		if !k.transferKeeper.HasDenomTrace(ctx, traceHash) {
			k.transferKeeper.SetDenomTrace(ctx, denomTrace)
		}
	}
	return &types.MsgCreateDenomMetadataResponse{}, nil
}

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

	for _, metadata := range msg.Metadatas {
		found := k.bankKeeper.HasDenomMetaData(ctx, metadata.TokenMetadata.Base)
		if !found {
			return nil, types.ErrDenomDoesNotExist
		}

		k.bankKeeper.SetDenomMetaData(ctx, metadata.TokenMetadata)

		// set hook after denom metadata update
		err := k.hooks.AfterDenomMetadataUpdate(ctx, metadata.TokenMetadata)
		if err != nil {
			return nil, fmt.Errorf("error in after denom metadata update hook: %w", err)
		}

		// construct the denomination trace from the full raw denomination
		denomTrace := transfertypes.ParseDenomTrace(metadata.DenomTrace)

		traceHash := denomTrace.Hash()
		if k.transferKeeper.HasDenomTrace(ctx, traceHash) {
			k.transferKeeper.SetDenomTrace(ctx, denomTrace)
		}
	}

	return &types.MsgUpdateDenomMetadataResponse{}, nil
}
