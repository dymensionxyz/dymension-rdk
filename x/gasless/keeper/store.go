package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

func (k Keeper) GetUsageIdentifierToGasTankIds(ctx sdk.Context, usageIdentifier string) (types.UsageIdentifierToGasTankIds, bool) {
	var gasTankIDs []uint64
	ranger := collections.NewPrefixedPairRange[string, uint64](usageIdentifier)
	err := k.usageIdentifierToGasTankIDSet.Walk(ctx, ranger, func(key collections.Pair[string, uint64]) (bool, error) {
		gasTankIDs = append(gasTankIDs, key.K2())
		return false, nil
	})
	if err != nil || len(gasTankIDs) == 0 {
		return types.UsageIdentifierToGasTankIds{}, false
	}
	return types.UsageIdentifierToGasTankIds{
		UsageIdentifier: usageIdentifier,
		GasTankIds:      gasTankIDs,
	}, true
}

func (k Keeper) GetAllUsageIdentifierToGasTankIds(ctx sdk.Context) (allUsageIdentifierToGasTankIds []types.UsageIdentifierToGasTankIds, err error) {
	err = k.usageIdentifierToGasTankIDSet.Walk(ctx, nil, func(key collections.Pair[string, uint64]) (stop bool, err error) {
		for i, usageIdentifierToGasTankIds := range allUsageIdentifierToGasTankIds {
			if usageIdentifierToGasTankIds.UsageIdentifier == key.K1() {
				allUsageIdentifierToGasTankIds[i].GasTankIds = append(usageIdentifierToGasTankIds.GasTankIds, key.K2())
				return
			}
		}
		allUsageIdentifierToGasTankIds = append(allUsageIdentifierToGasTankIds, types.UsageIdentifierToGasTankIds{
			UsageIdentifier: key.K1(),
			GasTankIds:      []uint64{key.K2()},
		})
		return
	})
	return
}

func (k Keeper) SetUsageIdentifierToGasTankIds(ctx sdk.Context, usageIdentifierToGasTankIds types.UsageIdentifierToGasTankIds) error {
	usageIdentifier := usageIdentifierToGasTankIds.UsageIdentifier
	err := k.RemoveAllGasTankIdsForUsageIdentifier(ctx, usageIdentifier)
	if err != nil {
		return err
	}
	for _, gasTankID := range usageIdentifierToGasTankIds.GasTankIds {
		key := collections.Join(usageIdentifier, gasTankID)
		err := k.usageIdentifierToGasTankIDSet.Set(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteUsageIdentifierToGasTankIds deletes an UsageIdentifierToGasTankIds.
func (k Keeper) DeleteUsageIdentifierToGasTankIds(ctx sdk.Context, usageIdentifierToGasTankIds types.UsageIdentifierToGasTankIds) error {
	return k.RemoveAllGasTankIdsForUsageIdentifier(ctx, usageIdentifierToGasTankIds.UsageIdentifier)
}

func (k Keeper) GetLastGasTankID(ctx sdk.Context) (id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastGasTankIDKey())
	if bz == nil {
		id = 0 // initialize the GasTankID
	} else {
		var val gogotypes.UInt64Value
		k.cdc.MustUnmarshal(bz, &val)
		id = val.GetValue()
	}
	return
}

func (k Keeper) SetLastGasTankID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})
	store.Set(types.GetLastGasTankIDKey(), bz)
}

func (k Keeper) GetNextGasTankIDWithUpdate(ctx sdk.Context) uint64 {
	id := k.GetLastGasTankID(ctx) + 1
	k.SetLastGasTankID(ctx, id)
	return id
}

func (k Keeper) GetGasTank(ctx sdk.Context, id uint64) (gasTank types.GasTank, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGasTankKey(id))
	if bz == nil {
		return
	}
	gasTank = types.MustUnmarshalGasTank(k.cdc, bz)
	return gasTank, true
}

func (k Keeper) IterateAllGasTanks(ctx sdk.Context, cb func(gasTank types.GasTank) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllGasTanksKey())
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		gasTank := types.MustUnmarshalGasTank(k.cdc, iter.Value())
		stop, err := cb(gasTank)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func (k Keeper) GetAllGasTanks(ctx sdk.Context) (gasTanks []types.GasTank) {
	gasTanks = []types.GasTank{}
	_ = k.IterateAllGasTanks(ctx, func(gasTank types.GasTank) (stop bool, err error) {
		gasTanks = append(gasTanks, gasTank)
		return false, nil
	})
	return gasTanks
}

func (k Keeper) SetGasTank(ctx sdk.Context, gasTank types.GasTank) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalGasTank(k.cdc, gasTank)
	store.Set(types.GetGasTankKey(gasTank.Id), bz)
}

