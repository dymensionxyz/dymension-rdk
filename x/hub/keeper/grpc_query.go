package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the x/hub-genesis keeper providing gRPC method handlers.
type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

func (q Querier) State(goCtx context.Context, _ *types.QueryStateRequest) (*types.QueryStateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	denoms, err := q.GetAllHubDenoms(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all hub denoms: %w", err)
	}
	var state types.State
	for _, denom := range denoms {
		state.Hub.RegisteredDenoms = append(state.Hub.RegisteredDenoms, &types.RegisteredDenom{
			Base: denom,
		})
	}
	return &types.QueryStateResponse{State: state}, nil
}
