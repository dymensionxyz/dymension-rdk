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

// Params returns params of the hub-genesis module.
func (q Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.Keeper.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (q Querier) Hub(goCtx context.Context, request *types.QueryGetHubRequest) (*types.QueryGetHubResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	hub, found := q.Keeper.GetHub(ctx, request.HubId)
	if !found {
		return nil, types.ErrUnknownHubID
	}

	return &types.QueryGetHubResponse{Hub: hub}, nil
}