func (k Keeper) GetGasConsumer(ctx sdk.Context, consumer sdk.AccAddress) (gasConsumer types.GasConsumer, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGasConsumerKey(consumer))
	if bz == nil {
		return
	}
	gasConsumer = types.MustUnmarshalGasConsumer(k.cdc, bz)
	return gasConsumer, true
}

func (k Keeper) IterateAllGasConsumers(ctx sdk.Context, cb func(gasConsumer types.GasConsumer) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllGasConsumersKey())
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		gasConsumer := types.MustUnmarshalGasConsumer(k.cdc, iter.Value())
		stop, err := cb(gasConsumer)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func (k Keeper) GetAllGasConsumers(ctx sdk.Context) (gasConsumers []types.GasConsumer) {
	gasConsumers = []types.GasConsumer{}
	_ = k.IterateAllGasConsumers(ctx, func(gasConsumer types.GasConsumer) (stop bool, err error) {
		gasConsumers = append(gasConsumers, gasConsumer)
		return false, nil
	})
	return gasConsumers
}

func (k Keeper) SetGasConsumer(ctx sdk.Context, gasConsumer types.GasConsumer) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalGasConsumer(k.cdc, gasConsumer)
	store.Set(types.GetGasConsumerKey(sdk.MustAccAddressFromBech32(gasConsumer.Consumer)), bz)
}

func (k Keeper) GetOrCreateGasConsumer(ctx sdk.Context, consumer sdk.AccAddress, gasTank types.GasTank) (gasConsumer types.GasConsumer, consumptionIndex uint64) {
	gasConsumer, found := k.GetGasConsumer(ctx, consumer)
	if !found {
		gasConsumer = types.NewGasConsumer(consumer)
	}

	consumptionLength := uint64(0)
	for consumptionIndex, consumption := range gasConsumer.Consumptions {
		if consumption.GasTankId == gasTank.Id {
			return gasConsumer, uint64(consumptionIndex)
		}
		consumptionLength++
	}

	gasConsumer.Consumptions = append(gasConsumer.Consumptions, types.NewConsumptionDetail(
		gasTank.Id,
		gasTank.MaxFeeUsagePerConsumer,
	))
	k.SetGasConsumer(ctx, gasConsumer)
	// eg. if length of existing consumption is 2, so after adding new consumption the index of appended consumption will also be 2 since sequence begins from 0
	return gasConsumer, consumptionLength
}

func (k Keeper) AddGasTankIdToUsageIdentifiers(ctx sdk.Context, usageIdentifiers []string, gasTankID uint64) error {
	for _, usageIdentifier := range usageIdentifiers {
		key := collections.Join(usageIdentifier, gasTankID)
		if err := k.usageIdentifierToGasTankIDSet.Set(ctx, key); err != nil {
			return fmt.Errorf("set gas tank id to usage identifier: %w", err)
		}
	}
	return nil
}

func (k Keeper) RemoveGasTankIdFromUsageIdentifiers(ctx sdk.Context, usageIdentifiers []string, gasTankID uint64) error {
	for _, usageIdentifier := range usageIdentifiers {
		key := collections.Join(usageIdentifier, gasTankID)
		if err := k.usageIdentifierToGasTankIDSet.Remove(ctx, key); err != nil {
			return fmt.Errorf("remove gas tank id from usage identifier: %w", err)
		}
	}
	return nil
}

func (k Keeper) RemoveAllGasTankIdsForUsageIdentifier(ctx sdk.Context, usageIdentifier string) error {
	ranger := collections.NewPrefixedPairRange[string, uint64](usageIdentifier)
	return k.usageIdentifierToGasTankIDSet.Clear(ctx, ranger)
}

func (k Keeper) UpdateConsumerAllowance(ctx sdk.Context, gasTank types.GasTank) {
	allConsumers := k.GetAllGasConsumers(ctx)
	for _, consumer := range allConsumers {
		for index, consumption := range consumer.Consumptions {
			if consumption.GasTankId == gasTank.Id {
				consumer.Consumptions[index].TotalFeeConsumptionAllowed = gasTank.MaxFeeUsagePerConsumer
				k.SetGasConsumer(ctx, consumer)
				break
			}
		}
	}
}

func (k Keeper) LastUsedGasTankID(ctx sdk.Context, usageIdentifier string) (uint64, error) {
	lastUsedGasTankID, err := k.lastUsedGasTankIDMap.Get(ctx, usageIdentifier)
	if err != nil {
		return 0, err
	}
	return lastUsedGasTankID, nil
}
