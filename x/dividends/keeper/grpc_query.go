package keeper

import (
	"context"
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the incentives module keeper providing gRPC method handlers.
type Querier struct {
	Keeper
}

// NewQuerier creates a new Querier struct.
func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

// GaugeByID takes a gaugeID and returns its respective gauge.
func (q Querier) GaugeByID(goCtx context.Context, req *types.GaugeByIDRequest) (*types.GaugeByIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	gauge, err := q.Keeper.GetGaugeByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.GaugeByIDResponse{Gauge: gauge}, nil
}

// Gauges returns all upcoming and active gauges.
func (q Querier) Gauges(goCtx context.Context, req *types.GaugesRequest) (*types.GaugesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	pageRes, gauges, err := q.filterByPrefixAndDenom(ctx, types.KeyPrefixGauges, "", req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.GaugesResponse{Data: gauges, Pagination: pageRes}, nil
}

// Params implements types.QueryServer.
func (q Querier) Params(goCtx context.Context, req *types.ParamsRequest) (*types.ParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.Keeper.GetParams(ctx)

	return &types.ParamsResponse{Params: &params}, nil
}
