package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// InitGenesis new hub-genesis genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	k.SetParams(ctx, genState.Params)
	k.SetState(ctx, genState.State)
}

func (k Keeper) mintCoins(ctx sdk.Context) {
	state := k.GetState(ctx)
	for _, ga := range state.GetGenesisAccounts() {
		coin := ga.GetAmount()
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{coin})
		if err != nil {
			// TODO: okay to panic?
			panic(fmt.Errorf("init genesis mint coins: %w", err))
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)
	genesis.State = k.GetState(ctx)

	return genesis
}
