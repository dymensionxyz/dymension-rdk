package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/dymensionxyz/dymension-rdk/x/converter/types"
)

func TestConvertAmount(t *testing.T) {
	tests := []struct {
		name          string
		amount        string
		fromDecimals  uint32
		toDecimals    uint32
		expectedAmt   string
		expectedError bool
	}{
		{
			name:          "scale up from 6 to 18 decimals",
			amount:        "1000000",
			fromDecimals:  6,
			toDecimals:    18,
			expectedAmt:   "1000000000000000000",
			expectedError: false,
		},
		{
			name:          "scale down from 18 to 6 decimals",
			amount:        "1000000000000000000",
			fromDecimals:  18,
			toDecimals:    6,
			expectedAmt:   "1000000",
			expectedError: false,
		},
		{
			name:          "same decimals - no conversion",
			amount:        "1000000",
			fromDecimals:  6,
			toDecimals:    6,
			expectedAmt:   "1000000",
			expectedError: false,
		},
		{
			name:          "scale down with precision loss (truncated)",
			amount:        "1000000000000000001",
			fromDecimals:  18,
			toDecimals:    6,
			expectedAmt:   "1000000",
			expectedError: false,
		},
		{
			name:          "negative amount",
			amount:        "-1000000",
			fromDecimals:  6,
			toDecimals:    18,
			expectedAmt:   "",
			expectedError: true,
		},
		{
			name:          "zero amount",
			amount:        "0",
			fromDecimals:  6,
			toDecimals:    18,
			expectedAmt:   "0",
			expectedError: false,
		},
		{
			name:          "scale up from 8 to 18 decimals (BTC-like)",
			amount:        "100000000",
			fromDecimals:  8,
			toDecimals:    18,
			expectedAmt:   "1000000000000000000",
			expectedError: false,
		},
		{
			name:          "large amount scale up",
			amount:        "1000000000000",
			fromDecimals:  6,
			toDecimals:    18,
			expectedAmt:   "1000000000000000000000000",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount, ok := sdk.NewIntFromString(tt.amount)
			require.True(t, ok, "failed to parse amount")

			result, err := types.ConvertAmount(amount, tt.fromDecimals, tt.toDecimals)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				expected, ok := sdk.NewIntFromString(tt.expectedAmt)
				require.True(t, ok, "failed to parse expected amount")
				require.Equal(t, expected, result)
			}
		})
	}
}

