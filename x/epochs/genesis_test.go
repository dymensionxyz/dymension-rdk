package epochs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testkeepers "github.com/dymensionxyz/rollapp/testutil/keepers"
	"github.com/dymensionxyz/rollapp/testutil/nullify"
	"github.com/dymensionxyz/rollapp/testutil/utils"
	"github.com/dymensionxyz/rollapp/x/epochs"
	"github.com/dymensionxyz/rollapp/x/epochs/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		// this line is used by starport scaffolding # genesis/test/state
	}

	app := utils.Setup(t, false)
	//EpochsKeeper
	k, ctx := testkeepers.NewTestEpochKeeperFromApp(t, app)

	epochs.InitGenesis(ctx, *k, genesisState)
	got := epochs.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
