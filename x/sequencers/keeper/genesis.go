package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// InitGenesis initializes the sequencers module's state from a provided genesis state.
// We return the ValidatorUpdate set by init chain
func (k *Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, genState.Params)

	updates := make([]abci.ValidatorUpdate, 1)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ValidatorUpdateKey)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &updates[0])
	}
	store.Delete(types.ValidatorUpdateKey)
	return updates
}

// ExportGenesis returns the sequencers module's exported genesis.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	sequencers := k.GetAllSequencers(ctx)
	// TODO:
	_ = sequencers

	return genesis
}
