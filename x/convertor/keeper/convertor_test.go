package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	convertorkeeper "github.com/dymensionxyz/dymension-rdk/x/convertor/keeper"
	hubtypes "github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

const (
	testDenom         = "ibc/ABC123"
	bridgeDecimals    = 6
	rollappDecimals   = 18
	decimalDifference = rollappDecimals - bridgeDecimals // 12
)

// setupTestApp creates a test app with all keepers initialized
func setupTestApp(t *testing.T) (*convertorkeeper.Keeper, *app.App, sdk.Context) {
	t.Helper()
	testApp := utils.Setup(t, false)
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "test-1"})

	// Register the test denom metadata with 18 decimals (rollapp standard)
	metadata := banktypes.Metadata{
		Base: testDenom,
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: testDenom, Exponent: 0},
			{Denom: "test", Exponent: rollappDecimals},
		},
		Display: "test",
	}
	testApp.BankKeeper.SetDenomMetaData(ctx, metadata)

	// Set up the decimal conversion pair in hub keeper
	pair := hubtypes.DecimalConversionPair{
		FromToken:    testDenom,
		FromDecimals: bridgeDecimals,
	}
	err := testApp.HubKeeper.SetDecimalConversionPair(ctx, pair)
	require.NoError(t, err)

	return &testApp.TransferKeeper, testApp, ctx
}

func TestConversionRequired(t *testing.T) {
	k, _, ctx := setupTestApp(t)

	tests := []struct {
		name     string
		denom    string
		expected bool
	}{
		{
			name:     "conversion required for registered denom",
			denom:    testDenom,
			expected: true,
		},
		{
			name:     "no conversion for other denom",
			denom:    "other/denom",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			required, err := k.ConversionRequired(ctx, tt.denom)
			require.NoError(t, err)
			require.Equal(t, tt.expected, required)
		})
	}
}

func TestConvertToBridgeAmt(t *testing.T) {
	k, _, ctx := setupTestApp(t)

	tests := []struct {
		name           string
		rollappAmount  string
		expectedBridge string
		description    string
	}{
		{
			name:           "exact conversion - no dust",
			rollappAmount:  "1000000000000000000", // 1e18
			expectedBridge: "1000000",             // 1e6
			description:    "1 token with 18 decimals converts to 1 token with 6 decimals",
		},
		{
			name:           "with dust - 1 wei",
			rollappAmount:  "1000000000000000001", // 1e18 + 1
			expectedBridge: "1000000",             // 1e6 (truncated)
			description:    "1 wei is lost in conversion",
		},
		{
			name:           "large dust",
			rollappAmount:  "1234567890123456789",
			expectedBridge: "1234567",
			description:    "loses 890123456789 wei",
		},
		{
			name:           "maximum dust - almost 1 bridge unit",
			rollappAmount:  "1000000999999999999", // 1e18 + (1e12 - 1)
			expectedBridge: "1000000",             // 1e6
			description:    "999999999999 wei lost (maximum possible dust)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rollappAmt, ok := math.NewIntFromString(tt.rollappAmount)
			require.True(t, ok)

			expectedBridgeAmt, ok := math.NewIntFromString(tt.expectedBridge)
			require.True(t, ok)

			// Convert from rollapp (18 decimals) to bridge (6 decimals)
			bridgeAmt, err := k.ConvertToBridgeAmt(ctx, rollappAmt)
			require.NoError(t, err)
			require.Equal(t, expectedBridgeAmt, bridgeAmt, tt.description)
		})
	}
}

func TestConvertFromBridgeAmt(t *testing.T) {
	k, _, ctx := setupTestApp(t)

	tests := []struct {
		name            string
		bridgeAmount    string
		expectedRollapp string
		description     string
	}{
		{
			name:            "exact conversion - no dust",
			bridgeAmount:    "1000000",             // 1e6
			expectedRollapp: "1000000000000000000", // 1e18
			description:     "1 token with 6 decimals converts to 1 token with 18 decimals",
		},
		{
			name:            "large amount",
			bridgeAmount:    "1234567",
			expectedRollapp: "1234567000000000000",
			description:     "larger amount scales up correctly",
		},
		{
			name:            "small amount",
			bridgeAmount:    "1",
			expectedRollapp: "1000000000000",
			description:     "1 smallest bridge unit = 1e12 rollapp units",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bridgeAmt, ok := math.NewIntFromString(tt.bridgeAmount)
			require.True(t, ok)

			expectedRollappAmt, ok := math.NewIntFromString(tt.expectedRollapp)
			require.True(t, ok)

			// Convert from bridge (6 decimals) to rollapp (18 decimals)
			rollappAmt, err := k.ConvertFromBridgeAmt(ctx, bridgeAmt)
			require.NoError(t, err)
			require.Equal(t, expectedRollappAmt, rollappAmt, tt.description)
		})
	}
}

