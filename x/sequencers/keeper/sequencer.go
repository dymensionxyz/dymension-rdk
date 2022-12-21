package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
)

// GetSequencer returns a sequencer from its index
func (k Keeper) GetSequencer(
	ctx sdk.Context,
	sequencerAddress string,

) (val types.Sequencer, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ValidatorsKey))

	b := store.Get(types.SequencerKey(
		sequencerAddress,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllSequencer returns all sequencer
func (k Keeper) GetAllSequencer(ctx sdk.Context) (list []stakingtypes.Validator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(string(types.SequencerKey())))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	// nolint: errcheck
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Sequencer
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
