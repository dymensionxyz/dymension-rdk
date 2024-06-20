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

func (q Querier) State(goCtx context.Context, request *types.QueryStateRequest) (*types.QueryStateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// TODO:
	return &types.QueryStateResponse{State: q.Keeper.GetState(ctx)}, nil
}
