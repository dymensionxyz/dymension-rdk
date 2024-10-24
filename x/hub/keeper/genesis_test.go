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
		State: types.State{
			Hub: types.Hub{
				RegisteredDenoms: []*types.RegisteredDenom{
					{
						Base: "adym",
					}, {
						Base: "ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2",
					},
				},
			},
		},
	}
	k.InitGenesis(ctx, expect)
	got := k.ExportGenesis(ctx)
	require.NotNil(t, got)

	require.Equal(t, expect.State, got.State)
}
