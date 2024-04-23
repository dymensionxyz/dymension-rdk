package vesting

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/x/vesting/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/vesting/types"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *types.GenesisState {
	return &types.GenesisState{
		Params: types.DefaultParams(),
	}
}

// InitGenesis import module genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	data types.GenesisState,
) {
	k.SetParams(ctx, data.Params)
}

// ExportGenesis export module state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
	}
}

func ValidateGenesis(gs types.GenesisState) error {
	return gs.Params.Validate()
}
