package keeper

import (
	"bytes"
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GetDelegation returns a specific delegation.
func (k Keeper) GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDelegationKey(delAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return delegation, false
	}

	delegation = stakingtypes.MustUnmarshalDelegation(k.cdc, value)

	return delegation, true
}

// IterateAllDelegations iterates through all of the delegations.
func (k Keeper) IterateAllDelegations(ctx sdk.Context, cb func(delegation stakingtypes.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, stakingtypes.DelegationKey)
	defer iterator.Close() //nolint:errcheck

	for ; iterator.Valid(); iterator.Next() {
		delegation := stakingtypes.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k Keeper) GetAllDelegations(ctx sdk.Context) (delegations []stakingtypes.Delegation) {
	k.IterateAllDelegations(ctx, func(delegation stakingtypes.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// GetGovernorDelegations returns all delegations to a specific governor.
// Useful for querier.
func (k Keeper) GetGovernorDelegations(ctx sdk.Context, valAddr sdk.ValAddress) (delegations []stakingtypes.Delegation) { //nolint:interfacer
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, stakingtypes.DelegationKey)
	defer iterator.Close() //nolint:errcheck

	for ; iterator.Valid(); iterator.Next() {
		delegation := stakingtypes.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if delegation.GetValidatorAddr().Equals(valAddr) {
			delegations = append(delegations, delegation)
		}
	}

	return delegations
}

// GetDelegatorDelegations returns a given amount of all the delegations from a
// delegator.
func (k Keeper) GetDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (delegations []stakingtypes.Delegation) {
	delegations = make([]stakingtypes.Delegation, maxRetrieve)
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close() //nolint:errcheck

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		delegation := stakingtypes.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations[i] = delegation
		i++
	}

	return delegations[:i] // trim if the array length < maxRetrieve
}

// SetDelegation sets a delegation.
func (k Keeper) SetDelegation(ctx sdk.Context, delegation stakingtypes.Delegation) {
	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)

	store := ctx.KVStore(k.storeKey)
	b := stakingtypes.MustMarshalDelegation(k.cdc, delegation)
	store.Set(types.GetDelegationKey(delegatorAddress, delegation.GetValidatorAddr()), b)
}

// RemoveDelegation removes a delegation
func (k Keeper) RemoveDelegation(ctx sdk.Context, delegation stakingtypes.Delegation) error {
	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)

	// TODO: Consider calling hooks outside of the store wrapper functions, it's unobvious.
	if err := k.BeforeDelegationRemoved(ctx, delegatorAddress, delegation.GetValidatorAddr()); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegationKey(delegatorAddress, delegation.GetValidatorAddr()))
	return nil
}

// GetUnbondingDelegations returns a given amount of all the delegator unbonding-delegations.
func (k Keeper) GetUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (unbondingDelegations []stakingtypes.UnbondingDelegation) {
	unbondingDelegations = make([]stakingtypes.UnbondingDelegation, maxRetrieve)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetUBDsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close() //nolint:errcheck

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		unbondingDelegation := stakingtypes.MustUnmarshalUBD(k.cdc, iterator.Value())
		unbondingDelegations[i] = unbondingDelegation
		i++
	}

	return unbondingDelegations[:i] // trim if the array length < maxRetrieve
}

// GetUnbondingDelegation returns a unbonding delegation.
func (k Keeper) GetUnbondingDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (ubd stakingtypes.UnbondingDelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(delAddr, valAddr)
	value := store.Get(key)

	if value == nil {
		return ubd, false
	}

	ubd = stakingtypes.MustUnmarshalUBD(k.cdc, value)

	return ubd, true
}

// GetUnbondingDelegationsFromGovernor returns all unbonding delegations from a
// particular governor.
func (k Keeper) GetUnbondingDelegationsFromGovernor(ctx sdk.Context, valAddr sdk.ValAddress) (ubds []stakingtypes.UnbondingDelegation) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetUBDsByValIndexKey(valAddr))
	defer iterator.Close() // nolint:errcheck

	for ; iterator.Valid(); iterator.Next() {
		key := types.GetUBDKeyFromValIndexKey(iterator.Key())
		value := store.Get(key)
		ubd := stakingtypes.MustUnmarshalUBD(k.cdc, value)
		ubds = append(ubds, ubd)
	}

	return ubds
}

