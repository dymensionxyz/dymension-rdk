package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Params queries the parameters of the gasless module.
func (k Querier) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var params types.Params
	k.Keeper.paramSpace.GetParamSet(ctx, &params)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Querier) AvailableUsageIdentifiers(goCtx context.Context, _ *types.QueryAvailableUsageIdentifiersRequest) (*types.QueryAvailableUsageIdentifiersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.QueryAvailableUsageIdentifiersResponse{
		UsageIdentifiers: k.GetAvailableUsageIdentifiers(ctx),
	}, nil
}

func (k Querier) GasTank(goCtx context.Context, req *types.QueryGasTankRequest) (*types.QueryGasTankResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GasTankId == 0 {
		return nil, status.Error(codes.InvalidArgument, "gas tank id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	gt, found := k.GetGasTank(ctx, req.GasTankId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "gas tank with id %d doesn't exist", req.GasTankId)
	}

	gasTankBalance := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(gt.Reserve), gt.FeeDenom)
	return &types.QueryGasTankResponse{
		GasTank: types.NewGasTankResponse(gt, gasTankBalance),
	}, nil
}

func (k Querier) GasTanks(goCtx context.Context, req *types.QueryGasTanksRequest) (*types.QueryGasTanksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)

	keyPrefix := types.GetAllGasTanksKey()
	gtGetter := func(_, value []byte) types.GasTank {
		return types.MustUnmarshalGasTank(k.cdc, value)
	}
	gtStore := prefix.NewStore(store, keyPrefix)
	var gasTanks []types.GasTankResponse

	pageRes, err := query.FilteredPaginate(gtStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		gt := gtGetter(key, value)
		if accumulate {
			gasTankBalance := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(gt.Reserve), gt.FeeDenom)
			gasTanks = append(gasTanks, types.NewGasTankResponse(gt, gasTankBalance))
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryGasTanksResponse{
		GasTanks:   gasTanks,
		Pagination: pageRes,
	}, nil
}

func (k Querier) GasTanksByProvider(goCtx context.Context, req *types.QueryGasTanksByProviderRequest) (*types.QueryGasTanksByProviderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if _, err := sdk.AccAddressFromBech32(req.Provider); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid provider address")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	allGasTanks := k.GetAllGasTanks(ctx)

	providerGasTanks := []types.GasTankResponse{}
	for _, tank := range allGasTanks {
		if tank.Provider == req.Provider {
			tankBalance := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(tank.Reserve), tank.FeeDenom)
			providerGasTanks = append(providerGasTanks, types.NewGasTankResponse(tank, tankBalance))
		}
	}
	return &types.QueryGasTanksByProviderResponse{
		GasTanks: providerGasTanks,
	}, nil
}

func (k Querier) GasConsumer(goCtx context.Context, req *types.QueryGasConsumerRequest) (*types.QueryGasConsumerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if _, err := sdk.AccAddressFromBech32(req.Consumer); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid consumer address")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	gc, found := k.GetGasConsumer(ctx, sdk.MustAccAddressFromBech32(req.Consumer))
	if !found {
		return nil, status.Errorf(codes.NotFound, "gas consumer %s not found", req.Consumer)
	}
	return &types.QueryGasConsumerResponse{
		GasConsumer: gc,
	}, nil
}

func (k Querier) GasConsumers(goCtx context.Context, req *types.QueryGasConsumersRequest) (*types.QueryGasConsumersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)

	keyPrefix := types.GetAllGasConsumersKey()
	gcGetter := func(_, value []byte) types.GasConsumer {
		return types.MustUnmarshalGasConsumer(k.cdc, value)
	}
	gcStore := prefix.NewStore(store, keyPrefix)
	var gasConsumers []types.GasConsumer

	pageRes, err := query.FilteredPaginate(gcStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		gc := gcGetter(key, value)
		if accumulate {
			gasConsumers = append(gasConsumers, gc)
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryGasConsumersResponse{
		GasConsumers: gasConsumers,
		Pagination:   pageRes,
	}, nil
}

func (k Querier) GasConsumersByGasTankID(goCtx context.Context, req *types.QueryGasConsumersByGasTankIDRequest) (*types.QueryGasConsumersByGasTankIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GasTankId == 0 {
		return nil, status.Error(codes.InvalidArgument, "gas tank id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	gt, found := k.GetGasTank(ctx, req.GasTankId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "gas tank with id %d doesn't exist", req.GasTankId)
	}

	tankConsumers := []types.GasConsumersByGasTankIDResponse{}
	overallFeesConsumed := sdk.NewCoin(gt.FeeDenom, sdk.ZeroInt())

	allConsumers := k.GetAllGasConsumers(ctx)
	for _, consumer := range allConsumers {
		for _, consumption := range consumer.Consumptions {
			if consumption.GasTankId == req.GasTankId {
				overallFeesConsumed.Amount = overallFeesConsumed.Amount.Add(consumption.TotalFeesConsumed)
				tankConsumers = append(tankConsumers, types.GasConsumersByGasTankIDResponse{
					Consumer:                   consumer.Consumer,
					IsBlocked:                  consumption.IsBlocked,
					TotalFeeConsumptionAllowed: sdk.NewCoin(gt.FeeDenom, consumption.TotalFeeConsumptionAllowed),
					TotalFeesConsumed:          sdk.NewCoin(gt.FeeDenom, consumption.TotalFeesConsumed),
					Usage:                      consumption.Usage,
				})
				break
			}
		}
	}

	return &types.QueryGasConsumersByGasTankIDResponse{
		GasTankId:           req.GasTankId,
		OverallFeesConsumed: overallFeesConsumed,
		GasConsumers:        tankConsumers,
	}, nil
}

func (k Querier) GasTankIdsForAllUsageIdentifiers(goCtx context.Context, _ *types.QueryGasTankIdsForAllUsageIdentifiersRequest) (*types.QueryGasTankIdsForAllUsageIdentifiersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	usageIdentifierToGasTankIds := []*types.UsageIdentifierToGasTankIds{}
	allusageIdentifierToGasTankIds := k.GetAllUsageIdentifierToGasTankIds(ctx)
	for _, val := range allusageIdentifierToGasTankIds {
		gtids := val
		usageIdentifierToGasTankIds = append(usageIdentifierToGasTankIds, &gtids)
	}
	return &types.QueryGasTankIdsForAllUsageIdentifiersResponse{
		UsageIdentifierToGastankIds: usageIdentifierToGasTankIds,
	}, nil
}
