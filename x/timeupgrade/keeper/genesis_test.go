package keeper_test

import (
	"testing"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/types"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/require"
)

func TestInitAndExportGenesis(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestTimeupgradeKeeperFromApp(app)

	exp := types.GenesisState{
		Plan: &upgradetypes.Plan{
			Name:   "fooname",
			Height: 7,
			Info:   "fooinfo",
		},
		Timestamp: &prototypes.Timestamp{Seconds: 7},
	}
	k.InitGenesis(ctx, &exp)
	got := k.ExportGenesis(ctx)
	require.Equal(t, &exp, got)
}
