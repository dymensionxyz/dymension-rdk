package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
)

/* -------------------------------------------------------------------------- */
/*                                    Alias func                              */
/* -------------------------------------------------------------------------- */
// Validator gets the Validator interface for a particular address
func (k Keeper) Validator(ctx sdk.Context, address sdk.ValAddress) stakingtypes.ValidatorI {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return nil
	}

	return val
}

// ValidatorByConsAddr gets the validator interface for a particular pubkey
func (k Keeper) ValidatorByConsAddr(ctx sdk.Context, addr sdk.ConsAddress) stakingtypes.ValidatorI {
	val, found := k.GetValidatorByConsAddr(ctx, addr)
	if !found {
		return nil
	}

	return val
}

/* -------------------------------------------------------------------------- */
/*                               implementation                              */
/* -------------------------------------------------------------------------- */

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
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)

	opAddr := store.Get(types.GetValidatorByConsAddrKey(consAddr))
	if opAddr == nil {
		return validator, false
	}

	return k.GetValidator(ctx, opAddr)
}

// set the main record holding validator details
func (k Keeper) SetValidator(ctx sdk.Context, validator stakingtypes.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := stakingtypes.MustMarshalValidator(k.cdc, &validator)
	store.Set(types.GetValidatorKey(validator.GetOperator()), bz)
}

// validator index
func (k Keeper) SetValidatorByConsAddr(ctx sdk.Context, validator stakingtypes.Validator) error {
	consPk, err := validator.GetConsAddr()
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorByConsAddrKey(consPk), validator.GetOperator())
	return nil
}

// get groups of validators

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

// // GetSequencer returns a sequencer from its index
// func (k Keeper) GetSequencer(ctx sdk.Context, sequencerAddress string) (val stakingtypes.Validator, found bool) {
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ValidatorsKey))

// 	// k.paramstore.Get(ctx, types.KeyHistoricalEntries, &res)

// 	store := ctx.KVStore(k.storeKey)
// 	validators = make([]types.Validator, maxRetrieve)

// 	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
// 	defer iterator.Close()

// 	value := store.Get(key)
// 	if value == nil {
// 		return stakingtypes.HistoricalInfo{}, false
// 	}

// 	*/

// 	b := store.Get(types.val(
// 		sequencerAddress,
// 	))
// 	if b == nil {
// 		return val, false
// 	}

// 	k.cdc.MustUnmarshal(b, &val)
// 	return val, true
// }

// // GetAllSequencer returns all sequencer
// func (k Keeper) GetAllSequencer(ctx sdk.Context) (list []stakingtypes.Validator) {
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetValidatorKey())
// 	iterator := sdk.KVStorePrefixIterator(store, []byte{})

// 	// nolint: errcheck
// 	defer iterator.Close()

// 	for ; iterator.Valid(); iterator.Next() {
// 		var val types.Sequencer
// 		k.cdc.MustUnmarshal(iterator.Value(), &val)
// 		list = append(list, val)
// 	}

// 	return
// }

// // unmarshal a redelegation from a store value
// func MustUnmarshalValidator(cdc codec.BinaryCodec, value []byte) stakingtypes.Validator {
// 	validator, err := UnmarshalValidator(cdc, value)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return validator
// }

// // unmarshal a redelegation from a store value
// func UnmarshalValidator(cdc codec.BinaryCodec, value []byte) (v stakingtypes.Validator, err error) {
// 	err = cdc.Unmarshal(value, &v)
// 	return v, err
// }

// // return the redelegation
// func MustMarshalValidator(cdc codec.BinaryCodec, validator *Validator) []byte {
// 	return cdc.MustMarshal(validator)
// }
