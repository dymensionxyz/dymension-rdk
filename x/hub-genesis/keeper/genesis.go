package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// InitGenesis new hub-genesis genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	/*
		TODO: there used to be this check to create the module account, I think I will need it now that I don't include
		it in the bank genesis state or anything
		if !k.accountKeeper.HasAccount(ctx, modAddress) {
		}
	*/
	acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	_ = acc
	// TODO: how to mint tokens?

	/*
		TODO: there used to be a check here to see if the balance which will later be needed for sending
		to the hub is available.
		But I think it's better to not check, potentially..
		Simply, send the tokens!
	*/
	for _, ga := range genState.State.GetGenesisAccounts() {
		coin := ga.GetAmount()
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{coin})
		if err != nil {
			panic(fmt.Errorf("init genesis mint coins: %w", err))
		}
	}

	k.SetState(ctx, genState.State)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)
	genesis.State = k.GetState(ctx)

	return genesis
}
