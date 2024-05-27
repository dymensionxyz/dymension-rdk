package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

func (k Keeper) GetTxGTIDs(ctx sdk.Context, txPathOrContractAddress string) (txGTIDs types.TxGTIDs, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTxGTIDsKey(txPathOrContractAddress))
	if bz == nil {
		return
	}
	txGTIDs = types.MustUnmarshalTxGTIDs(k.cdc, bz)
	return txGTIDs, true
}

func (k Keeper) IterateAllTxGTIDs(ctx sdk.Context, cb func(txGTIDs types.TxGTIDs) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllTxGTIDsKey())
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		txGTIDs := types.MustUnmarshalTxGTIDs(k.cdc, iter.Value())
		stop, err := cb(txGTIDs)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func (k Keeper) GetAllTxGTIDs(ctx sdk.Context) (txGTIDss []types.TxGTIDs) {
	txGTIDss = []types.TxGTIDs{}
	_ = k.IterateAllTxGTIDs(ctx, func(txGTIDs types.TxGTIDs) (stop bool, err error) {
		txGTIDss = append(txGTIDss, txGTIDs)
		return false, nil
	})
	return txGTIDss
}

func (k Keeper) SetTxGTIDs(ctx sdk.Context, txGTIDs types.TxGTIDs) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalTxGTIDs(k.cdc, txGTIDs)
	store.Set(types.GetTxGTIDsKey(txGTIDs.TxPathOrContractAddress), bz)
}

// DeleteTxGTIDs deletes an TxGTIDs.
func (k Keeper) DeleteTxGTIDs(ctx sdk.Context, txGTIDs types.TxGTIDs) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetTxGTIDsKey(txGTIDs.TxPathOrContractAddress))
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
		gasTank.MaxTxsCountPerConsumer,
		gasTank.MaxFeeUsagePerConsumer,
	))
	k.SetGasConsumer(ctx, gasConsumer)
	// eg. if length of existing consumption is 2, so after adding new consumption the index of appended consuption will also be 2 since sequence begins from 0
	return gasConsumer, consumptionLength
}

func (k Keeper) AddToTxGtids(ctx sdk.Context, txs, contracts []string, gtid uint64) {
	for _, txPath := range txs {
		txGtids, found := k.GetTxGTIDs(ctx, txPath)
		if !found {
			txGtids = types.NewTxGTIDs(txPath)
		}
		txGtids.GasTankIds = append(txGtids.GasTankIds, gtid)
		txGtids.GasTankIds = types.RemoveDuplicates(txGtids.GasTankIds)
		k.SetTxGTIDs(ctx, txGtids)
	}

	for _, c := range contracts {
		txGtids, found := k.GetTxGTIDs(ctx, c)
		if !found {
			txGtids = types.NewTxGTIDs(c)
		}
		txGtids.GasTankIds = append(txGtids.GasTankIds, gtid)
		txGtids.GasTankIds = types.RemoveDuplicates(txGtids.GasTankIds)
		k.SetTxGTIDs(ctx, txGtids)
	}
}

func (k Keeper) RemoveFromTxGtids(ctx sdk.Context, txs, contracts []string, gtid uint64) {
	for _, txPath := range txs {
		txGtids, found := k.GetTxGTIDs(ctx, txPath)
		if !found {
			continue
		}
		txGtids.GasTankIds = types.RemoveValueFromList(txGtids.GasTankIds, gtid)
		if len(txGtids.GasTankIds) == 0 {
			k.DeleteTxGTIDs(ctx, txGtids)
			continue
		}
		k.SetTxGTIDs(ctx, txGtids)
	}

	for _, c := range contracts {
		txGtids, found := k.GetTxGTIDs(ctx, c)
		if !found {
			continue
		}
		txGtids.GasTankIds = types.RemoveValueFromList(txGtids.GasTankIds, gtid)
		if len(txGtids.GasTankIds) == 0 {
			k.DeleteTxGTIDs(ctx, txGtids)
			continue
		}
		k.SetTxGTIDs(ctx, txGtids)
	}
}

func (k Keeper) UpdateConsumerAllowance(ctx sdk.Context, gasTank types.GasTank) {
	allConsumers := k.GetAllGasConsumers(ctx)
	for _, consumer := range allConsumers {
		for index, consumption := range consumer.Consumptions {
			if consumption.GasTankId == gasTank.Id {
				consumer.Consumptions[index].TotalTxsAllowed = gasTank.MaxTxsCountPerConsumer
				consumer.Consumptions[index].TotalFeeConsumptionAllowed = gasTank.MaxFeeUsagePerConsumer
				k.SetGasConsumer(ctx, consumer)
				break
			}
		}
	}
}
