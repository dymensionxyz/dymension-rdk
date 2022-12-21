package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//FIXME
var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// Validators queries all sequencers that match the given status.
func (k Keeper) Validators(c context.Context, req *types.QueryValidatorsRequest) (*types.QueryValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryValidatorsResponse{
		Sequencers: k.GetAllSequencer(ctx),
	}, nil
}

// Validator queries validator info for given validator address.
func (k Keeper) Validator(c context.Context, req *types.QueryValidatorRequest) (*types.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	val, ok := k.GetSequencer(ctx, req.ValidatorAddr)

	return &types.QueryValidatorResponse{
		Validator: ,
	}, nil
}

// HistoricalInfo queries the historical info for given height.
//Implement required methods for IBC expected keeper
func (k Keeper) HistoricalInfo(ctx sdk.Context, height int64) (types.HistoricalInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetHistoricalInfoKey(height)

	value := store.Get(key)
	if value == nil {
		return stakingtypes.HistoricalInfo{}, false
	}

	return types.MustUnmarshalHistoricalInfo(k.cdc, value), true
}
