package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

type msgServer struct {
	keeper Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateGauge creates a gauge.
func (s msgServer) CreateGauge(goCtx context.Context, msg *types.MsgCreateGauge) (*types.MsgCreateGaugeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: validate query and vesting conditions

	gauge := types.NewGauge(msg.Id, msg.Address, msg.QueryCondition, msg.VestingCondition)

	err := s.keeper.gauges.Set(ctx, gauge.Id, gauge)
	if err != nil {
		return nil, fmt.Errorf("set gauge: %w", err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtCreateGauge,
			sdk.NewAttribute(types.AttributeGaugeID, osmoutils.Uint64ToString(gaugeID)),
		),
	})

	return &types.MsgCreateGaugeResponse{}, nil
}
