package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the x/hub-genesis keeper providing gRPC method
// handlers.
type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

func (q Querier) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{Params: q.Keeper.GetParams(sdk.UnwrapSDKContext(ctx))}, nil
}

func (q Querier) State(ctx context.Context, _ *types.QueryStateRequest) (*types.QueryStateResponse, error) {
	return &types.QueryStateResponse{State: q.Keeper.GetState(sdk.UnwrapSDKContext(ctx))}, nil
}

func (q Querier) GenesisInfo(ctx context.Context, _ *types.QueryGenesisInfoRequest) (*types.QueryGenesisInfoResponse, error) {
	return &types.QueryGenesisInfoResponse{GenesisInfo: q.Keeper.GetGenesisInfo(sdk.UnwrapSDKContext(ctx))}, nil
}

func (q Querier) GenesisBridgeData(goCtx context.Context, _ *types.QueryGenesisBridgeDataRequest) (*types.QueryGenesisBridgeDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	data, err := q.PrepareGenesisBridgeData(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryGenesisBridgeDataResponse{Data: data}, nil
}
