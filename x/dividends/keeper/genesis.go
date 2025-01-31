package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return nil
}
