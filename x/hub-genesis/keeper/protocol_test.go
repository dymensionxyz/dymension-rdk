package keeper

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
)

func TestProtocol(t *testing.T) {
	/*
		Cases:
		1. Transfers are enabled if genesis accounts is zero
	*/

	t.Run("transfers are enabled immediately if there are no genesis accounts", func(t *testing.T) {
		app := utils.Setup(t, false)
		k, ctx := testkeepers.NewTestHubGenesisKeeperFromApp(app)
		k.on
	})
}
