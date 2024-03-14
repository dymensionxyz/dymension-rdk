package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

/* ---------------------------------- alias --------------------------------- */

// get a single validator by consensus address
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.ValidatorI, found bool) {
	return k.GetSequencerByConsAddr(ctx, consAddr)
}

/* --------------------------------- GETTERS -------------------------------- */
// get a single validator
func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)

	value := store.Get(types.GetValidatorKey(addr))
	if value == nil {
		return validator, false
	}

	validator = stakingtypes.MustUnmarshalValidator(k.cdc, value)
	return validator, true
}

// get a single validator by consensus address
func (k Keeper) GetSequencerByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (sequencer stakingtypes.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	opAddr := store.Get(types.GetValidatorByConsAddrKey(consAddr))
	if opAddr == nil {
		return sequencer, false
	}

	return k.GetValidator(ctx, opAddr)
}

/* --------------------------------- SETTERS -------------------------------- */
// set the main record holding validator details
func (k Keeper) SetValidator(ctx sdk.Context, validator stakingtypes.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := stakingtypes.MustMarshalValidator(k.cdc, &validator)
	store.Set(types.GetValidatorKey(validator.GetOperator()), bz)
}

// delete the main record holding validator details
func (k Keeper) DeletetValidator(ctx sdk.Context, validator stakingtypes.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorKey(validator.GetOperator()))
}

// validator index
func (k Keeper) SetValidatorByConsAddr(ctx sdk.Context, validator stakingtypes.Validator) error {
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorByConsAddrKey(consAddr), validator.GetOperator())

	return nil
}

// get the set of all validators with no limits, used during genesis dump
func (k Keeper) GetAllValidators(ctx sdk.Context) (validators []stakingtypes.Validator) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := stakingtypes.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}

	return validators
}

func (k Keeper) SetOperatorAddressForGenesisSequencer(ctx sdk.Context, proposerValAddr sdk.ValAddress) {
	val, ok := k.GetValidator(ctx, sdk.ValAddress(types.GenesisOperatorAddrStub))
	if !ok {
		return
	}

	k.DeletetValidator(ctx, val)
	val.OperatorAddress = proposerValAddr.String()
	k.SetValidator(ctx, val)
	k.SetValidatorByConsAddr(ctx, val)
}
