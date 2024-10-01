package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

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
