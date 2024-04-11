package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// InitGenesis new hub-genesis genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	modAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	if !k.accountKeeper.HasAccount(ctx, modAddress) {
		k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	}

	// if tokens provided and waiting to be locked, verify the balance
	if !genState.State.IsLocked && !genState.State.GenesisTokens.IsZero() {
		// get spendable coins in the module account
		spendable := k.bankKeeper.SpendableCoins(ctx, modAddress)
		// we expect the genesis balance of the module account to be equal to required genesis tokens
		if !spendable.IsEqual(genState.State.GenesisTokens) {
			panic(types.ErrWrongGenesisBalance)
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
