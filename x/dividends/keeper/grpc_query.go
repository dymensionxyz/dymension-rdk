package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Querier{}

type Querier struct {
	keeper Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{keeper: k}
}

func (q Querier) GaugeByID(goCtx context.Context, req *types.GaugeByIDRequest) (*types.GaugeByIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	gauge, err := q.keeper.GetGauge(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.GaugeByIDResponse{Gauge: gauge}, nil
}

func (q Querier) Gauges(goCtx context.Context, req *types.GaugesRequest) (*types.GaugesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	gauges, err := q.keeper.GetAllGauges(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.GaugesResponse{Data: gauges, Pagination: nil}, nil
}

func (q Querier) Params(goCtx context.Context, req *types.ParamsRequest) (*types.ParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.keeper.MustGetParams(ctx)

	return &types.ParamsResponse{Params: &params}, nil
}
