package keeper_test

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	utils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/mint/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestMintKeeperFromApp(app)

	state := types.GenesisState{}

	state.Params = types.DefaultParams()
	state.Minter = types.InitialMinter()

	k.InitGenesis(ctx, &state)
	got := k.ExportGenesis(ctx)
	require.NotNil(t, got)

	require.Equal(t, state.Params, got.Params)
	require.Equal(t, state.Minter, got.Minter)
}
