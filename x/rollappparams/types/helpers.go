package types

import (
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IntersectMinGasPrices returns the intersection of two MinGasPrices.
// The intersection is the maximum of the two prices for each denom.
// Operation is commutative.
func IntersectMinGasPrices(a, b sdk.DecCoins) sdk.DecCoins {
	res := make([]sdk.DecCoin, 0, len(a))
	for _, coin := range a {
		bAmount := b.AmountOf(coin.Denom)
		if bAmount.IsZero() {
			continue
		}

		res = append(res, sdk.DecCoin{
			Denom:  coin.Denom,
			Amount: sdk.MaxDec(coin.Amount, bAmount),
		})
	}
	return slices.Clip(res)
}
