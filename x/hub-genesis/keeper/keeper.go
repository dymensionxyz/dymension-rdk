package keeper

import (
	"fmt"

	"github.com/dymensionxyz/dymension-rdk/utils/whitelistedrelayer"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramstore paramtypes.Subspace

	ak types.AccountKeeper
	bk types.BankKeeper
	mk types.MintKeeper
	dk whitelistedrelayer.DistrK
	sk whitelistedrelayer.SeqK
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	mk types.MintKeeper,
	dk whitelistedrelayer.DistrK,
	sk whitelistedrelayer.SeqK,

) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		table := types.ParamKeyTable()
		ps = ps.WithKeyTable(table)
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

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,
		ak:         ak,
		bk:         bk,
		mk:         mk,
		sk:         sk,
		dk:         dk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
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