func TestConvertAmountRoundTrip(t *testing.T) {
	k, _, ctx := setupTestApp(t)

	// Test that converting to bridge and back gives us the expected values for dust calculation
	tests := []struct {
		name                  string
		originalAmount        string
		expectedBridgeAmount  string
		expectedRollappAmount string
		expectedDust          string
	}{
		{
			name:                  "no dust - exact conversion",
			originalAmount:        "1000000000000000000", // 1e18
			expectedBridgeAmount:  "1000000",             // 1e6
			expectedRollappAmount: "1000000000000000000", // 1e18
			expectedDust:          "0",
		},
		{
			name:                  "1 wei dust",
			originalAmount:        "1000000000000000001", // 1e18 + 1
			expectedBridgeAmount:  "1000000",             // 1e6
			expectedRollappAmount: "1000000000000000000", // 1e18
			expectedDust:          "1",
		},
		{
			name:                  "maximum dust - almost 1 smallest bridge unit",
			originalAmount:        "1000000999999999999", // 1e18 + 999999999999
			expectedBridgeAmount:  "1000000",             // 1e6
			expectedRollappAmount: "1000000000000000000", // 1e18
			expectedDust:          "999999999999",        // almost 1e12
		},
		{
			name:                  "larger amount with dust",
			originalAmount:        "123456789012345678",
			expectedBridgeAmount:  "123456",
			expectedRollappAmount: "123456000000000000",
			expectedDust:          "789012345678",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalAmount, ok := math.NewIntFromString(tt.originalAmount)
			require.True(t, ok, "failed to parse original amount")

			expectedBridgeAmount, ok := math.NewIntFromString(tt.expectedBridgeAmount)
			require.True(t, ok, "failed to parse expected bridge amount")

			expectedRollappAmount, ok := math.NewIntFromString(tt.expectedRollappAmount)
			require.True(t, ok, "failed to parse expected rollapp amount")

			expectedDust, ok := math.NewIntFromString(tt.expectedDust)
			require.True(t, ok, "failed to parse expected dust")

			// 1. Convert to bridge decimals (this is what ConvertToBridgeAmt does)
			bridgeAmt, err := k.ConvertToBridgeAmt(ctx, originalAmount)
			require.NoError(t, err)
			require.Equal(t, expectedBridgeAmount, bridgeAmt, "bridge amount mismatch")

			// 2. Convert back to rollapp decimals (this is what ConvertFromBridgeAmt does)
			rollappAmt, err := k.ConvertFromBridgeAmt(ctx, bridgeAmt)
			require.NoError(t, err)
			require.Equal(t, expectedRollappAmount, rollappAmt, "rollapp amount mismatch")

			// 3. Calculate dust
			dust := originalAmount.Sub(rollappAmt)
			if expectedDust.IsZero() {
				require.True(t, dust.IsZero(), "dust should be zero")
			} else {
				require.Equal(t, expectedDust, dust, "dust mismatch")
			}
		})
	}
}

func TestBurnAndMintCoins(t *testing.T) {
	k, testApp, ctx := setupTestApp(t)

	// Create a test account and fund it
	testAddr := sdk.AccAddress("test_address______")
	initialAmount := math.NewInt(1000000)

	// Mint coins to the test address
	err := k.MintCoins(ctx, testAddr, sdk.NewCoin(testDenom, initialAmount))
	require.NoError(t, err)

	// Verify the balance using the bank keeper
	balance := testApp.BankKeeper.GetBalance(ctx, testAddr, testDenom)
	require.Equal(t, initialAmount, balance.Amount)

	// Burn some coins
	burnAmount := math.NewInt(500000)
	err = k.BurnCoins(ctx, testAddr, sdk.NewCoin(testDenom, burnAmount))
	require.NoError(t, err)

	// Verify the new balance
	balance = testApp.BankKeeper.GetBalance(ctx, testAddr, testDenom)
	expectedBalance := initialAmount.Sub(burnAmount)
	require.Equal(t, expectedBalance, balance.Amount)
}
