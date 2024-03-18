package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/testutils"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

func TestCreateDenomMetadata(t *testing.T) {
	k, ctx := testutils.NewTestDenommetadataKeeper(t)

	// Prepare the test message for creating denom metadata
	createMsg := &types.MsgCreateDenomMetadata{
		SenderAddress: "cosmos1s77x8wr2gzdhq8gt8c085vate0s23xu9u80wtx",
		TokenMetadata: banktypes.Metadata{
			Name:        "Dymension Hub token",
			Symbol:      "DYM",
			Description: "Denom metadata for DYM.",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "adym", Exponent: uint32(0), Aliases: []string{}},
				{Denom: "DYM", Exponent: uint32(18), Aliases: []string{}},
			},
			Base:    "adym",
			Display: "DYM",
		},
	}

	// Test permission error
	createMsg.SenderAddress = "cosmos1s77x8wr2gzdhq8gt8c085vate0s23xu9u80wtx"
	_, err := k.CreateDenomMetadata(sdk.WrapSDKContext(ctx), createMsg)
	require.ErrorIs(t, err, types.ErrNoPermission, "should return permission error")

	// Set allowed addresses
	initialParams := types.DefaultParams()
	initialParams.AllowedAddresses = []string{"cosmos1s77x8wr2gzdhq8gt8c085vate0s23xu9u80wtx", "cosmos1gusne8eh37myphx09hgdsy85zpl2t0kzdvu3en"}
	k.SetParams(ctx, initialParams)

	// Test creating denom metadata successfully
	_, err = k.CreateDenomMetadata(sdk.WrapSDKContext(ctx), createMsg)
	require.NoError(t, err, "creating denom metadata with allowed address should not error")

	// Test creating duplicate denom metadata
	_, err = k.CreateDenomMetadata(sdk.WrapSDKContext(ctx), createMsg)
	require.ErrorIs(t, err, types.ErrDenomAlreadyExists, "creating duplicate denom metadata should fail")
}
