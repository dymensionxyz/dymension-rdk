package testutils

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/keeper"
)

// NewTestDenommetadataKeeper creates a new denommetadata keeper for testing
func NewTestDenommetadataKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestDenommetadataKeeperFromApp(app)
	return k, ctx
}