// IterateUnbondingDelegations iterates through all of the unbonding delegations.
func (k Keeper) IterateUnbondingDelegations(ctx sdk.Context, fn func(index int64, ubd stakingtypes.UnbondingDelegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.UnbondingDelegationKey)
	defer iterator.Close() //nolint:errcheck

	for i := int64(0); iterator.Valid(); iterator.Next() {
		ubd := stakingtypes.MustUnmarshalUBD(k.cdc, iterator.Value())
		if stop := fn(i, ubd); stop {
			break
		}
		i++
	}
}

// GetDelegatorUnbonding returns the total amount a delegator has unbonding.
func (k Keeper) GetDelegatorUnbonding(ctx sdk.Context, delegator sdk.AccAddress) math.Int {
	unbonding := sdk.ZeroInt()
	k.IterateDelegatorUnbondingDelegations(ctx, delegator, func(ubd stakingtypes.UnbondingDelegation) bool {
		for _, entry := range ubd.Entries {
			unbonding = unbonding.Add(entry.Balance)
		}
		return false
	})
	return unbonding
}

// IterateDelegatorUnbondingDelegations iterates through a delegator's unbonding delegations.
func (k Keeper) IterateDelegatorUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress, cb func(ubd stakingtypes.UnbondingDelegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetUBDsKey(delegator))
	defer iterator.Close() //nolint:errcheck

	for ; iterator.Valid(); iterator.Next() {
		ubd := stakingtypes.MustUnmarshalUBD(k.cdc, iterator.Value())
		if cb(ubd) {
			break
		}
	}
}

// GetDelegatorBonded returs the total amount a delegator has bonded.
func (k Keeper) GetDelegatorBonded(ctx sdk.Context, delegator sdk.AccAddress) math.Int {
	bonded := sdk.ZeroDec()

	k.IterateDelegatorDelegations(ctx, delegator, func(delegation stakingtypes.Delegation) bool {
		governorAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err) // shouldn't happen
		}
		governor, found := k.GetGovernor(ctx, governorAddr)
		if found {
			shares := delegation.Shares
			tokens := governor.TokensFromSharesTruncated(shares)
			bonded = bonded.Add(tokens)
		}
		return false
	})
	return bonded.RoundInt()
}

