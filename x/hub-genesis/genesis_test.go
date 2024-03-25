package hub_genesis_test

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/nullify"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	hub_genesis "github.com/dymensionxyz/dymension-rdk/x/hub-genesis"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := &types.GenesisState{}

	app := utils.Setup(t, false)
	// HubGenesis Keeper
	k, ctx := testkeepers.NewTestHubGenesisKeeperFromApp(app)

	ak := app.AccountKeeper

	hub_genesis.InitGenesis(ctx, *k, ak, genesisState)
	got := hub_genesis.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
