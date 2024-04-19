package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dymensionxyz/dymension-rdk/testutil/nullify"
	"github.com/dymensionxyz/dymension-rdk/x/epochs/types"
)

// TODO: add specific test and scenario
func TestEpochsInitAndExportGenesis(t *testing.T) {
	genesisState := types.GenesisState{}
	ctx, epochsKeeper := Setup(t)

	epochsKeeper.InitGenesis(ctx, genesisState)
	genesis := epochsKeeper.ExportGenesis(ctx)
	require.NotNil(t, genesis)
	require.Len(t, genesis.Epochs, 5)
	nullify.Fill(&genesis)
	nullify.Fill(genesis)
}