// IterateDelegatorDelegations iterates through one delegator's delegations.
func (k Keeper) IterateDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, cb func(delegation stakingtypes.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close() //nolint:errcheck
	for ; iterator.Valid(); iterator.Next() {
		delegation := stakingtypes.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// IterateDelegatorRedelegations iterates through one delegator's redelegations.
func (k Keeper) IterateDelegatorRedelegations(ctx sdk.Context, delegator sdk.AccAddress, cb func(red stakingtypes.Redelegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetREDsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close() //nolint:errcheck
	for ; iterator.Valid(); iterator.Next() {
		red := stakingtypes.MustUnmarshalRED(k.cdc, iterator.Value())
		if cb(red) {
			break
		}
	}
}

// HasMaxUnbondingDelegationEntries - check if unbonding delegation has maximum number of entries.
func (k Keeper) HasMaxUnbondingDelegationEntries(ctx sdk.Context, delegatorAddr sdk.AccAddress, governorAddr sdk.ValAddress) bool {
	ubd, found := k.GetUnbondingDelegation(ctx, delegatorAddr, governorAddr)
	if !found {
		return false
	}

	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
}

// SetUnbondingDelegation sets the unbonding delegation and associated index.
func (k Keeper) SetUnbondingDelegation(ctx sdk.Context, ubd stakingtypes.UnbondingDelegation) {
	delegatorAddress := sdk.MustAccAddressFromBech32(ubd.DelegatorAddress)

	store := ctx.KVStore(k.storeKey)
	bz := stakingtypes.MustMarshalUBD(k.cdc, ubd)
	addr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	key := types.GetUBDKey(delegatorAddress, addr)
	store.Set(key, bz)
	store.Set(types.GetUBDByValIndexKey(delegatorAddress, addr), []byte{}) // index, store empty bytes
}

// RemoveUnbondingDelegation removes the unbonding delegation object and associated index.
func (k Keeper) RemoveUnbondingDelegation(ctx sdk.Context, ubd stakingtypes.UnbondingDelegation) {
	delegatorAddress := sdk.MustAccAddressFromBech32(ubd.DelegatorAddress)

	store := ctx.KVStore(k.storeKey)
	addr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	key := types.GetUBDKey(delegatorAddress, addr)
	store.Delete(key)
	store.Delete(types.GetUBDByValIndexKey(delegatorAddress, addr))
}

// SetUnbondingDelegationEntry adds an entry to the unbonding delegation at
// the given addresses. It creates the unbonding delegation if it does not exist.
func (k Keeper) SetUnbondingDelegationEntry(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, governorAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, balance math.Int,
) stakingtypes.UnbondingDelegation {
	ubd, found := k.GetUnbondingDelegation(ctx, delegatorAddr, governorAddr)
	if found {
		ubd.AddEntry(creationHeight, minTime, balance)
	} else {
		ubd = stakingtypes.NewUnbondingDelegation(delegatorAddr, governorAddr, creationHeight, minTime, balance)
	}

	k.SetUnbondingDelegation(ctx, ubd)

	return ubd
}

// unbonding delegation queue timeslice operations

// GetUBDQueueTimeSlice gets a specific unbonding queue timeslice. A timeslice
// is a slice of DVPairs corresponding to unbonding delegations that expire at a
// certain time.
func (k Keeper) GetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvPairs []stakingtypes.DVPair) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetUnbondingDelegationTimeKey(timestamp))
	if bz == nil {
		return []stakingtypes.DVPair{}
	}

	pairs := stakingtypes.DVPairs{}
	k.cdc.MustUnmarshal(bz, &pairs)

	return pairs.Pairs
}

// SetUBDQueueTimeSlice sets a specific unbonding queue timeslice.
func (k Keeper) SetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []stakingtypes.DVPair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stakingtypes.DVPairs{Pairs: keys})
	store.Set(types.GetUnbondingDelegationTimeKey(timestamp), bz)
}

// InsertUBDQueue inserts an unbonding delegation to the appropriate timeslice
// in the unbonding queue.
func (k Keeper) InsertUBDQueue(ctx sdk.Context, ubd stakingtypes.UnbondingDelegation, completionTime time.Time) {
	dvPair := stakingtypes.DVPair{DelegatorAddress: ubd.DelegatorAddress, ValidatorAddress: ubd.ValidatorAddress}

	timeSlice := k.GetUBDQueueTimeSlice(ctx, completionTime)
	if len(timeSlice) == 0 {
		k.SetUBDQueueTimeSlice(ctx, completionTime, []stakingtypes.DVPair{dvPair})
	} else {
		timeSlice = append(timeSlice, dvPair)
		k.SetUBDQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// UBDQueueIterator returns all the unbonding queue timeslices from time 0 until endTime.
func (k Keeper) UBDQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnbondingQueueKey,
		sdk.InclusiveEndBytes(types.GetUnbondingDelegationTimeKey(endTime)))
}

// DequeueAllMatureUBDQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue.
func (k Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context, currTime time.Time) (matureUnbonds []stakingtypes.DVPair) {
	store := ctx.KVStore(k.storeKey)

	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator := k.UBDQueueIterator(ctx, currTime)
	defer unbondingTimesliceIterator.Close() // nolint: errcheck

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := stakingtypes.DVPairs{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		matureUnbonds = append(matureUnbonds, timeslice.Pairs...)

		store.Delete(unbondingTimesliceIterator.Key())
	}

	return matureUnbonds
}

// GetRedelegations returns a given amount of all the delegator redelegations.
func (k Keeper) GetRedelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (redelegations []stakingtypes.Redelegation) {
	redelegations = make([]stakingtypes.Redelegation, maxRetrieve)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetREDsKey(delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close() //nolint:errcheck
	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		redelegation := stakingtypes.MustUnmarshalRED(k.cdc, iterator.Value())
		redelegations[i] = redelegation
		i++
	}

	return redelegations[:i] // trim if the array length < maxRetrieve
}

