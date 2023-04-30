package testutils

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/dymensionxyz/rollapp/app"
	"github.com/dymensionxyz/rollapp/x/sequencers/keeper"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
)

func NewTestContext() sdk.Context {
	return sdk.NewContext(nil, tmproto.Header{}, false, log.NewNopLogger())
}

func NewTestSequencerKeeperFromApp(t *testing.T, app *app.App) (*keeper.Keeper, sdk.Context) {
	cdc := app.AppCodec()
	k := keeper.NewKeeper(
		cdc,
		app.GetKey(types.StoreKey),
		app.GetSubspace(types.ModuleName),
	)
	ctx := app.GetBaseApp().NewContext(false, tmproto.Header{})
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}

func NewTestSequencerKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	t_storeKey := sdk.NewTransientStoreKey("t_" + types.StoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(t_storeKey, sdk.StoreTypeTransient, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		t_storeKey,
		"SequencerParams",
	)
	k := keeper.NewKeeper(
		cdc,
		storeKey,
		paramsSubspace,
	)

	ctx := NewTestContext().WithMultiStore(stateStore)
	// Initialize default params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}

// NewValidator is a testing helper method to create validators in tests
func NewValidator(t testing.TB, operator sdk.ValAddress, pubKey cryptotypes.PubKey) stakingtypes.Validator {
	v, err := stakingtypes.NewValidator(operator, pubKey, stakingtypes.Description{})
	require.NoError(t, err)
	return v
}
