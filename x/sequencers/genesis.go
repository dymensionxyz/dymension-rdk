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
	sequencers := k.GetAllSequencers(ctx)
	if len(sequencers) > 2 {
		panic(types.ErrMultipleDymintSequencers)
	}
	if len(sequencers) == 0 {
		panic(types.ErrNoSequencerOnInitChain)
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

	val, ok := k.GetSequencer(ctx, sdk.ValAddress(types.GenesisOperatorAddrStub))
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