// GetRedelegation returns a redelegation.
func (k Keeper) GetRedelegation(ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress) (red stakingtypes.Redelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetREDKey(delAddr, valSrcAddr, valDstAddr)

	value := store.Get(key)
	if value == nil {
		return red, false
	}

	red = stakingtypes.MustUnmarshalRED(k.cdc, value)

	return red, true
}

// GetRedelegationsFromSrcGovernor returns all redelegations from a particular
// governor.
func (k Keeper) GetRedelegationsFromSrcGovernor(ctx sdk.Context, valAddr sdk.ValAddress) (reds []stakingtypes.Redelegation) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetREDsFromValSrcIndexKey(valAddr))
	defer iterator.Close() //nolint:errcheck
	for ; iterator.Valid(); iterator.Next() {
		key := types.GetREDKeyFromValSrcIndexKey(iterator.Key())
		value := store.Get(key)
		red := stakingtypes.MustUnmarshalRED(k.cdc, value)
		reds = append(reds, red)
	}

	return reds
}

// HasReceivingRedelegation checks if governor is receiving a redelegation.
func (k Keeper) HasReceivingRedelegation(ctx sdk.Context, delAddr sdk.AccAddress, valDstAddr sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetREDsByDelToValDstIndexKey(delAddr, valDstAddr)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close() //nolint:errcheck
	return iterator.Valid()
}

// HasMaxRedelegationEntries checks if redelegation has maximum number of entries.
func (k Keeper) HasMaxRedelegationEntries(ctx sdk.Context, delegatorAddr sdk.AccAddress, governorSrcAddr, governorDstAddr sdk.ValAddress) bool {
	red, found := k.GetRedelegation(ctx, delegatorAddr, governorSrcAddr, governorDstAddr)
	if !found {
		return false
	}

	return len(red.Entries) >= int(k.MaxEntries(ctx))
}

// SetRedelegation set a redelegation and associated index.
func (k Keeper) SetRedelegation(ctx sdk.Context, red stakingtypes.Redelegation) {
	delegatorAddress := sdk.MustAccAddressFromBech32(red.DelegatorAddress)

	store := ctx.KVStore(k.storeKey)
	bz := stakingtypes.MustMarshalRED(k.cdc, red)
	valSrcAddr, err := sdk.ValAddressFromBech32(red.ValidatorSrcAddress)
	if err != nil {
		panic(err)
	}
	valDestAddr, err := sdk.ValAddressFromBech32(red.ValidatorDstAddress)
	if err != nil {
		panic(err)
	}
	key := types.GetREDKey(delegatorAddress, valSrcAddr, valDestAddr)
	store.Set(key, bz)
	store.Set(types.GetREDByValSrcIndexKey(delegatorAddress, valSrcAddr, valDestAddr), []byte{})
	store.Set(types.GetREDByValDstIndexKey(delegatorAddress, valSrcAddr, valDestAddr), []byte{})
}

// SetRedelegationEntry adds an entry to the unbonding delegation at the given
// addresses. It creates the unbonding delegation if it does not exist.
func (k Keeper) SetRedelegationEntry(ctx sdk.Context,
	delegatorAddr sdk.AccAddress, governorSrcAddr,
	governorDstAddr sdk.ValAddress, creationHeight int64,
	minTime time.Time, balance math.Int,
	sharesSrc, sharesDst sdk.Dec,
) stakingtypes.Redelegation {
	red, found := k.GetRedelegation(ctx, delegatorAddr, governorSrcAddr, governorDstAddr)
	if found {
		red.AddEntry(creationHeight, minTime, balance, sharesDst)
	} else {
		red = stakingtypes.NewRedelegation(delegatorAddr, governorSrcAddr,
			governorDstAddr, creationHeight, minTime, balance, sharesDst)
	}

	k.SetRedelegation(ctx, red)

	return red
}

