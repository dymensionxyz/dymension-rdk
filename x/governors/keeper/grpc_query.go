package keeper

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Governors queries all governors that match the given status
func (k Querier) Governors(c context.Context, req *types.QueryGovernorsRequest) (*types.QueryGovernorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// validate the provided status, return all the governors if the status is empty
	if req.Status != "" && !(req.Status == types.Bonded.String() || req.Status == types.Unbonded.String() || req.Status == types.Unbonding.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid governor status %s", req.Status)
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	valStore := prefix.NewStore(store, types.GovernorsKey)

	governors, pageRes, err := query.GenericFilteredPaginate(k.cdc, valStore, req.Pagination, func(key []byte, val *types.Governor) (*types.Governor, error) {
		if req.Status != "" && !strings.EqualFold(val.GetStatus().String(), req.Status) {
			return nil, nil
		}

		return val, nil
	}, func() *types.Governor {
		return &types.Governor{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	vals := types.Governors{}
	for _, val := range governors {
		vals = append(vals, *val)
	}

	return &types.QueryGovernorsResponse{Governors: vals, Pagination: pageRes}, nil
}

// Governor queries governor info for given governor address
func (k Querier) Governor(c context.Context, req *types.QueryGovernorRequest) (*types.QueryGovernorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GovernorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "governor address cannot be empty")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.GovernorAddr)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	governor, found := k.GetGovernor(ctx, valAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "governor %s not found", req.GovernorAddr)
	}

	return &types.QueryGovernorResponse{Governor: governor}, nil
}

// GovernorDelegations queries delegate info for given governor
func (k Querier) GovernorDelegations(c context.Context, req *types.QueryGovernorDelegationsRequest) (*types.QueryGovernorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GovernorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "governor address cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	valStore := prefix.NewStore(store, types.DelegationKey)
	delegations, pageRes, err := query.GenericFilteredPaginate(k.cdc, valStore, req.Pagination, func(key []byte, delegation *types.Delegation) (*types.Delegation, error) {
		valAddr, err := sdk.ValAddressFromBech32(req.GovernorAddr)
		if err != nil {
			return nil, err
		}

		if !delegation.GetGovernorAddr().Equals(valAddr) {
			return nil, nil
		}

		return delegation, nil
	}, func() *types.Delegation {
		return &types.Delegation{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	dels := types.Delegations{}
	for _, d := range delegations {
		dels = append(dels, *d)
	}

	delResponses, err := DelegationsToDelegationResponses(ctx, k.Keeper, dels)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGovernorDelegationsResponse{
		DelegationResponses: delResponses, Pagination: pageRes,
	}, nil
}

// GovernorUnbondingDelegations queries unbonding delegations of a governor
func (k Querier) GovernorUnbondingDelegations(c context.Context, req *types.QueryGovernorUnbondingDelegationsRequest) (*types.QueryGovernorUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GovernorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "governor address cannot be empty")
	}
	var ubds types.UnbondingDelegations
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)

	valAddr, err := sdk.ValAddressFromBech32(req.GovernorAddr)
	if err != nil {
		return nil, err
	}

	srcValPrefix := types.GetUBDsByValIndexKey(valAddr)
	ubdStore := prefix.NewStore(store, srcValPrefix)
	pageRes, err := query.Paginate(ubdStore, req.Pagination, func(key []byte, value []byte) error {
		storeKey := types.GetUBDKeyFromValIndexKey(append(srcValPrefix, key...))
		storeValue := store.Get(storeKey)

		ubd, err := types.UnmarshalUBD(k.cdc, storeValue)
		if err != nil {
			return err
		}
		ubds = append(ubds, ubd)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGovernorUnbondingDelegationsResponse{
		UnbondingResponses: ubds,
		Pagination:         pageRes,
	}, nil
}

// Delegation queries delegate info for given governor delegator pair
func (k Querier) Delegation(c context.Context, req *types.QueryDelegationRequest) (*types.QueryDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.GovernorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "governor address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(req.GovernorAddr)
	if err != nil {
		return nil, err
	}

	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"delegation with delegator %s not found for governor %s",
			req.DelegatorAddr, req.GovernorAddr)
	}

	delResponse, err := DelegationToDelegationResponse(ctx, k.Keeper, delegation)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegationResponse{DelegationResponse: &delResponse}, nil
}

// UnbondingDelegation queries unbonding info for give governor delegator pair
func (k Querier) UnbondingDelegation(c context.Context, req *types.QueryUnbondingDelegationRequest) (*types.QueryUnbondingDelegationResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.GovernorAddr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "governor address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(req.GovernorAddr)
	if err != nil {
		return nil, err
	}

	unbond, found := k.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"unbonding delegation with delegator %s not found for governor %s",
			req.DelegatorAddr, req.GovernorAddr)
	}

	return &types.QueryUnbondingDelegationResponse{Unbond: unbond}, nil
}

// DelegatorDelegations queries all delegations of a give delegator address
func (k Querier) DelegatorDelegations(c context.Context, req *types.QueryDelegatorDelegationsRequest) (*types.QueryDelegatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	var delegations types.Delegations
	ctx := sdk.UnwrapSDKContext(c)

	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.GetDelegationsKey(delAddr))
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	delegationResps, err := DelegationsToDelegationResponses(ctx, k.Keeper, delegations)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorDelegationsResponse{DelegationResponses: delegationResps, Pagination: pageRes}, nil
}

