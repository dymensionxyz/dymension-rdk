package keeper

import (
	"errors"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

// creates a querier for staking REST endpoints
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGovernors:
			return queryGovernors(ctx, req, k, legacyQuerierCdc)

		case types.QueryGovernor:
			return queryGovernor(ctx, req, k, legacyQuerierCdc)

		case types.QueryGovernorDelegations:
			return queryGovernorDelegations(ctx, req, k, legacyQuerierCdc)

		case types.QueryGovernorUnbondingDelegations:
			return queryGovernorUnbondingDelegations(ctx, req, k, legacyQuerierCdc)

		case types.QueryDelegation:
			return queryDelegation(ctx, req, k, legacyQuerierCdc)

		case types.QueryUnbondingDelegation:
			return queryUnbondingDelegation(ctx, req, k, legacyQuerierCdc)

		case types.QueryDelegatorDelegations:
			return queryDelegatorDelegations(ctx, req, k, legacyQuerierCdc)

		case types.QueryDelegatorUnbondingDelegations:
			return queryDelegatorUnbondingDelegations(ctx, req, k, legacyQuerierCdc)

		case types.QueryRedelegations:
			return queryRedelegations(ctx, req, k, legacyQuerierCdc)

		case types.QueryDelegatorGovernors:
			return queryDelegatorGovernors(ctx, req, k, legacyQuerierCdc)

		case types.QueryDelegatorGovernor:
			return queryDelegatorGovernor(ctx, req, k, legacyQuerierCdc)

		case types.QueryPool:
			return queryPool(ctx, k, legacyQuerierCdc)

		case types.QueryParameters:
			return queryParameters(ctx, k, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func queryGovernors(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryGovernorsParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	governors := k.GetAllGovernors(ctx)
	filteredVals := make(types.Governors, 0, len(governors))

	for _, val := range governors {
		if strings.EqualFold(val.GetStatus().String(), params.Status) {
			filteredVals = append(filteredVals, val)
		}
	}

	start, end := client.Paginate(len(filteredVals), params.Page, params.Limit, int(k.GetParams(ctx).MaxValidators))
	if start < 0 || end < 0 {
		filteredVals = []types.Governor{}
	} else {
		filteredVals = filteredVals[start:end]
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, filteredVals)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryGovernor(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryGovernorParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	governor, found := k.GetGovernor(ctx, params.GovernorAddr)
	if !found {
		return nil, types.ErrNoGovernorFound
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, governor)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryGovernorDelegations(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryGovernorParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delegations := k.GetGovernorDelegations(ctx, params.GovernorAddr)

	start, end := client.Paginate(len(delegations), params.Page, params.Limit, int(k.GetParams(ctx).MaxValidators))
	if start < 0 || end < 0 {
		delegations = []stakingtypes.Delegation{}
	} else {
		delegations = delegations[start:end]
	}

	delegationResps, err := DelegationsToDelegationResponses(ctx, k, delegations)
	if err != nil {
		return nil, err
	}

	if delegationResps == nil {
		delegationResps = stakingtypes.DelegationResponses{}
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, delegationResps)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryGovernorUnbondingDelegations(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryGovernorParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	unbonds := k.GetUnbondingDelegationsFromGovernor(ctx, params.GovernorAddr)
	if unbonds == nil {
		unbonds = stakingtypes.UnbondingDelegations{}
	}

	start, end := client.Paginate(len(unbonds), params.Page, params.Limit, int(k.GetParams(ctx).MaxValidators))
	if start < 0 || end < 0 {
		unbonds = stakingtypes.UnbondingDelegations{}
	} else {
		unbonds = unbonds[start:end]
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, unbonds)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDelegatorDelegations(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryDelegatorParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delegations := k.GetAllDelegatorDelegations(ctx, params.DelegatorAddr)
	delegationResps, err := DelegationsToDelegationResponses(ctx, k, delegations)
	if err != nil {
		return nil, err
	}

	if delegationResps == nil {
		delegationResps = stakingtypes.DelegationResponses{}
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, delegationResps)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDelegatorUnbondingDelegations(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryDelegatorParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	unbondingDelegations := k.GetAllUnbondingDelegations(ctx, params.DelegatorAddr)
	if unbondingDelegations == nil {
		unbondingDelegations = stakingtypes.UnbondingDelegations{}
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, unbondingDelegations)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDelegatorGovernors(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryDelegatorParams

	stakingParams := k.GetParams(ctx)

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	governors := k.GetDelegatorGovernors(ctx, params.DelegatorAddr, stakingParams.MaxValidators)
	if governors == nil {
		governors = types.Governors{}
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, governors)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDelegatorGovernor(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryDelegatorGovernorRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delAddr, err := sdk.AccAddressFromBech32(params.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(params.GovernorAddr)
	if err != nil {
		return nil, err
	}

	governor, err := k.GetDelegatorGovernor(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, governor)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryDelegation(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryDelegatorGovernorRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delAddr, err := sdk.AccAddressFromBech32(params.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(params.GovernorAddr)
	if err != nil {
		return nil, err
	}

	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, types.ErrNoDelegation
	}

	delegationResp, err := DelegationToDelegationResponse(ctx, k, delegation)
	if err != nil {
		return nil, err
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, delegationResp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryUnbondingDelegation(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryDelegatorGovernorRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delAddr, err := sdk.AccAddressFromBech32(params.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	valAddr, err := sdk.ValAddressFromBech32(params.GovernorAddr)
	if err != nil {
		return nil, err
	}

	unbond, found := k.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, types.ErrNoUnbondingDelegation
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, unbond)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryRedelegations(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryRedelegationParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var redels []stakingtypes.Redelegation

	switch {
	case !params.DelegatorAddr.Empty() && !params.SrcGovernorAddr.Empty() && !params.DstGovernorAddr.Empty():
		redel, found := k.GetRedelegation(ctx, params.DelegatorAddr, params.SrcGovernorAddr, params.DstGovernorAddr)
		if !found {
			return nil, types.ErrNoRedelegation
		}

		redels = []stakingtypes.Redelegation{redel}
	case params.DelegatorAddr.Empty() && !params.SrcGovernorAddr.Empty() && params.DstGovernorAddr.Empty():
		redels = k.GetRedelegationsFromSrcGovernor(ctx, params.SrcGovernorAddr)
	default:
		redels = k.GetAllRedelegations(ctx, params.DelegatorAddr, params.SrcGovernorAddr, params.DstGovernorAddr)
	}

	redelResponses, err := RedelegationsToRedelegationResponses(ctx, k, redels)
	if err != nil {
		return nil, err
	}

	if redelResponses == nil {
		redelResponses = stakingtypes.RedelegationResponses{}
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, redelResponses)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
func queryPool(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	bondDenom := k.BondDenom(ctx)
	bondedPool := k.GetBondedPool(ctx)
	notBondedPool := k.GetNotBondedPool(ctx)

	if bondedPool == nil || notBondedPool == nil {
		return nil, errors.New("pool accounts haven't been set")
	}

	pool := types.NewPool(
		k.bankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount,
		k.bankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount,
	)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, pool)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryParameters(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// util

func DelegationToDelegationResponse(ctx sdk.Context, k Keeper, del stakingtypes.Delegation) (stakingtypes.DelegationResponse, error) {
	val, found := k.GetGovernor(ctx, del.GetValidatorAddr())
	if !found {
		return stakingtypes.DelegationResponse{}, types.ErrNoGovernorFound
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
	if err != nil {
		return stakingtypes.DelegationResponse{}, err
	}

	return stakingtypes.NewDelegationResp(
		delegatorAddress,
		del.GetValidatorAddr(),
		del.Shares,
		sdk.NewCoin(k.BondDenom(ctx), val.TokensFromShares(del.Shares).TruncateInt()),
	), nil
}

func DelegationsToDelegationResponses(
	ctx sdk.Context, k Keeper, delegations stakingtypes.Delegations,
) (stakingtypes.DelegationResponses, error) {
	resp := make(stakingtypes.DelegationResponses, len(delegations))

	for i, del := range delegations {
		delResp, err := DelegationToDelegationResponse(ctx, k, del)
		if err != nil {
			return nil, err
		}

		resp[i] = delResp
	}

	return resp, nil
}

func RedelegationsToRedelegationResponses(
	ctx sdk.Context, k Keeper, redels stakingtypes.Redelegations,
) (stakingtypes.RedelegationResponses, error) {
	resp := make(stakingtypes.RedelegationResponses, len(redels))

	for i, redel := range redels {
		valSrcAddr, err := sdk.ValAddressFromBech32(redel.ValidatorSrcAddress)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := sdk.ValAddressFromBech32(redel.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}

		delegatorAddress := sdk.MustAccAddressFromBech32(redel.DelegatorAddress)

		val, found := k.GetGovernor(ctx, valDstAddr)
		if !found {
			return nil, types.ErrNoGovernorFound
		}

		entryResponses := make([]stakingtypes.RedelegationEntryResponse, len(redel.Entries))
		for j, entry := range redel.Entries {
			entryResponses[j] = stakingtypes.NewRedelegationEntryResponse(
				entry.CreationHeight,
				entry.CompletionTime,
				entry.SharesDst,
				entry.InitialBalance,
				val.TokensFromShares(entry.SharesDst).TruncateInt(),
			)
		}

		resp[i] = stakingtypes.NewRedelegationResponse(
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
			entryResponses,
		)
	}

	return resp, nil
}
