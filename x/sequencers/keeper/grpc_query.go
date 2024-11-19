package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
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
		return nil, fmt.Errorf("ValAddressFromBech32: %w", err)
	}

	val, found := k.GetSequencer(ctx, addr)
	if !found {
		return nil, types.ErrSequencerNotFound
	}

	// don't return error if the reward address or whitelisted relayers are not found
	rewardAddr, _ := k.GetRewardAddr(ctx, addr)
	wlr, _ := k.GetWhitelistedRelayers(ctx, addr)

	return &types.QuerySequencerResponse{
		Sequencer:  val,
		RewardAddr: rewardAddr.String(),
		Relayers:   wlr.Relayers,
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

func (k Querier) RewardAddress(goCtx context.Context, req *types.QueryRewardAddressRequest) (*types.QueryRewardAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	addr, err := sdk.ValAddressFromBech32(req.SequencerAddr)
	if err != nil {
		return nil, fmt.Errorf("ValAddressFromBech32: %s", err.Error())
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	rewardAddr, found := k.GetRewardAddr(ctx, addr)
	if !found {
		return nil, types.ErrRewardAddressNotFound
	}

	return &types.QueryRewardAddressResponse{
		RewardAddr: rewardAddr.String(),
	}, nil
}

func (k Querier) WhitelistedRelayers(goCtx context.Context, req *types.QueryWhitelistedRelayersRequest) (*types.QueryWhitelistedRelayersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	addr, err := sdk.ValAddressFromBech32(req.SequencerAddr)
	if err != nil {
		return nil, fmt.Errorf("ValAddressFromBech32: %w", err)
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	wlr, err := k.GetWhitelistedRelayers(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("GetWhitelistedRelayers: %w", err)
	}

	return &types.QueryWhitelistedRelayersResponse{
		Relayers: wlr.Relayers,
	}, nil
}
