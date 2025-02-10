package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

func (k Keeper) NextGaugeId(ctx sdk.Context) (uint64, error) {
	return k.lastGaugeID.Next(ctx)
}

func (k Keeper) SetLastGaugeId(ctx sdk.Context, id uint64) error {
	return k.lastGaugeID.Set(ctx, id)
}

func (k Keeper) GetAllGauges(ctx sdk.Context) ([]types.Gauge, error) {
	i, err := k.gauges.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer i.Close()
	return i.Values()
}

func (k Keeper) IterateActiveGauges(ctx sdk.Context, fn func(gauge types.Gauge) (stop bool, err error)) error {
	return k.gauges.Walk(
		ctx,
		collections.NewPrefixedPairRange[bool, uint64](true), // only active gauges
		func(_ collections.Pair[bool, uint64], g types.Gauge) (stop bool, err error) { return fn(g) },
	)
}

func (k Keeper) SetGauge(ctx sdk.Context, gauge types.Gauge) error {
	gaugeKey := collections.Join(gauge.Active, gauge.Id)
	return k.gauges.Set(ctx, gaugeKey, gauge)
}

func (k Keeper) GetGauge(ctx sdk.Context, gaugeId uint64) (types.Gauge, error) {
	activeGaugeKey := collections.Join(true, gaugeId)
	gauge, err := k.gauges.Get(ctx, activeGaugeKey)
	if err == nil {
		return gauge, nil
	}
	if !errors.Is(err, collections.ErrNotFound) {
		return types.Gauge{}, err
	}

	inactiveGaugeKey := collections.Join(false, gaugeId)
	return k.gauges.Get(ctx, inactiveGaugeKey)
}

func (k Keeper) GetActiveGauge(ctx sdk.Context, gaugeId uint64) (types.Gauge, error) {
	activeGaugeKey := collections.Join(true, gaugeId)
	return k.gauges.Get(ctx, activeGaugeKey)
}

func (k Keeper) DeactivateGauge(ctx sdk.Context, gaugeId uint64) error {
	activeGaugeKey := collections.Join(true, gaugeId)
	gauge, err := k.gauges.Get(ctx, activeGaugeKey)
	if err != nil {
		return err
	}
	if err = k.gauges.Remove(ctx, activeGaugeKey); err != nil {
		return err
	}
	inactiveGaugeKey := collections.Join(false, gaugeId)
	gauge.Active = false
	return k.gauges.Set(ctx, inactiveGaugeKey, gauge)
}

func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
	return k.params.Set(ctx, p)
}

func (k Keeper) MustGetParams(ctx sdk.Context) types.Params {
	p, err := k.params.Get(ctx)
	if err != nil {
		panic(err)
	}
	return p
}
