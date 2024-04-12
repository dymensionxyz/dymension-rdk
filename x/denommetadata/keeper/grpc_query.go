package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the x/FeeShare keeper providing gRPC method
// handlers.
type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

// Params returns denommetadata module params
func (q Querier) Params(
	c context.Context,
	_ *types.QueryParamsRequest,
) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// IBCDenomByDenomTrace returns IBC denom base on denom trace
func (q Querier) IBCDenomByDenomTrace(
	_ context.Context,
	req *types.QueryGetIBCDenomByDenomTraceRequest,
) (*types.QueryIBCDenomByDenomTraceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	paths := strings.Split(req.Path, "/")
	if len(paths) < 2 {
		return nil, status.Error(codes.InvalidArgument, "input denom traces invalid, need to have at least 2 elements")
	}

	if len(paths)%2 != 0 {
		return nil, status.Error(codes.InvalidArgument, "denom traces must follow this format [port-id-1]/[channel-id-1]/.../[port-id-n]/[channel-id-n]")
	}

	denom := req.Denom

	for i := 0; i < len(paths); i += 2 {
		portID := paths[i]
		channelID := paths[i+1]
		if !strings.Contains(channelID, "channel-") {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("channel %s must contain channel-", channelID))
		}
		tokenDenom := transfertypes.GetPrefixedDenom(portID, channelID, denom)
		denom = transfertypes.ParseDenomTrace(tokenDenom).IBCDenom()
	}

	ibcDenomResponse := &types.QueryIBCDenomByDenomTraceResponse{IbcDenom: denom}
	return ibcDenomResponse, nil
}
