package testutils

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/dymensionxyz/rollapp/x/sequencers/keeper"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
)

func NewTestSequencer(ctx sdk.Context) *keeper.Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		storeKey,
		"SequencerParams",
	)
	k := keeper.NewKeeper(
		cdc,
		storeKey,
		paramsSubspace,
	)

	// Initialize default params
	k.SetParams(ctx, types.DefaultParams())

	return k
}
