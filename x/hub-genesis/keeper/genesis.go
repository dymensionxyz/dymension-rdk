package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// InitGenesis new hub-genesis genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	/*
		Mint coins which will later be transferred to the hub
		TODO: need to not do it if it not height 0 genesis

		TODO: move to foo
	*/
	for _, ga := range genState.State.GetGenesisAccounts() {
		coin := ga.GetAmount()
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{coin})
		if err != nil {
			// TODO: okay to panic?
			panic(fmt.Errorf("init genesis mint coins: %w", err))
		}
	}

	k.SetParams(ctx, genState.Params)
	k.SetState(ctx, genState.State)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)
	genesis.State = k.GetState(ctx)

	return genesis
}
