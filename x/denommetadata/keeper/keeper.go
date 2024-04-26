package keeper

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k *Keeper) isAddressPermissioned(ctx sdk.Context, address string) bool {
	logger := k.Logger(ctx)
	accAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		logger.Error("failed to extract account address from bech32: ", err)
		return false
	}

	return k.sequencerKeeper.HasPermission(ctx, accAddr, types.ModuleName)
}
