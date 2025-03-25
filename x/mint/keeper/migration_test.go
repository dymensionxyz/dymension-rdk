package keeper_test

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	utils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
	keeper "github.com/dymensionxyz/dymension-rdk/x/mint/keeper"
	v2 "github.com/dymensionxyz/dymension-rdk/x/mint/types/migrations/v2"
	"github.com/stretchr/testify/require"
)

func TestMigrate1to2(t *testing.T) {
	/* ---------------------------------- setup --------------------------------- */
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestMintKeeperFromApp(app)

	// set old style params
	oldParams := v2.DefaultParams()
	m := keeper.NewMigrator(*k, app.ParamsKeeper.Subspace("empty"))
	m.SetOldParams(ctx, oldParams)

	// no denom in minter
	minter := k.GetMinter(ctx)
	minter.MintDenom = ""
	k.SetMinter(ctx, minter)

	err := m.Migrate1to2(ctx)
	require.NoError(t, err)

	minter = k.GetMinter(ctx)
	require.Equal(t, oldParams.MintDenom, minter.MintDenom)
}
