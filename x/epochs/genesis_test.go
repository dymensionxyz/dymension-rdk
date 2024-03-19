package epochs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/nullify"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/epochs"
	"github.com/dymensionxyz/dymension-rdk/x/epochs/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{}

	app := utils.Setup(t, false)
	//EpochsKeeper
	k, ctx := testkeepers.NewTestEpochKeeperFromApp(app)

	epochs.InitGenesis(ctx, *k, genesisState)
	got := epochs.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