// IterateRedelegations iterates through all redelegations.
func (k Keeper) IterateRedelegations(ctx sdk.Context, fn func(index int64, red stakingtypes.Redelegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.RedelegationKey)
	defer iterator.Close() //nolint:errcheck
	for i := int64(0); iterator.Valid(); iterator.Next() {
		red := stakingtypes.MustUnmarshalRED(k.cdc, iterator.Value())
		if stop := fn(i, red); stop {
			break
		}
		i++
	}
}

// RemoveRedelegation removes a redelegation object and associated index.
func (k Keeper) RemoveRedelegation(ctx sdk.Context, red stakingtypes.Redelegation) {
	delegatorAddress := sdk.MustAccAddressFromBech32(red.DelegatorAddress)

	store := ctx.KVStore(k.storeKey)
	valSrcAddr, err := sdk.ValAddressFromBech32(red.ValidatorSrcAddress)
	if err != nil {
		panic(err)
	}
	valDestAddr, err := sdk.ValAddressFromBech32(red.ValidatorDstAddress)
	if err != nil {
		panic(err)
	}
	redKey := types.GetREDKey(delegatorAddress, valSrcAddr, valDestAddr)
	store.Delete(redKey)
	store.Delete(types.GetREDByValSrcIndexKey(delegatorAddress, valSrcAddr, valDestAddr))
	store.Delete(types.GetREDByValDstIndexKey(delegatorAddress, valSrcAddr, valDestAddr))
}

// redelegation queue timeslice operations

// GetRedelegationQueueTimeSlice gets a specific redelegation queue timeslice. A
// timeslice is a slice of DVVTriplets corresponding to redelegations that
// expire at a certain time.
func (k Keeper) GetRedelegationQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvvTriplets []stakingtypes.DVVTriplet) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetRedelegationTimeKey(timestamp))
	if bz == nil {
		return []stakingtypes.DVVTriplet{}
	}

	triplets := stakingtypes.DVVTriplets{}
	k.cdc.MustUnmarshal(bz, &triplets)

	return triplets.Triplets
}

// SetRedelegationQueueTimeSlice sets a specific redelegation queue timeslice.
func (k Keeper) SetRedelegationQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []stakingtypes.DVVTriplet) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stakingtypes.DVVTriplets{Triplets: keys})
	store.Set(types.GetRedelegationTimeKey(timestamp), bz)
}

// InsertRedelegationQueue insert an redelegation delegation to the appropriate
// timeslice in the redelegation queue.
func (k Keeper) InsertRedelegationQueue(ctx sdk.Context, red stakingtypes.Redelegation, completionTime time.Time) {
	timeSlice := k.GetRedelegationQueueTimeSlice(ctx, completionTime)
	dvvTriplet := stakingtypes.DVVTriplet{
		DelegatorAddress:    red.DelegatorAddress,
		ValidatorSrcAddress: red.ValidatorSrcAddress,
		ValidatorDstAddress: red.ValidatorDstAddress,
	}

	if len(timeSlice) == 0 {
		k.SetRedelegationQueueTimeSlice(ctx, completionTime, []stakingtypes.DVVTriplet{dvvTriplet})
	} else {
		timeSlice = append(timeSlice, dvvTriplet)
		k.SetRedelegationQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// RedelegationQueueIterator returns all the redelegation queue timeslices from
// time 0 until endTime.
func (k Keeper) RedelegationQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.RedelegationQueueKey, sdk.InclusiveEndBytes(types.GetRedelegationTimeKey(endTime)))
}

// DequeueAllMatureRedelegationQueue returns a concatenated list of all the
// timeslices inclusively previous to currTime, and deletes the timeslices from
// the queue.
func (k Keeper) DequeueAllMatureRedelegationQueue(ctx sdk.Context, currTime time.Time) (matureRedelegations []stakingtypes.DVVTriplet) {
	store := ctx.KVStore(k.storeKey)

	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	redelegationTimesliceIterator := k.RedelegationQueueIterator(ctx, ctx.BlockHeader().Time)
	defer redelegationTimesliceIterator.Close() // nolint: errcheck

	for ; redelegationTimesliceIterator.Valid(); redelegationTimesliceIterator.Next() {
		timeslice := stakingtypes.DVVTriplets{}
		value := redelegationTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		matureRedelegations = append(matureRedelegations, timeslice.Triplets...)

		store.Delete(redelegationTimesliceIterator.Key())
	}

	return matureRedelegations
}

