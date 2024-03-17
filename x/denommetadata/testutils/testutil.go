package testutils

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
	rollapp "github.com/dymensionxyz/rollapp/app"
)

// NewTestDenommetadataKeeper creates a new denommetadata keeper for testing
func NewTestDenommetadataKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	app := utils.Setup(t, false)
	bankKeeper, _ := testkeepers.NewTestBankKeeperFromApp(t, app)

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	fmt.Println("storeKey: ", storeKey)
	t_storeKey := sdk.NewTransientStoreKey("t_" + types.StoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(t_storeKey, storetypes.StoreTypeTransient, nil)
	require.NoError(t, stateStore.LoadLatestVersion(), "loading latest version failed")

	encCdc := rollapp.MakeEncodingConfig()
	cdc := encCdc.Codec

	paramsSubspace := typesparams.NewSubspace(
		cdc,
		types.Amino,
		storeKey,
		t_storeKey,
		"DenommetadataParams",
	)

	denommetadataKeeper := keeper.NewKeeper(
		storeKey,
		cdc,
		*bankKeeper,
		paramsSubspace,
	)

	ctx := sdk.NewContext(nil, tmproto.Header{}, false, log.NewNopLogger()).WithMultiStore(stateStore)
	// Initialize default params
	denommetadataKeeper.SetParams(ctx, types.DefaultParams())

	return &denommetadataKeeper, ctx

}
