package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// InitGenesis initializes the sequencers module's state from a provided genesis state.
// We return the ValidatorUpdate set by init chain
func (k *Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) []abci.ValidatorUpdate {
	k.SetParams(ctx, genState.Params)

	for _, s := range genState.GetSequencers() {
		k.SetSequencer(ctx, *s.Validator)
		if s.RewardAddr != "" {
			k.SetRewardAddr(ctx, *s.Validator, sdk.AccAddress(s.RewardAddr))
		}
	}

	// return (and delete) the update from init chain
	updates := make([]abci.ValidatorUpdate, 1)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ValidatorUpdateKey)
	k.cdc.MustUnmarshal(bz, &updates[0])
	store.Delete(types.ValidatorUpdateKey)
	return updates
}

// ExportGenesis returns the sequencers module's exported genesis.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	sequencersAsValidators := k.GetAllSequencers(ctx)
	genesis.Sequencers = make([]types.Sequencer, len(sequencersAsValidators))
	for i, v := range sequencersAsValidators {
		genesis.Sequencers[i].Validator = &v
		rewardAddr, ok := k.GetRewardAddr(ctx, v.GetOperator())
		if ok {
			genesis.Sequencers[i].RewardAddr = rewardAddr.String()
		}
	}

	return genesis
}

// MustSetDymintValidatorUpdates  - ABCI expects the result of init genesis to return the same value as passed in InitChainer,
// so we save it to return later.
func (k Keeper) MustSetDymintValidatorUpdates(ctx sdk.Context, updates []abci.ValidatorUpdate) {
	if len(updates) != 1 {
		panic(errors.Wrapf(gerrc.ErrOutOfRange, "expect 1 abci validator update: got: %d", len(updates)))
	}
	u := updates[0]
	k.cdc.MustMarshal(&u)
	ctx.KVStore(k.storeKey).Set(types.ValidatorUpdateKey, k.cdc.MustMarshal(&u))
}
