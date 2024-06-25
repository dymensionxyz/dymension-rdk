package keeper_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

func TestGenesis(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestHubKeeperFromApp(app)

	expect := &types.GenesisState{
		State: types.State{},
	}
	k.InitGenesis(ctx, expect)
	got := k.ExportGenesis(ctx)
	require.NotNil(t, got)

	require.Equal(t, expect.State, got.State)
}
