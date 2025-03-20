package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// FilterDenoms filters out coins with denoms not in the provided list
func FilterDenoms(coins sdk.Coins, denoms []string) sdk.Coins {
	d := make(map[string]struct{}, len(denoms))
	for _, denom := range denoms {
		d[denom] = struct{}{}
	}
	filtered := make(sdk.Coins, 0, len(denoms))
	for _, coin := range coins {
		if _, ok := d[coin.Denom]; ok {
			filtered = filtered.Add(coin)
		}
	}
	return filtered
}
