package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := genState.Validate(); err != nil {
		panic(err)
	}
	k.SetParams(ctx, genState.Params)

	for _, uigids := range genState.UsageIdentifierToGastankIds {
		if err := k.SetUsageIdentifierToGasTankIds(ctx, uigids); err != nil {
			panic(err)
		}
	}

	k.SetLastGasTankID(ctx, genState.LastGasTankId)

	for _, tank := range genState.GasTanks {
		k.SetGasTank(ctx, tank)
	}

	for _, consumer := range genState.GasConsumers {
		k.SetGasConsumer(ctx, consumer)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	gasTankIds, err := k.GetAllUsageIdentifierToGasTankIds(ctx)
	if err != nil {
		panic(err)
	}
	return &types.GenesisState{
		Params:                      k.GetParams(ctx),
		UsageIdentifierToGastankIds: gasTankIds,
		LastGasTankId:               k.GetLastGasTankID(ctx),
		GasTanks:                    k.GetAllGasTanks(ctx),
		GasConsumers:                k.GetAllGasConsumers(ctx),
	}
}