func TestClearPrecisionLoss(t *testing.T) {
	tests := []struct {
		name          string
		amount        string
		fromDecimals  uint32
		toDecimals    uint32
		expectedAmt   string
		expectedDust  string
		description   string
		expectedError bool
	}{
		{
			name:          "no dust - exact conversion",
			amount:        "1000000000000000000", // 1e18
			fromDecimals:  18,
			toDecimals:    6,
			expectedAmt:   "1000000000000000000", // 1e18 (unchanged)
			expectedDust:  "0",
			description:   "1 token with no fractional part - no precision loss",
			expectedError: false,
		},
		{
			name:          "1 wei dust",
			amount:        "1000000000000000001", // 1e18 + 1
			fromDecimals:  18,
			toDecimals:    6,
			expectedAmt:   "1000000000000000000", // 1e18 (1 wei cleared)
			expectedDust:  "1",
			description:   "1 wei is lost in conversion and cleared",
			expectedError: false,
		},
		{
			name:          "maximum dust - almost 1 bridge unit",
			amount:        "1000000999999999999", // 1e18 + (1e12 - 1)
			fromDecimals:  18,
			toDecimals:    6,
			expectedAmt:   "1000000000000000000", // 1e18
			expectedDust:  "999999999999",        // almost 1e12
			description:   "maximum possible dust for 18->6 conversion",
			expectedError: false,
		},
		{
			name:          "large dust amount",
			amount:        "1234567890123456789",
			fromDecimals:  18,
			toDecimals:    6,
			expectedAmt:   "1234567000000000000",
			expectedDust:  "890123456789",
			description:   "larger amount with significant dust",
			expectedError: false,
		},
		{
			name:          "scaling up - no precision loss",
			amount:        "1000000", // 1e6
			fromDecimals:  6,
			toDecimals:    18,
			expectedAmt:   "1000000", // 1e6 (unchanged)
			expectedDust:  "0",
			description:   "scaling up has no precision loss, amount returned as-is",
			expectedError: false,
		},
		{
			name:          "same decimals - no change",
			amount:        "1234567890",
			fromDecimals:  18,
			toDecimals:    18,
			expectedAmt:   "1234567890",
			expectedDust:  "0",
			description:   "no conversion needed when decimals are equal",
			expectedError: false,
		},
		{
			name:          "zero amount",
			amount:        "0",
			fromDecimals:  18,
			toDecimals:    6,
			expectedAmt:   "0",
			expectedDust:  "0",
			description:   "zero amount has no dust",
			expectedError: false,
		},
		{
			name:          "different decimal pairs - 18 to 8",
			amount:        "1000000000000000001", // 1e18 + 1
			fromDecimals:  18,
			toDecimals:    8,
			expectedAmt:   "1000000000000000000", // 1e18
			expectedDust:  "1",
			description:   "18->8 decimals conversion (BTC-like)",
			expectedError: false,
		},
		{
			name:          "small amount all dust",
			amount:        "999", // less than 1e12, so entire amount is dust
			fromDecimals:  18,
			toDecimals:    6,
			expectedAmt:   "0",
			expectedDust:  "999",
			description:   "amount smaller than precision gets cleared to zero",
			expectedError: false,
		},
		{
			name:          "negative amount - error",
			amount:        "-1000000000000000000",
			fromDecimals:  18,
			toDecimals:    6,
			expectedAmt:   "",
			expectedDust:  "",
			description:   "negative amounts should error",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount, ok := sdk.NewIntFromString(tt.amount)
			require.True(t, ok, "failed to parse amount")

			result, err := types.ClearPrecisionLoss(amount, tt.fromDecimals, tt.toDecimals)

			if tt.expectedError {
				require.Error(t, err, tt.description)
			} else {
				require.NoError(t, err, tt.description)

				expected, ok := sdk.NewIntFromString(tt.expectedAmt)
				require.True(t, ok, "failed to parse expected amount")
				require.Equal(t, expected, result, tt.description)

				// Verify dust calculation
				expectedDust, ok := sdk.NewIntFromString(tt.expectedDust)
				require.True(t, ok, "failed to parse expected dust")

				actualDust := amount.Sub(result)
				require.True(t, expectedDust.Equal(actualDust), "dust calculation mismatch: expected %s, got %s - %s",
					expectedDust.String(), actualDust.String(), tt.description)
			}
		})
	}
}

func TestClearPrecisionLoss_RoundTrip(t *testing.T) {
	// Test that ClearPrecisionLoss produces the same result as a manual round-trip conversion
	testCases := []struct {
		amount       string
		fromDecimals uint32
		toDecimals   uint32
	}{
		{"1000000000000000000", 18, 6},
		{"1000000000000000001", 18, 6},
		{"123456789012345678", 18, 6},
		{"1000000000000000000", 18, 8},
		{"999999999999", 18, 6},
	}

	for _, tc := range testCases {
		t.Run(tc.amount, func(t *testing.T) {
			amount, ok := sdk.NewIntFromString(tc.amount)
			require.True(t, ok)

			// Method 1: Use ClearPrecisionLoss
			cleared, err := types.ClearPrecisionLoss(amount, tc.fromDecimals, tc.toDecimals)
			require.NoError(t, err)

			// Method 2: Manual round-trip
			down, err := types.ConvertAmount(amount, tc.fromDecimals, tc.toDecimals)
			require.NoError(t, err)
			up, err := types.ConvertAmount(down, tc.toDecimals, tc.fromDecimals)
			require.NoError(t, err)

			// Both methods should produce the same result
			require.Equal(t, up, cleared, "ClearPrecisionLoss should match manual round-trip")
		})
	}
}
