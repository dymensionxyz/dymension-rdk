package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"
)

func TestIntersectMinGasPrices(t *testing.T) {
	tests := []struct {
		name     string
		a        sdk.DecCoins
		b        sdk.DecCoins
		expected sdk.DecCoins
	}{
		{
			name:     "non-overlapping denoms",
			a:        sdk.DecCoins{sdk.NewDecCoinFromDec("denom1", sdk.NewDec(1))},
			b:        sdk.DecCoins{sdk.NewDecCoinFromDec("denom2", sdk.NewDec(2))},
			expected: sdk.DecCoins{},
		},
		{
			name:     "overlapping denoms with different amounts",
			a:        sdk.DecCoins{sdk.NewDecCoinFromDec("denom1", sdk.NewDec(1))},
			b:        sdk.DecCoins{sdk.NewDecCoinFromDec("denom1", sdk.NewDec(2))},
			expected: sdk.DecCoins{sdk.NewDecCoinFromDec("denom1", sdk.NewDec(2))},
		},
		{
			name:     "overlapping denoms with same amounts",
			a:        sdk.DecCoins{sdk.NewDecCoinFromDec("denom1", sdk.NewDec(1))},
			b:        sdk.DecCoins{sdk.NewDecCoinFromDec("denom1", sdk.NewDec(1))},
			expected: sdk.DecCoins{sdk.NewDecCoinFromDec("denom1", sdk.NewDec(1))},
		},
		{
			name:     "multiple overlapping and non-overlapping denoms",
			a:        sdk.DecCoins{sdk.NewDecCoinFromDec("denom1", sdk.NewDec(1)), sdk.NewDecCoinFromDec("denom2", sdk.NewDec(2))},
			b:        sdk.DecCoins{sdk.NewDecCoinFromDec("denom1", sdk.NewDec(2)), sdk.NewDecCoinFromDec("denom3", sdk.NewDec(3))},
			expected: sdk.DecCoins{sdk.NewDecCoinFromDec("denom1", sdk.NewDec(2))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1 := types.IntersectMinGasPrices(tt.a, tt.b)
			require.Equal(t, tt.expected, result1)

			result2 := types.IntersectMinGasPrices(tt.b, tt.a)
			require.Equal(t, tt.expected, result2)
		})
	}
}
