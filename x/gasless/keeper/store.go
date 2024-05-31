package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/dymensionxyz/dymension-rdk/utils"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

func (k Keeper) GetUsageIdentifierToGasTankIds(ctx sdk.Context, usageIdentifier string) (usageIdentifierToGasTankIds types.UsageIdentifierToGasTankIds, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUsageIdentifierToGasTankIdsKey(usageIdentifier))
	if bz == nil {
		return
	}
	usageIdentifierToGasTankIds = types.MustUnmarshalUsageIdentifierToGastankIds(k.cdc, bz)
	return usageIdentifierToGasTankIds, true
}

func (k Keeper) IterateAllUsageIdentifierToGasTankIds(ctx sdk.Context, cb func(usageIdentifierToGasTankIds types.UsageIdentifierToGasTankIds) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllUsageIdentifierToGasTankIdsKey())
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		usageIdentifierToGasTankIds := types.MustUnmarshalUsageIdentifierToGastankIds(k.cdc, iter.Value())
		stop, err := cb(usageIdentifierToGasTankIds)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func (k Keeper) GetAllUsageIdentifierToGasTankIds(ctx sdk.Context) (allUsageIdentifierToGasTankIds []types.UsageIdentifierToGasTankIds) {
	allUsageIdentifierToGasTankIds = []types.UsageIdentifierToGasTankIds{}
	_ = k.IterateAllUsageIdentifierToGasTankIds(ctx, func(usageIdentifierToGasTankIds types.UsageIdentifierToGasTankIds) (stop bool, err error) {
		allUsageIdentifierToGasTankIds = append(allUsageIdentifierToGasTankIds, usageIdentifierToGasTankIds)
		return false, nil
	})
	return allUsageIdentifierToGasTankIds
}

func (k Keeper) SetUsageIdentifierToGasTankIds(ctx sdk.Context, usageIdentifierToGasTankIds types.UsageIdentifierToGasTankIds) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalUsageIdentifierToGastankIds(k.cdc, usageIdentifierToGasTankIds)
	store.Set(types.GetUsageIdentifierToGasTankIdsKey(usageIdentifierToGasTankIds.UsageIdentifier), bz)
}

// DeleteUsageIdentifierToGasTankIds deletes an UsageIdentifierToGasTankIds.
func (k Keeper) DeleteUsageIdentifierToGasTankIds(ctx sdk.Context, usageIdentifierToGasTankIds types.UsageIdentifierToGasTankIds) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetUsageIdentifierToGasTankIdsKey(usageIdentifierToGasTankIds.UsageIdentifier))
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
	// eg. if length of existing consumption is 2, so after adding new consumption the index of appended consuption will also be 2 since sequence begins from 0
	return gasConsumer, consumptionLength
}

func (k Keeper) AddGasTankIdToUsageIdentifiers(ctx sdk.Context, usageIdentifiers []string, gasTankID uint64) {
	for _, usageIdentifier := range usageIdentifiers {
		usageIdentifierToGasTankIds, found := k.GetUsageIdentifierToGasTankIds(ctx, usageIdentifier)
		if !found {
			usageIdentifierToGasTankIds = types.NewUsageIdentifierToGastankIds(usageIdentifier)
		}
		usageIdentifierToGasTankIds.GasTankIds = append(usageIdentifierToGasTankIds.GasTankIds, gasTankID)
		usageIdentifierToGasTankIds.GasTankIds = utils.RemoveDuplicates(usageIdentifierToGasTankIds.GasTankIds)
		k.SetUsageIdentifierToGasTankIds(ctx, usageIdentifierToGasTankIds)
	}
}

func (k Keeper) RemoveGasTankIdFromUsageIdentifiers(ctx sdk.Context, usageIdentifiers []string, gasTankID uint64) {
	for _, usageIdentifier := range usageIdentifiers {
		usageIdentifierToGasTankIds, found := k.GetUsageIdentifierToGasTankIds(ctx, usageIdentifier)
		if !found {
			continue
		}
		usageIdentifierToGasTankIds.GasTankIds = utils.RemoveValueFromSlice(usageIdentifierToGasTankIds.GasTankIds, gasTankID)
		if len(usageIdentifierToGasTankIds.GasTankIds) == 0 {
			k.DeleteUsageIdentifierToGasTankIds(ctx, usageIdentifierToGasTankIds)
			continue
		}
		k.SetUsageIdentifierToGasTankIds(ctx, usageIdentifierToGasTankIds)
	}
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
