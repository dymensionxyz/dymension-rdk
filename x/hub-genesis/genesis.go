package hub_genesis

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// InitGenesis new hub-genesis genesis.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, ak types.AccountKeeper, genState *types.GenesisState) {
	keeper.SetParams(ctx, genState.Params)

	if !ak.HasAccount(ctx, ak.GetModuleAddress(types.ModuleName)) {
		ak.GetModuleAccount(ctx, types.ModuleName)
	}

	// TODO: check genesis balance is enough for expected tokens!
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = keeper.GetParams(ctx)
	genesis.State = keeper.GetState(ctx)

	return genesis
}
