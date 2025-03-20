package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
	"github.com/stretchr/testify/require"
)

func TestFilterDenoms(t *testing.T) {
	tests := []struct {
		name   string
		coins  sdk.Coins
		denoms []string
		result sdk.Coins
	}{
		{
			name:   "all denoms match",
			coins:  sdk.NewCoins(sdk.NewInt64Coin("adym", 10), sdk.NewInt64Coin("uosmo", 10)),
			denoms: []string{"adym", "uosmo"},
			result: sdk.NewCoins(sdk.NewInt64Coin("adym", 10), sdk.NewInt64Coin("uosmo", 10)),
		},
		{
			name:   "one denom matches",
			coins:  sdk.NewCoins(sdk.NewInt64Coin("adym", 10), sdk.NewInt64Coin("uosmo", 10)),
			denoms: []string{"uosmo"},
			result: sdk.NewCoins(sdk.NewInt64Coin("uosmo", 10)),
		},
		{
			name:   "no denoms match",
			coins:  sdk.NewCoins(sdk.NewInt64Coin("adym", 10), sdk.NewInt64Coin("uosmo", 10)),
			denoms: []string{},
			result: sdk.Coins{},
		},
		{
			name:   "single coin matches",
			coins:  sdk.NewCoins(sdk.NewInt64Coin("adym", 10)),
			denoms: []string{"adym", "uosmo"},
			result: sdk.NewCoins(sdk.NewInt64Coin("adym", 10)),
		},
		{
			name:   "no coins",
			coins:  sdk.Coins{},
			denoms: []string{"adym", "uosmo"},
			result: sdk.Coins{},
		},
		{
			name:   "no coins and no denoms",
			coins:  sdk.Coins{},
			denoms: []string{},
			result: sdk.Coins{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := types.FilterDenoms(tt.coins, tt.denoms)
			require.Equal(t, tt.result, filtered)
		})
	}
}
