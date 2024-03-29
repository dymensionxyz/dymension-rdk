package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MintHooks interface {
	AfterDistributeMintedCoin(ctx sdk.Context, mintedCoins sdk.Coins)
}

var _ MintHooks = MultiMintHooks{}

// combine multiple mint hooks, all hook functions are run in array sequence.
type MultiMintHooks []MintHooks

func NewMultiMintHooks(hooks ...MintHooks) MultiMintHooks {
	return hooks
}

func (h MultiMintHooks) AfterDistributeMintedCoin(ctx sdk.Context, mintedCoins sdk.Coins) {
	for i := range h {
		h[i].AfterDistributeMintedCoin(ctx, mintedCoins)
	}
}
