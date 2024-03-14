package sequencers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// InitGenesis initializes the capability module's state from a provided genesis state.
// We return the for ValidatorUpdate only the sequencers set by dymint
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) []abci.ValidatorUpdate {
	var updates []abci.ValidatorUpdate

	k.SetParams(ctx, genState.Params)

	// Set all the sequencer if exists on genesis file
	for _, elem := range genState.Sequencers {
		pk, _ := elem.ConsPubKey()
		if _, err := k.CreateSequencer(ctx, elem.OperatorAddress, pk); err != nil {
			panic(err)
		}
		updates = append(updates, elem.ABCIValidatorUpdate(sdk.DefaultPowerReduction))
	}

	//get the dymint sequencers if it's clean genesis
	if len(updates) == 0 {
		val, ok := k.GetValidator(ctx, sdk.ValAddress(types.GenesisOperatorAddrStub))
		if ok {
			updates = append(updates, val.ABCIValidatorUpdate(sdk.DefaultPowerReduction))
		}
	}

	return updates
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.Sequencers = k.GetAllValidators(ctx)

	return genesis
}
