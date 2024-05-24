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

	/*
		TODO: there used to be a check here to see if the balance which will later be needed for sending
		to the hub is available.
		But I think it's better to not check, potentially..
		Simply, send the tokens!
	*/

	k.SetState(ctx, genState.State)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)
	genesis.State = k.GetState(ctx)

	return genesis
}
