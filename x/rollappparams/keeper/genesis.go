package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	if err := types.ValidateGenesis(genState); err != nil {
		panic(err)
	}

	if err := k.SetDA(ctx, genState.Params.Da); err != nil {
		panic(err)
	}
	if err := k.SetVersion(ctx, genState.Params.DrsVersion); err != nil {
		panic(err)
	}
	if err := k.SetMinGasPrices(ctx, genState.Params.MinGasPrices); err != nil {
		panic(err)
	}
	k.SetFreeIBC(ctx, genState.Params.FreeIbc)
}

// ExportGenesis returns the capability module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
	}
}
