package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

// Keeper of this module maintains distributing tokens to all stakers.
type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	sequencerKeeper types.SequencerKeeper
	bankKeeper      types.BankKeeper
	transferKeeper  types.TransferKeeper
	hooks           types.MultiDenomMetadataHooks
}

// NewKeeper creates new instances of the Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	sk types.SequencerKeeper,
	bk types.BankKeeper,
	tk types.TransferKeeper,
	hooks types.MultiDenomMetadataHooks,
) Keeper {
	return Keeper{
		storeKey:        storeKey,
		cdc:             cdc,
		sequencerKeeper: sk,
		bankKeeper:      bk,
		transferKeeper:  tk,
		hooks:           hooks,
	}
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
