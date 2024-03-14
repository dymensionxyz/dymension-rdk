package testutils

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
)

// FIXME: remove.
func NewTestSequencerKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(t, app)
	return k, ctx
}
