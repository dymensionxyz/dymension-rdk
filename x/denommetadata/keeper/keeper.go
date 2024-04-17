package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

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
