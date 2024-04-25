package denommetadata

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *types.GenesisState {
	return &types.GenesisState{}
}

// ExportGenesis export module state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{}
}
