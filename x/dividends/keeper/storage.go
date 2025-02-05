package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

func (k Keeper) NextGaugeId(ctx sdk.Context) (uint64, error) {
	return k.lastGaugeID.Next(ctx)
}

func (k Keeper) GetAllGauges(ctx sdk.Context) ([]types.Gauge, error) {
	i, err := k.gauges.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer i.Close()
	return i.Values()
}

func (k Keeper) IterateGauges(ctx sdk.Context, fn func(gaugeId uint64, gauge types.Gauge) (stop bool, err error)) error {
	return k.gauges.Walk(ctx, nil, fn)
}

func (k Keeper) SetGauge(ctx sdk.Context, gauge types.Gauge) error {
	return k.gauges.Set(ctx, gauge.Id, gauge)
}

func (k Keeper) GetGauge(ctx sdk.Context, gaugeId uint64) (types.Gauge, error) {
	return k.gauges.Get(ctx, gaugeId)
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
