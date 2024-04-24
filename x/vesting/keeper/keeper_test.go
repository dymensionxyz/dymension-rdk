package keeper_test

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/vesting/types"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	// Setup the test environment
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestVestingKeeperFromApp(app)

	// Set some initial parameters
	initialParams := types.DefaultParams()
	initialParams.AllowedAddresses = []string{"cosmos19crd4fwzm9qtf5ln5l3e2vmquhevjwprk8tgxp", "cosmos1gusne8eh37myphx09hgdsy85zpl2t0kzdvu3en"} // Example addresses
	k.SetParams(ctx, initialParams)

	// Retrieve the parameters
	retrievedParams := k.GetParams(ctx)

	// Assert that the retrieved parameters match the initial ones
	require.Equal(t, initialParams, retrievedParams, "retrieved parameters should match the initial ones")

	// Test setting and getting a different set of parameters
	updatedParams := initialParams
	updatedParams.AllowedAddresses = append(updatedParams.AllowedAddresses, "cosmos1s77x8wr2gzdhq8gt8c085vate0s23xu9u80wtx")
	k.SetParams(ctx, updatedParams)
	retrievedParams = k.GetParams(ctx)
	require.Equal(t, updatedParams, retrievedParams, "retrieved parameters should match the updated ones")
}