// DelegatorGovernor queries governor info for given delegator governor pair
func (k Querier) DelegatorGovernor(c context.Context, req *types.QueryDelegatorGovernorRequest) (*types.QueryDelegatorGovernorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.GovernorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "governor address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(req.GovernorAddr)
	if err != nil {
		return nil, err
	}

	governor, err := k.GetDelegatorGovernor(ctx, delAddr, valAddr)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorGovernorResponse{Governor: governor}, nil
}

// DelegatorUnbondingDelegations queries all unbonding delegations of a given delegator address
func (k Querier) DelegatorUnbondingDelegations(c context.Context, req *types.QueryDelegatorUnbondingDelegationsRequest) (*types.QueryDelegatorUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	var unbondingDelegations types.UnbondingDelegations
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	unbStore := prefix.NewStore(store, types.GetUBDsKey(delAddr))
	pageRes, err := query.Paginate(unbStore, req.Pagination, func(key []byte, value []byte) error {
		unbond, err := types.UnmarshalUBD(k.cdc, value)
		if err != nil {
			return err
		}
		unbondingDelegations = append(unbondingDelegations, unbond)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorUnbondingDelegationsResponse{
		UnbondingResponses: unbondingDelegations, Pagination: pageRes,
	}, nil
}

// Redelegations queries redelegations of given address
func (k Querier) Redelegations(c context.Context, req *types.QueryRedelegationsRequest) (*types.QueryRedelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	var redels types.Redelegations
	var pageRes *query.PageResponse
	var err error

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	switch {
	case req.DelegatorAddr != "" && req.SrcGovernorAddr != "" && req.DstGovernorAddr != "":
		redels, err = queryRedelegation(ctx, k, req)
	case req.DelegatorAddr == "" && req.SrcGovernorAddr != "" && req.DstGovernorAddr == "":
		redels, pageRes, err = queryRedelegationsFromSrcGovernor(store, k, req)
	default:
		redels, pageRes, err = queryAllRedelegations(store, k, req)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	redelResponses, err := RedelegationsToRedelegationResponses(ctx, k.Keeper, redels)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRedelegationsResponse{RedelegationResponses: redelResponses, Pagination: pageRes}, nil
}

// DelegatorGovernors queries all governors info for given delegator address
func (k Querier) DelegatorGovernors(c context.Context, req *types.QueryDelegatorGovernorsRequest) (*types.QueryDelegatorGovernorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	var governors types.Governors
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	delStore := prefix.NewStore(store, types.GetDelegationsKey(delAddr))
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}

		governor, found := k.GetGovernor(ctx, delegation.GetGovernorAddr())
		if !found {
			return types.ErrNoGovernorFound
		}

		governors = append(governors, governor)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorGovernorsResponse{Governors: governors, Pagination: pageRes}, nil
}

// Pool queries the pool info
func (k Querier) Pool(c context.Context, _ *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	bondDenom := k.BondDenom(ctx)
	bondedPool := k.GetBondedPool(ctx)
	notBondedPool := k.GetNotBondedPool(ctx)

	pool := types.NewPool(
		k.bankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount,
		k.bankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount,
	)

	return &types.QueryPoolResponse{Pool: pool}, nil
}

// Params queries the staking parameters
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func queryRedelegation(ctx sdk.Context, k Querier, req *types.QueryRedelegationsRequest) (redels types.Redelegations, err error) {
	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	srcValAddr, err := sdk.ValAddressFromBech32(req.SrcGovernorAddr)
	if err != nil {
		return nil, err
	}

	dstValAddr, err := sdk.ValAddressFromBech32(req.DstGovernorAddr)
	if err != nil {
		return nil, err
	}

	redel, found := k.GetRedelegation(ctx, delAddr, srcValAddr, dstValAddr)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"redelegation not found for delegator address %s from governor address %s",
			req.DelegatorAddr, req.SrcGovernorAddr)
	}
	redels = []types.Redelegation{redel}

	return redels, err
}

func queryRedelegationsFromSrcGovernor(store sdk.KVStore, k Querier, req *types.QueryRedelegationsRequest) (redels types.Redelegations, res *query.PageResponse, err error) {
	valAddr, err := sdk.ValAddressFromBech32(req.SrcGovernorAddr)
	if err != nil {
		return nil, nil, err
	}

	srcValPrefix := types.GetREDsFromValSrcIndexKey(valAddr)
	redStore := prefix.NewStore(store, srcValPrefix)
	res, err = query.Paginate(redStore, req.Pagination, func(key []byte, value []byte) error {
		storeKey := types.GetREDKeyFromValSrcIndexKey(append(srcValPrefix, key...))
		storeValue := store.Get(storeKey)
		red, err := types.UnmarshalRED(k.cdc, storeValue)
		if err != nil {
			return err
		}
		redels = append(redels, red)
		return nil
	})

	return redels, res, err
}

func queryAllRedelegations(store sdk.KVStore, k Querier, req *types.QueryRedelegationsRequest) (redels types.Redelegations, res *query.PageResponse, err error) {
	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, nil, err
	}

	redStore := prefix.NewStore(store, types.GetREDsKey(delAddr))
	res, err = query.Paginate(redStore, req.Pagination, func(key []byte, value []byte) error {
		redelegation, err := types.UnmarshalRED(k.cdc, value)
		if err != nil {
			return err
		}
		redels = append(redels, redelegation)
		return nil
	})

	return redels, res, err
}
