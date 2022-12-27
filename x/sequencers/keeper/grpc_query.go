package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// var _ types.QueryServer = Keeper{}

func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// Validators queries all sequencers that match the given status.
func (k Querier) Validators(c context.Context, req *types.QueryValidatorsRequest) (*types.QueryValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryValidatorsResponse{
		Sequencers: k.GetAllValidators(ctx),
	}, nil
}

// Validator queries validator info for given validator address.
func (k Querier) Validator(c context.Context, req *types.QueryValidatorRequest) (*types.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	addr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		panic(err)
	}

	val, found := k.GetValidator(ctx, addr)
	if !found {
		return nil, types.ErrSequencerNotFound
	}

	return &types.QueryValidatorResponse{
		Validator: val,
	}, nil
}

// Validator queries validator info for given validator address.
func (k Querier) HistoricalInfo(c context.Context, req *types.QueryHistoricalInfoRequest) (*types.QueryHistoricalInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	histInfo, found := k.GetHistoricalInfo(ctx, req.Height)
	if !found {
		return nil, types.ErrHistoricalInfoNotFound
	}

	return &types.QueryHistoricalInfoResponse{
		Hist: &histInfo,
	}, nil
}
