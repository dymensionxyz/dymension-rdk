package keeper_test

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitAndExportGenesis(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	params := k.GetParams(ctx)
	seqs := k.GetAllSequencers(ctx)
	require.Equal(t, 1, len(seqs))
	expectedOperator := seqs[0].GetOperator().String()
	require.NotEmpty(t, expectedOperator)

	genState := k.ExportGenesis(ctx)
	assert.Equal(t, expectedOperator, genState.GenesisOperatorAddress)
	assert.Equal(t, params, genState.Params)

	// Test InitGenesis
	genState.Params.HistoricalEntries = 100

	_ = k.InitGenesis(ctx, *genState)
	assert.Equal(t, genState.Params, k.GetParams(ctx))
}
