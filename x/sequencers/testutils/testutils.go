package testutils

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
)

func NewTestSequencerKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(t, app)
	return k, ctx
}

/*
func NewTestSequencerKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	t_storeKey := sdk.NewTransientStoreKey("t_" + types.StoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(t_storeKey, storetypes.StoreTypeTransient, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	encCdc := testutils.MakeEncodingConfig()
	cdc := encCdc.Codec

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

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, log.NewNopLogger()).WithMultiStore(stateStore)
	// Initialize default params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
*/
