package keeper_test

import (
	"testing"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	utils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
	keeper "github.com/dymensionxyz/dymension-rdk/x/mint/keeper"
	minttypes "github.com/dymensionxyz/dymension-rdk/x/mint/types"
	"github.com/stretchr/testify/require"
)

func TestMigrate1to2(t *testing.T) {
	/* ---------------------------------- setup --------------------------------- */
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestMintKeeperFromApp(app)

	// raw access to the mint params subspace
	tmpMintSubspace := paramstypes.NewSubspace(app.AppCodec(), app.LegacyAmino(), app.GetKey(paramstypes.StoreKey), app.GetTKey(paramstypes.TStoreKey), minttypes.ModuleName)

	// migrator sets old v1 keyTable on this subspace
	m := keeper.NewMigrator(*k, tmpMintSubspace)
	tmpMintSubspace.Set(ctx, []byte("MintDenom"), "denom")           // simulate v1 denom in params
	tmpMintSubspace.Set(ctx, []byte("MintEpochIdentifier"), "epoch") // assertion that the subpsace is actually changed in the x/mint keeper

	// make sure the above changes reflected in the keeper
	params := k.GetParams(ctx)
	require.Equal(t, "epoch", params.MintEpochIdentifier)

	err := m.Migrate1to2(ctx)
	require.NoError(t, err)

	require.Equal(t, "denom", k.GetMinter(ctx).MintDenom)
}
