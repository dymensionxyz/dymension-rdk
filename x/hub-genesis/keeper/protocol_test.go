package keeper

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
)

func TestProtocol(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestHubGenesisKeeperFromApp(app)
}
