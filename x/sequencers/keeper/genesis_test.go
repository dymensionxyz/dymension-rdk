package keeper_test

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"

	"github.com/stretchr/testify/require"
)

func TestInitAndExportGenesis(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	exp := types.GenesisState{
		Params: types.DefaultParams(),
		Sequencers: []types.Sequencer{
			{
				Validator:  &utils.Proposer,
				RewardAddr: utils.OperatorAcc().String(),
			},
		},
	}
	k.InitGenesis(ctx, exp)
	got := k.ExportGenesis(ctx)
	require.Equal(t, &exp, got)
}
