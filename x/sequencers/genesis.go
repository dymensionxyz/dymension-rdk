package sequencers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/rollapp/x/sequencers/keeper"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, genState.Params)

	valUpdates := []abci.ValidatorUpdate{}

	// Set all the sequencer
	for _, elem := range genState.Sequencers {
		if elem.OperatorAddress == "" {
			if err := k.SetDymintSequencerByAddr(ctx, elem); err != nil {
				panic(err)
			}
		} else {
			//TODO: refactor to call create sequecncer from keeper
			consAddr, _ := elem.GetConsAddr()
			power, found := k.GetDymintSequencerByAddr(ctx, sdk.ConsAddress(consAddr))
			if !found {
				panic("trying to register unknown sequencer")
			}

			pk, _ := elem.ConsPubKey()
			seq, err := types.NewSequencer(elem.GetOperator(), pk, uint64(power))
			if err != nil {
				panic(err)
			}

			k.SetValidator(ctx, seq)
			if err := k.SetValidatorByConsAddr(ctx, seq); err != nil {
				panic(sdkerrors.Wrapf(err, "failed to InitGenesis for sequencers"))
			}
			valUp := seq.ABCIValidatorUpdate(sdk.DefaultPowerReduction)
			valUpdates = append(valUpdates, valUp)
		}
	}

	if len(valUpdates) == 0 {
		panic("no sequencer registered on genesis")
	}

	return valUpdates
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.Sequencers = k.GetAllValidators(ctx)

	return genesis
}
