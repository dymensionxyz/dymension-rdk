package sequencers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/rollapp/x/sequencers/keeper"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, genState.Params)

	var updates []abci.ValidatorUpdate

	// Set all the sequencer
	for _, elem := range genState.Sequencers {
		if elem.OperatorAddress == "" {
			if err := k.SetDymintSequencerByAddr(ctx, elem); err != nil {
				panic(err)
			}
		} else {
			pk, _ := elem.ConsPubKey()
			if _, err := k.CreateSequencer(ctx, elem.OperatorAddress, pk); err != nil {
				panic(err)
			}
		}

		updates = append(updates, elem.ABCIValidatorUpdate(sdk.DefaultPowerReduction))
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
