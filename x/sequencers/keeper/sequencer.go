package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// GetValidatorByConsAddr get a single validator by consensus address - used for interface compat
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, found bool) {
	return k.GetSequencerByConsAddr(ctx, consAddr)
}

// GetSequencer get a single sequencer
func (k Keeper) GetSequencer(ctx sdk.Context, operator sdk.ValAddress) (sequencer stakingtypes.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)

	value := store.Get(types.GetSequencerKey(operator))
	if value == nil {
		return sequencer, false
	}

	sequencer = stakingtypes.MustUnmarshalValidator(k.cdc, value)
	return sequencer, true
}

// GetSequencerByConsAddr get a single sequencer by consensus address
func (k Keeper) GetSequencerByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (sequencer stakingtypes.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	opAddr := store.Get(types.GetSequencerByConsAddrKey(consAddr))
	if opAddr == nil {
		return sequencer, false
	}

	return k.GetSequencer(ctx, opAddr)
}

// SetSequencer set the main record holding sequencer details
func (k Keeper) SetSequencer(ctx sdk.Context, sequencer stakingtypes.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := stakingtypes.MustMarshalValidator(k.cdc, &sequencer)
	store.Set(types.GetSequencerKey(sequencer.GetOperator()), bz)
	k.MustSetSequencerByConsAddr(ctx, sequencer)
}

// DeleteSequencer delete the main record holding sequencer details
func (k Keeper) DeleteSequencer(ctx sdk.Context, sequencer stakingtypes.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetSequencerKey(sequencer.GetOperator()))
}

func (k Keeper) MustSetSequencerByConsAddr(ctx sdk.Context, sequencer stakingtypes.Validator) {
	consAddr, err := sequencer.GetConsAddr()
	if err != nil {
		panic(err)
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetSequencerByConsAddrKey(consAddr), sequencer.GetOperator())
}

// GetAllSequencers get the set of all sequencers with no limits, used during genesis dump
func (k Keeper) GetAllSequencers(ctx sdk.Context) (sequencers []stakingtypes.Validator) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.SequencersKey)
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		sequencer := stakingtypes.MustUnmarshalValidator(k.cdc, iterator.Value())
		sequencers = append(sequencers, sequencer)
	}

	return sequencers
}

// SetRewardAddr sets the address that rewards will be allocated to
func (k Keeper) SetRewardAddr(ctx sdk.Context, sequencer stakingtypes.Validator, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetSequencerRewardAddrKey(sequencer.GetOperator()), addr)
}

// GetRewardAddr gets the addr that the sequencer wants to allocate rewards to
func (k Keeper) GetRewardAddr(ctx sdk.Context, operator sdk.ValAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	res := store.Get(types.GetSequencerRewardAddrKey(operator))
	return res, res != nil
}

func (k Keeper) GetRewardAddrByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (sdk.AccAddress, bool) {
	seq, ok := k.GetSequencerByConsAddr(ctx, consAddr)
	if !ok {
		return sdk.AccAddress{}, false
	}
	return k.GetRewardAddr(ctx, seq.GetOperator())
}
