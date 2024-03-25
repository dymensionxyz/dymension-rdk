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

	// TODO: move to validateGenesis
	if genState.GenesisOperatorAddress == "" {
		panic("genesis operator address not set")
	}

	operatorAddr, err := sdk.ValAddressFromBech32(genState.GenesisOperatorAddress)
	if err != nil {
		panic(err)
	}

	// get the sequencer we set on InitChain. and delete it as it's not needed in store anymore
	seq, ok := k.GetSequencer(ctx, sdk.ValAddress(types.InitChainStubAddr))
	if !ok {
		panic("genesis sequencer not found")
	}
	k.DeleteSequencer(ctx, seq)

	pubkey, err := seq.ConsPubKey()
	if err != nil {
		panic(err)
	}
	power := seq.ConsensusPower(sdk.DefaultPowerReduction)

	sequencer, err := types.NewSequencer(operatorAddr, pubkey, power)
	if err != nil {
		panic(err)
	}

	k.SetSequencer(ctx, sequencer)
	err = k.SetSequencerByConsAddr(ctx, sequencer)
	if err != nil {
		panic(err)
	}

	tmPubkey, err := seq.TmConsPublicKey()
	if err != nil {
		panic(err)
	}
	updateConsPubkey := abci.ValidatorUpdate{
		PubKey: tmPubkey,
		Power:  power,
	}
	updates = append(updates, updateConsPubkey)

	return updates
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	return genesis
}
