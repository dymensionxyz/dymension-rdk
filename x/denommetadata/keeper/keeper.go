package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

// Keeper of this module maintains distributing tokens to all stakers.
type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	paramSpace paramtypes.Subspace

	bankKeeper     types.BankKeeper
	transferKeeper types.TransferKeeper
	hooks          types.MultiDenomMetadataHooks
}

// NewKeeper creates new instances of the Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bk types.BankKeeper,
	tk types.TransferKeeper,
	hooks types.MultiDenomMetadataHooks,
	paramSpace paramtypes.Subspace,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:       storeKey,
		cdc:            cdc,
		paramSpace:     paramSpace,
		bankKeeper:     bk,
		transferKeeper: tk,
		hooks:          hooks,
	}
}

func (k Keeper) GetDenomMetadata(ctx sdk.Context, denomHash transfertypes.DenomTrace) (md types.DenomMetadata, err error) {
	tokenMetadata, ok := k.bankKeeper.GetDenomMetaData(ctx, denomHash.IBCDenom())
	if !ok {
		err = banktypes.ErrDenomMetadataNotFound
		return
	}

	denomTrace, ok := k.transferKeeper.GetDenomTrace(ctx, denomHash.Hash())
	if !ok {
		err = fmt.Errorf("denom trace not found for denom %s", denomHash.IBCDenom())
		return
	}

	md = types.DenomMetadata{
		TokenMetadata: tokenMetadata,
		DenomTrace:    denomTrace.GetFullDenomPath(),
	}

	return
}

// CreateDenomMetadata create the denom metadata in bank module
func (k Keeper) CreateDenomMetadata(ctx sdk.Context, metadatas ...types.DenomMetadata) error {
	for _, metadata := range metadatas {
		found := k.bankKeeper.HasDenomMetaData(ctx, metadata.TokenMetadata.Base)
		if found {
			return types.ErrDenomAlreadyExists
		}

		k.bankKeeper.SetDenomMetaData(ctx, metadata.TokenMetadata)
		// set hook after denom metadata creation
		if err := k.hooks.AfterDenomMetadataCreation(ctx, metadata.TokenMetadata); err != nil {
			return fmt.Errorf("error in after denom metadata creation hook: %w", err)
		}

		// construct the denomination trace from the full raw denomination
		denomTrace := transfertypes.ParseDenomTrace(metadata.DenomTrace)

		traceHash := denomTrace.Hash()
		if !k.transferKeeper.HasDenomTrace(ctx, traceHash) {
			k.transferKeeper.SetDenomTrace(ctx, denomTrace)
		}
	}

	return nil
}

// UpdateDenomMetadata update the denom metadata in bank module
func (k Keeper) UpdateDenomMetadata(ctx sdk.Context, metadatas ...types.DenomMetadata) error {
	for _, metadata := range metadatas {
		found := k.bankKeeper.HasDenomMetaData(ctx, metadata.TokenMetadata.Base)
		if !found {
			return types.ErrDenomDoesNotExist
		}

		k.bankKeeper.SetDenomMetaData(ctx, metadata.TokenMetadata)

		// set hook after denom metadata update
		if err := k.hooks.AfterDenomMetadataUpdate(ctx, metadata.TokenMetadata); err != nil {
			return fmt.Errorf("error in after denom metadata update hook: %w", err)
		}

		// construct the denomination trace from the full raw denomination
		denomTrace := transfertypes.ParseDenomTrace(metadata.DenomTrace)

		traceHash := denomTrace.Hash()
		if k.transferKeeper.HasDenomTrace(ctx, traceHash) {
			k.transferKeeper.SetDenomTrace(ctx, denomTrace) // nolint: errcheck
		}
	}

	return nil
}

// GetParams returns the total set of denommetadata parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// SetParams sets the total set of denommetadata parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// IsAddressPermissioned checks if the given address is permissioned to create or update denom metadata
func (k Keeper) IsAddressPermissioned(ctx sdk.Context, address string) bool {
	params := k.GetParams(ctx)
	for _, PermissionedAddress := range params.AllowedAddresses {
		if PermissionedAddress == address {
			return true
		}
	}
	return false
}

// SetHooks set the denommetadata hooks
func (k *Keeper) SetHooks(sh types.MultiDenomMetadataHooks) {
	if k.hooks != nil {
		panic("cannot set rollapp hooks twice")
	}
	k.hooks = sh
}

// GetHooks get the denommetadata hooks
func (k *Keeper) GetHooks() types.MultiDenomMetadataHooks {
	return k.hooks
}
