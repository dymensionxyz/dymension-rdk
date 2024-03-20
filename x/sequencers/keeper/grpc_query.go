package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// Sequencers queries all sequencers that match the given status.
func (k Querier) Sequencers(c context.Context, req *types.QuerySequencersRequest) (*types.QuerySequencersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QuerySequencersResponse{
		Sequencers: k.GetAllSequencers(ctx),
	}, nil
}

// Sequencer queries sequencer info for given sequencer address.
func (k Querier) Sequencer(c context.Context, req *types.QuerySequencerRequest) (*types.QuerySequencerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	addr, err := sdk.ValAddressFromBech32(req.SequencerAddr)
	if err != nil {
		panic(err)
	}

	val, found := k.GetSequencer(ctx, addr)
	if !found {
		return nil, types.ErrSequencerNotFound
	}

	return &types.QuerySequencerResponse{
		Sequencer: val,
	}, nil
}

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