// Delegate performs a delegation, set/update everything necessary within the store.
// tokenSrc indicates the bond status of the incoming funds.
func (k Keeper) Delegate(
	ctx sdk.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc types.BondStatus,
	governor types.Governor, subtractAccount bool,
) (newShares sdk.Dec, err error) {
	// In some situations, the exchange rate becomes invalid, e.g. if
	// Governor loses all tokens due to slashing. In this case,
	// make all future delegations invalid.
	if governor.InvalidExRate() {
		return sdk.ZeroDec(), types.ErrDelegatorShareExRateInvalid
	}

	// Get or create the delegation object
	delegation, found := k.GetDelegation(ctx, delAddr, governor.GetOperator())
	if !found {
		delegation = stakingtypes.NewDelegation(delAddr, governor.GetOperator(), sdk.ZeroDec())
	}

	// call the appropriate hook if present
	if found {
		err = k.BeforeDelegationSharesModified(ctx, delAddr, governor.GetOperator())
	} else {
		err = k.BeforeDelegationCreated(ctx, delAddr, governor.GetOperator())
	}

	if err != nil {
		return sdk.ZeroDec(), err
	}

	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)

	// if subtractAccount is true then we are
	// performing a delegation and not a redelegation, thus the source tokens are
	// all non bonded
	if subtractAccount {
		if tokenSrc == types.Bonded {
			panic("delegation token source cannot be bonded")
		}

		var sendName string

		switch {
		case governor.IsBonded():
			sendName = types.BondedPoolName
		case governor.IsUnbonding(), governor.IsUnbonded():
			sendName = types.NotBondedPoolName
		default:
			panic("invalid governor status")
		}

		coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), bondAmt))
		if err := k.bankKeeper.DelegateCoinsFromAccountToModule(ctx, delegatorAddress, sendName, coins); err != nil {
			return sdk.Dec{}, err
		}
	} else {
		// potentially transfer tokens between pools, if
		switch {
		case tokenSrc == types.Bonded && governor.IsBonded():
			// do nothing
		case (tokenSrc == types.Unbonded || tokenSrc == types.Unbonding) && !governor.IsBonded():
			// do nothing
		case (tokenSrc == types.Unbonded || tokenSrc == types.Unbonding) && governor.IsBonded():
			// transfer pools
			k.notBondedTokensToBonded(ctx, bondAmt)
		case tokenSrc == types.Bonded && !governor.IsBonded():
			// transfer pools
			k.bondedTokensToNotBonded(ctx, bondAmt)
		default:
			panic("unknown token source bond status")
		}
	}

	_, newShares = k.AddGovernorTokensAndShares(ctx, governor, bondAmt)

	// Update delegation
	delegation.Shares = delegation.Shares.Add(newShares)
	k.SetDelegation(ctx, delegation)

	// Call the after-modification hook
	if err := k.AfterDelegationModified(ctx, delegatorAddress, delegation.GetValidatorAddr()); err != nil {
		return newShares, err
	}

	return newShares, nil
}

