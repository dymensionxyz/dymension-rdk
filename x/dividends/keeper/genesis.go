package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) error {
	err := k.SetLastGaugeId(ctx, genState.LastGaugeId)
	if err != nil {
		return fmt.Errorf("set last gauge id: %w", err)
	}
	err = k.SetParams(ctx, genState.Params)
	if err != nil {
		return fmt.Errorf("set params: %w", err)
	}
	for _, gauge := range genState.Gauges {
		err = k.SetGauge(ctx, gauge)
		if err != nil {
			return fmt.Errorf("set gauge: %w", err)
		}
	}
	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	gauges, err := k.GetAllGauges(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all gauges: %w", err)
	}
	lastGaugeId, err := k.NextGaugeId(ctx)
	if err != nil {
		return nil, fmt.Errorf("next gauge id: %w", err)
	}
	return &types.GenesisState{
		LastGaugeId: lastGaugeId,
		Params:      k.MustGetParams(ctx),
		Gauges:      gauges,
	}, nil
}
