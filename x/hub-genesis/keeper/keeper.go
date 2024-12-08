package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"

	"cosmossdk.io/collections"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/dymension-rdk/utils/collcompat"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramstore paramtypes.Subspace

	ak types.AccountKeeper
	bk types.BankKeeper
	mk types.MintKeeper

	gb types.GenesisBridgeSubmitter

	// key is port/channel. value is types.ChannelState
	PendingChannels collections.Map[string, uint64]
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	mk types.MintKeeper,
	gb types.GenesisBridgeSubmitter,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	if ak == nil {
		panic("account keeper cannot be nil")
	}
	if bk == nil {
		panic("bank keeper cannot be nil")
	}
	if mk == nil {
		panic("mint keeper cannot be nil")
	}

	sb := collections.NewSchemaBuilder(collcompat.NewKVStoreService(storeKey))

	k := Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,
		ak:         ak,
		bk:         bk,
		mk:         mk,
		gb:         gb,
		PendingChannels: collections.NewMap(
			sb,
			types.OngoingChannelsPrefix(),
			"ongoing_channels",
			collections.StringKey,
			collections.Uint64Value,
		),
	}
	return k
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) SetICS4Submitter(submitter types.GenesisBridgeSubmitter) {
	k.gb = submitter
}

// SetState sets the state.
func (k Keeper) SetState(ctx sdk.Context, state types.State) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.StateKey, k.cdc.MustMarshal(&state))
}

// GetState returns the state.
func (k Keeper) GetState(ctx sdk.Context) types.State {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.StateKey)
	if bz == nil {
		return types.State{}
	}
	var state types.State
	k.cdc.MustUnmarshal(bz, &state)
	return state
}

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// SetGenesisInfo sets the genesis info.
func (k Keeper) SetGenesisInfo(ctx sdk.Context, gInfo types.GenesisInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GenesisInfoKey, k.cdc.MustMarshal(&gInfo))
}

// GetGenesisInfo returns the genesis info.
func (k Keeper) GetGenesisInfo(ctx sdk.Context) types.GenesisInfo {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GenesisInfoKey)
	if bz == nil {
		return types.GenesisInfo{}
	}
	var gInfo types.GenesisInfo
	k.cdc.MustUnmarshal(bz, &gInfo)
	return gInfo
}

func (k Keeper) SetPendingChannel(ctx sdk.Context, portChannel types.PortAndChannel, status types.ChannelState) error {
	return k.PendingChannels.Set(ctx, portChannel.Key(), uint64(status))
}

func (k Keeper) ClearPendingChannels(ctx sdk.Context) error {
	return k.PendingChannels.Clear(ctx, nil)
}

func (k Keeper) IsPendingChannel(ctx sdk.Context, portChannel types.PortAndChannel) (bool, error) {
	return k.PendingChannels.Has(ctx, portChannel.Key())
}