// Unbond unbonds a particular delegation and perform associated store operations.
func (k Keeper) Unbond(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec,
) (amount math.Int, err error) {
	// check if a delegation object exists in the store
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return amount, types.ErrNoDelegatorForAddress
	}

	// call the before-delegation-modified hook
	if err := k.BeforeDelegationSharesModified(ctx, delAddr, valAddr); err != nil {
		return amount, err
	}

	// ensure that we have enough shares to remove
	if delegation.Shares.LT(shares) {
		return amount, sdkerrors.Wrap(types.ErrNotEnoughDelegationShares, delegation.Shares.String())
	}

	// get governor
	governor, found := k.GetGovernor(ctx, valAddr)
	if !found {
		return amount, types.ErrNoGovernorFound
	}

	// subtract shares from delegation
	delegation.Shares = delegation.Shares.Sub(shares)

	delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		return amount, err
	}

	isGovernorOperator := delegatorAddress.Equals(governor.GetOperator())

	// If the delegation is the operator of the governor and undelegating will decrease the governor's
	// self-delegation below their minimum, we jail the governor.
	if isGovernorOperator && !governor.Jailed &&
		governor.TokensFromShares(delegation.Shares).TruncateInt().LT(governor.MinSelfDelegation) {
		k.jailGovernor(ctx, governor)
		governor = k.mustGetGovernor(ctx, governor.GetOperator())
	}

	if delegation.Shares.IsZero() {
		err = k.RemoveDelegation(ctx, delegation)
	} else {
		k.SetDelegation(ctx, delegation)
		// call the after delegation modification hook
		err = k.AfterDelegationModified(ctx, delegatorAddress, delegation.GetValidatorAddr())
	}

	if err != nil {
		return amount, err
	}

	// remove the shares and coins from the governor
	// NOTE that the amount is later (in keeper.Delegation) moved between staking module pools
	governor, amount = k.RemoveGovernorTokensAndShares(ctx, governor, shares)
	if governor.DelegatorShares.IsZero() && governor.IsUnbonded() {
		// if not unbonded, we must instead remove governor in EndBlocker once it finishes its unbonding period
		k.RemoveGovernor(ctx, governor.GetOperator())
	}

	return amount, nil
}

// getBeginInfo returns the completion time and height of a redelegation, along
// with a boolean signaling if the redelegation is complete based on the source
// governor.
func (k Keeper) getBeginInfo(
	ctx sdk.Context, valSrcAddr sdk.ValAddress,
) (completionTime time.Time, height int64, completeNow bool) {
	governor, found := k.GetGovernor(ctx, valSrcAddr)

	// TODO: When would the governor not be found?
	switch {
	case !found || governor.IsBonded():
		// the longest wait - just unbonding period from now
		completionTime = ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
		height = ctx.BlockHeight()

		return completionTime, height, false

	case governor.IsUnbonded():
		return completionTime, height, true

	case governor.IsUnbonding():
		return governor.UnbondingTime, governor.UnbondingHeight, false

	default:
		panic(fmt.Sprintf("unknown governor status: %s", governor.Status))
	}
}

// Undelegate unbonds an amount of delegator shares from a given governor. It
// will verify that the unbonding entries between the delegator and governor
// are not exceeded and unbond the staked tokens (based on shares) by creating
// an unbonding object and inserting it into the unbonding queue which will be
// processed during the staking EndBlocker.
func (k Keeper) Undelegate(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount sdk.Dec,
) (time.Time, error) {
	governor, found := k.GetGovernor(ctx, valAddr)
	if !found {
		return time.Time{}, types.ErrNoDelegatorForAddress
	}

	if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
		return time.Time{}, types.ErrMaxUnbondingDelegationEntries
	}

	returnAmount, err := k.Unbond(ctx, delAddr, valAddr, sharesAmount)
	if err != nil {
		return time.Time{}, err
	}

	// transfer the governor tokens to the not bonded pool
	if governor.IsBonded() {
		k.bondedTokensToNotBonded(ctx, returnAmount)
	}

	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	ubd := k.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	k.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
}

// CompleteUnbonding completes the unbonding of all mature entries in the
// retrieved unbonding delegation object and returns the total unbonding balance
// or an error upon failure.
func (k Keeper) CompleteUnbonding(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	ubd, found := k.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, types.ErrNoUnbondingDelegation
	}

	bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	delegatorAddress, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			ubd.RemoveEntry(int64(i))
			i--

			// track undelegation only when remaining or truncated shares are non-zero
			if !entry.Balance.IsZero() {
				amt := sdk.NewCoin(bondDenom, entry.Balance)
				if err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(
					ctx, types.NotBondedPoolName, delegatorAddress, sdk.NewCoins(amt),
				); err != nil {
					return nil, err
				}

				balances = balances.Add(amt)
			}
		}
	}

	// set the unbonding delegation or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUnbondingDelegation(ctx, ubd)
	} else {
		k.SetUnbondingDelegation(ctx, ubd)
	}

	return balances, nil
}

