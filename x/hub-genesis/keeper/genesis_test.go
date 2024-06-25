package keeper_test

import (
	_ "embed"
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestHubGenesisKeeperFromApp(app)

	expect := &types.GenesisState{
		State: types.State{
			NumUnackedTransfers: 2,
		},
		UnackedTransferSeqNums: []uint64{42, 43},
	}
	k.InitGenesis(ctx, expect)
	got := k.ExportGenesis(ctx)
	require.NotNil(t, got)

	require.ElementsMatch(t, expect.UnackedTransferSeqNums, got.UnackedTransferSeqNums)
	require.Equal(t, expect.State.NumUnackedTransfers, got.State.NumUnackedTransfers)
}
