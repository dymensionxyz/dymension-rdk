package sequencers_test

import (
	"testing"

	"github.com/dymensionxyz/rollapp/testutil/nullify"
	"github.com/dymensionxyz/rollapp/x/sequencers"
	"github.com/dymensionxyz/rollapp/x/sequencers/testutils"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	k, ctx := testutils.NewTestSequencer(t)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	sequencers.InitGenesis(ctx, *k, genesisState)
	got := sequencers.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
