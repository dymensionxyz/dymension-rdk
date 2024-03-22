package testutils

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
)

// NewTestDenommetadataKeeper creates a new denommetadata keeper for testing
func NewTestDenommetadataKeeper(t *testing.T) (*app.App, sdk.Context) {
	tapp := utils.Setup(t, false)
	_, ctx := testkeepers.NewTestDenommetadataKeeperFromApp(tapp)
	return tapp, ctx
}
