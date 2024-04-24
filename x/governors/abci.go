package staking

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/governors/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	k.BlockGovernorUpdates(ctx)
	return []abci.ValidatorUpdate{}
}
