package coinsutils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetCoinByDenomFromCoins(denom string, coins sdk.Coins) (sdk.Coin, bool) {
	for _, coin := range coins {
		if coin.Denom == denom {
			return coin, true
		}
	}
	return sdk.Coin{}, false
}
