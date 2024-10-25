package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

// InitGenesis new hub genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	for _, denom := range genState.State.Hub.RegisteredDenoms {
		if err := k.SetHubDenom(ctx, denom.Base); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	denoms, err := k.GetAllHubDenoms(ctx)
	if err != nil {
		panic(err)
	}

	for _, denom := range denoms {
		genesis.State.Hub.RegisteredDenoms = append(genesis.State.Hub.RegisteredDenoms, &types.RegisteredDenom{
			Base: denom,
		})
	}
	return genesis
}
