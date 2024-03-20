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

	// Required code as the cosmos sdk validates that the InitChain request is equal to the result
	// so we need to retun here the same valUpdates as we received from the InitChain request
	// reminder: dymint passes two objects, one with the operator address and one with the consensus pubkey
	// the operator address object needs to be removed as
	sequencers := k.GetAllSequencers(ctx)
	if len(sequencers) != 2 {
		panic(types.ErrFailedInitGenesis)
	}

	for _, seq := range sequencers {
		pubkey, err := seq.TmConsPublicKey()
		if err != nil {
			panic(err)
		}

		updateConsPubkey := abci.ValidatorUpdate{
			PubKey: pubkey,
			Power:  seq.ConsensusPower(sdk.DefaultPowerReduction),
		}
		updates = append(updates, updateConsPubkey)
	}

	// delete the genesis sequencer, which we hackly used to keep the data from the InitChain request
	// we stored it only to have the operatorPubKey available to return it in the ValidatorUpdate
	val, ok := k.GetSequencer(ctx, sdk.ValAddress(types.InitChainStubAddr))
	if !ok {
		panic("genesis sequencer not found")
	}
	k.DeleteSequencer(ctx, val)

	return updates
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	return genesis
}