// BeginRedelegation begins unbonding / redelegation and creates a redelegation
// record.
func (k Keeper) BeginRedelegation(
	ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount sdk.Dec,
) (completionTime time.Time, err error) {
	if bytes.Equal(valSrcAddr, valDstAddr) {
		return time.Time{}, types.ErrSelfRedelegation
	}

	dstGovernor, found := k.GetGovernor(ctx, valDstAddr)
	if !found {
		return time.Time{}, types.ErrBadRedelegationDst
	}

	srcGovernor, found := k.GetGovernor(ctx, valSrcAddr)
	if !found {
		return time.Time{}, types.ErrBadRedelegationDst
	}

	// check if this is a transitive redelegation
	if k.HasReceivingRedelegation(ctx, delAddr, valSrcAddr) {
		return time.Time{}, types.ErrTransitiveRedelegation
	}

	if k.HasMaxRedelegationEntries(ctx, delAddr, valSrcAddr, valDstAddr) {
		return time.Time{}, types.ErrMaxRedelegationEntries
	}

	returnAmount, err := k.Unbond(ctx, delAddr, valSrcAddr, sharesAmount)
	if err != nil {
		return time.Time{}, err
	}

	if returnAmount.IsZero() {
		return time.Time{}, types.ErrTinyRedelegationAmount
	}

	sharesCreated, err := k.Delegate(ctx, delAddr, returnAmount, srcGovernor.GetStatus(), dstGovernor, false)
	if err != nil {
		return time.Time{}, err
	}

	// create the unbonding delegation
	completionTime, height, completeNow := k.getBeginInfo(ctx, valSrcAddr)

	if completeNow { // no need to create the redelegation object
		return completionTime, nil
	}

	red := k.SetRedelegationEntry(
		ctx, delAddr, valSrcAddr, valDstAddr,
		height, completionTime, returnAmount, sharesAmount, sharesCreated,
	)
	k.InsertRedelegationQueue(ctx, red, completionTime)

	return completionTime, nil
}

// CompleteRedelegation completes the redelegations of all mature entries in the
// retrieved redelegation object and returns the total redelegation (initial)
// balance or an error upon failure.
func (k Keeper) CompleteRedelegation(
	ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress,
) (sdk.Coins, error) {
	red, found := k.GetRedelegation(ctx, delAddr, valSrcAddr, valDstAddr)
	if !found {
		return nil, types.ErrNoRedelegation
	}

	bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	// loop through all the entries and complete mature redelegation entries
	for i := 0; i < len(red.Entries); i++ {
		entry := red.Entries[i]
		if entry.IsMature(ctxTime) {
			red.RemoveEntry(int64(i))
			i--

			if !entry.InitialBalance.IsZero() {
				balances = balances.Add(sdk.NewCoin(bondDenom, entry.InitialBalance))
			}
		}
	}

	// set the redelegation or remove it if there are no more entries
	if len(red.Entries) == 0 {
		k.RemoveRedelegation(ctx, red)
	} else {
		k.SetRedelegation(ctx, red)
	}

	return balances, nil
}

// ValidateUnbondAmount validates that a given unbond or redelegation amount is
// valied based on upon the converted shares. If the amount is valid, the total
// amount of respective shares is returned, otherwise an error is returned.
func (k Keeper) ValidateUnbondAmount(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amt math.Int,
) (shares sdk.Dec, err error) {
	governor, found := k.GetGovernor(ctx, valAddr)
	if !found {
		return shares, types.ErrNoGovernorFound
	}

	del, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return shares, types.ErrNoDelegation
	}

	shares, err = governor.SharesFromTokens(amt)
	if err != nil {
		return shares, err
	}

	sharesTruncated, err := governor.SharesFromTokensTruncated(amt)
	if err != nil {
		return shares, err
	}

	delShares := del.GetShares()
	if sharesTruncated.GT(delShares) {
		return shares, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid shares amount")
	}

	// Cap the shares at the delegation's shares. Shares being greater could occur
	// due to rounding, however we don't want to truncate the shares or take the
	// minimum because we want to allow for the full withdraw of shares from a
	// delegation.
	if shares.GT(delShares) {
		shares = delShares
	}

	return shares, nil
}
