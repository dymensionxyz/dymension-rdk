package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/dymensionxyz/dymension-rdk/x/convertor/types"
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
