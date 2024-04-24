package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/governors/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/governors/teststaking"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

func TestNewQuerier(t *testing.T) {
	cdc, app, ctx := createTestInput(t)

	addrs := utils.AddTestAddrs(app, ctx, 500, sdk.NewInt(10000))
	_, addrAcc2 := addrs[0], addrs[1]
	addrVal1, _ := sdk.ValAddress(addrs[0]), sdk.ValAddress(addrs[1])

	// Create Governors
	amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8)}
	var governors [2]types.Governor
	for i, amt := range amts {
		governors[i] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[i]))
		governors[i], _ = governors[i].AddTokensFromDel(amt)
		app.StakingKeeper.SetGovernor(ctx, governors[i])
		app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[i])
	}

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := keeper.NewQuerier(app.StakingKeeper, legacyQuerierCdc.LegacyAmino)

	bz, err := querier(ctx, []string{"other"}, query)
	require.Error(t, err)
	require.Nil(t, bz)

	_, err = querier(ctx, []string{"pool"}, query)
	require.NoError(t, err)

	_, err = querier(ctx, []string{"parameters"}, query)
	require.NoError(t, err)

	queryValParams := types.NewQueryGovernorParams(addrVal1, 0, 0)
	bz, errRes := cdc.MarshalJSON(queryValParams)
	require.NoError(t, errRes)

	query.Path = "/custom/staking/governor"
	query.Data = bz

	_, err = querier(ctx, []string{"governor"}, query)
	require.NoError(t, err)

	_, err = querier(ctx, []string{"governorDelegations"}, query)
	require.NoError(t, err)

	_, err = querier(ctx, []string{"governorUnbondingDelegations"}, query)
	require.NoError(t, err)

	queryDelParams := types.NewQueryDelegatorParams(addrAcc2)
	bz, errRes = cdc.MarshalJSON(queryDelParams)
	require.NoError(t, errRes)

	query.Path = "/custom/staking/governor"
	query.Data = bz

	_, err = querier(ctx, []string{"delegatorDelegations"}, query)
	require.NoError(t, err)

	_, err = querier(ctx, []string{"delegatorUnbondingDelegations"}, query)
	require.NoError(t, err)

	_, err = querier(ctx, []string{"delegatorGovernors"}, query)
	require.NoError(t, err)

	bz, errRes = cdc.MarshalJSON(types.NewQueryRedelegationParams(nil, nil, nil))
	require.NoError(t, errRes)
	query.Data = bz

	_, err = querier(ctx, []string{"redelegations"}, query)
	require.NoError(t, err)
}

func TestQueryParametersPool(t *testing.T) {
	cdc, app, ctx := createTestInput(t)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := keeper.NewQuerier(app.StakingKeeper, legacyQuerierCdc.LegacyAmino)

	bondDenom := sdk.DefaultBondDenom

	res, err := querier(ctx, []string{types.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, err)

	var params types.Params
	errRes := cdc.UnmarshalJSON(res, &params)
	require.NoError(t, errRes)
	require.Equal(t, app.StakingKeeper.GetParams(ctx), params)

	res, err = querier(ctx, []string{types.QueryPool}, abci.RequestQuery{})
	require.NoError(t, err)

	var pool types.Pool
	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	require.NoError(t, cdc.UnmarshalJSON(res, &pool))
	require.Equal(t, app.BankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount, pool.NotBondedTokens)
	require.Equal(t, app.BankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount, pool.BondedTokens)
}

func TestQueryGovernors(t *testing.T) {
	cdc, app, ctx := createTestInput(t)
	params := app.StakingKeeper.GetParams(ctx)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := keeper.NewQuerier(app.StakingKeeper, legacyQuerierCdc.LegacyAmino)

	addrs := utils.AddTestAddrs(app, ctx, 500, app.StakingKeeper.TokensFromConsensusPower(ctx, 10000))

	// Create Governors
	amts := []sdk.Int{sdk.NewInt(8), sdk.NewInt(7)}
	status := []types.BondStatus{types.Unbonded, types.Unbonding}
	var governors [2]types.Governor
	for i, amt := range amts {
		governors[i] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[i]))
		governors[i], _ = governors[i].AddTokensFromDel(amt)
		governors[i] = governors[i].UpdateStatus(status[i])
	}

	app.StakingKeeper.SetGovernor(ctx, governors[0])
	app.StakingKeeper.SetGovernor(ctx, governors[1])

	// Query Governors
	queriedGovernors := app.StakingKeeper.GetGovernors(ctx, params.MaxValidators)
	require.Len(t, queriedGovernors, 3)

	for i, s := range status {
		queryValsParams := types.NewQueryGovernorsParams(1, int(params.MaxValidators), s.String())
		bz, err := cdc.MarshalJSON(queryValsParams)
		require.NoError(t, err)

		req := abci.RequestQuery{
			Path: fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryGovernors),
			Data: bz,
		}

		res, err := querier(ctx, []string{types.QueryGovernors}, req)
		require.NoError(t, err)

		var governorsResp []types.Governor
		err = cdc.UnmarshalJSON(res, &governorsResp)
		require.NoError(t, err)

		require.Equal(t, 1, len(governorsResp))
		require.Equal(t, governors[i].OperatorAddress, governorsResp[0].OperatorAddress)
	}

	// Query each governor
	for _, governor := range governors {
		queryParams := types.NewQueryGovernorParams(governor.GetOperator(), 0, 0)
		bz, err := cdc.MarshalJSON(queryParams)
		require.NoError(t, err)

		query := abci.RequestQuery{
			Path: "/custom/staking/governor",
			Data: bz,
		}
		res, err := querier(ctx, []string{types.QueryGovernor}, query)
		require.NoError(t, err)

		var queriedGovernor types.Governor
		err = cdc.UnmarshalJSON(res, &queriedGovernor)
		require.NoError(t, err)

		require.True(t, governor.Equal(&queriedGovernor))
	}
}

func TestQueryDelegation(t *testing.T) {
	cdc, app, ctx := createTestInput(t)
	params := app.StakingKeeper.GetParams(ctx)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := keeper.NewQuerier(app.StakingKeeper, legacyQuerierCdc.LegacyAmino)

	addrs := utils.AddTestAddrs(app, ctx, 2, app.StakingKeeper.TokensFromConsensusPower(ctx, 10000))
	addrAcc1, addrAcc2 := addrs[0], addrs[1]
	addrVal1, addrVal2 := sdk.ValAddress(addrAcc1), sdk.ValAddress(addrAcc2)

	// Create Governors and Delegation
	val1 := teststaking.NewGovernor(t, addrVal1)
	app.StakingKeeper.SetGovernor(ctx, val1)
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, val1)

	val2 := teststaking.NewGovernor(t, addrVal2)
	app.StakingKeeper.SetGovernor(ctx, val2)
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, val2)

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 20)
	_, err := app.StakingKeeper.Delegate(ctx, addrAcc2, delTokens, types.Unbonded, val1, true)
	require.NoError(t, err)

	// apply TM updates
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, -1)

	// Query Delegator bonded governors
	queryParams := types.NewQueryDelegatorParams(addrAcc2)
	bz, errRes := cdc.MarshalJSON(queryParams)
	require.NoError(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/staking/delegatorGovernors",
		Data: bz,
	}

	delGovernors := app.StakingKeeper.GetDelegatorGovernors(ctx, addrAcc2, params.MaxValidators)

	res, err := querier(ctx, []string{types.QueryDelegatorGovernors}, query)
	require.NoError(t, err)

	var governorsResp types.Governors
	errRes = cdc.UnmarshalJSON(res, &governorsResp)
	require.NoError(t, errRes)

	require.Equal(t, len(delGovernors), len(governorsResp))
	require.ElementsMatch(t, delGovernors, governorsResp)

	// error unknown request
	query.Data = bz[:len(bz)-1]

	_, err = querier(ctx, []string{types.QueryDelegatorGovernors}, query)
	require.Error(t, err)

	// Query bonded governor
	queryBondParams := types.QueryDelegatorGovernorRequest{DelegatorAddr: addrAcc2.String(), GovernorAddr: addrVal1.String()}
	bz, errRes = cdc.MarshalJSON(queryBondParams)
	require.NoError(t, errRes)

	query = abci.RequestQuery{
		Path: "/custom/staking/delegatorGovernor",
		Data: bz,
	}

	res, err = querier(ctx, []string{types.QueryDelegatorGovernor}, query)
	require.NoError(t, err)

	var governor types.Governor
	errRes = cdc.UnmarshalJSON(res, &governor)
	require.NoError(t, errRes)
	require.True(t, governor.Equal(&delGovernors[0]))

	// error unknown request
	query.Data = bz[:len(bz)-1]

	_, err = querier(ctx, []string{types.QueryDelegatorGovernor}, query)
	require.Error(t, err)

	// Query delegation

	query = abci.RequestQuery{
		Path: "/custom/staking/delegation",
		Data: bz,
	}

	delegation, found := app.StakingKeeper.GetDelegation(ctx, addrAcc2, addrVal1)
	require.True(t, found)

	res, err = querier(ctx, []string{types.QueryDelegation}, query)
	require.NoError(t, err)

	var delegationRes stakingtypes.DelegationResponse
	errRes = cdc.UnmarshalJSON(res, &delegationRes)
	require.NoError(t, errRes)

	require.Equal(t, delegation.ValidatorAddress, delegationRes.Delegation.ValidatorAddress)
	require.Equal(t, delegation.DelegatorAddress, delegationRes.Delegation.DelegatorAddress)
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), delegationRes.Balance)

	// Query Delegator Delegations
	bz, errRes = cdc.MarshalJSON(queryParams)
	require.NoError(t, errRes)

	query = abci.RequestQuery{
		Path: "/custom/staking/delegatorDelegations",
		Data: bz,
	}

	res, err = querier(ctx, []string{types.QueryDelegatorDelegations}, query)
	require.NoError(t, err)

	var delegatorDelegations stakingtypes.DelegationResponses
	errRes = cdc.UnmarshalJSON(res, &delegatorDelegations)
	require.NoError(t, errRes)
	require.Len(t, delegatorDelegations, 1)
	require.Equal(t, delegation.ValidatorAddress, delegatorDelegations[0].Delegation.ValidatorAddress)
	require.Equal(t, delegation.DelegatorAddress, delegatorDelegations[0].Delegation.DelegatorAddress)
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), delegatorDelegations[0].Balance)

	// error unknown request
	query.Data = bz[:len(bz)-1]

	_, err = querier(ctx, []string{types.QueryDelegation}, query)
	require.Error(t, err)

	// Query governor delegations
	bz, errRes = cdc.MarshalJSON(types.NewQueryGovernorParams(addrVal1, 1, 100))
	require.NoError(t, errRes)

	query = abci.RequestQuery{
		Path: "custom/staking/governorDelegations",
		Data: bz,
	}

	res, err = querier(ctx, []string{types.QueryGovernorDelegations}, query)
	require.NoError(t, err)

	var delegationsRes stakingtypes.DelegationResponses
	errRes = cdc.UnmarshalJSON(res, &delegationsRes)
	require.NoError(t, errRes)
	require.Len(t, delegatorDelegations, 1)
	require.Equal(t, delegation.ValidatorAddress, delegationsRes[0].Delegation.ValidatorAddress)
	require.Equal(t, delegation.DelegatorAddress, delegationsRes[0].Delegation.DelegatorAddress)
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), delegationsRes[0].Balance)

	// Query unbonding delegation
	unbondingTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	_, err = app.StakingKeeper.Undelegate(ctx, addrAcc2, val1.GetOperator(), sdk.NewDecFromInt(unbondingTokens))
	require.NoError(t, err)

	queryBondParams = types.QueryDelegatorGovernorRequest{DelegatorAddr: addrAcc2.String(), GovernorAddr: addrVal1.String()}
	bz, errRes = cdc.MarshalJSON(queryBondParams)
	require.NoError(t, errRes)

	query = abci.RequestQuery{
		Path: "/custom/staking/unbondingDelegation",
		Data: bz,
	}

	unbond, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrAcc2, addrVal1)
	require.True(t, found)

	res, err = querier(ctx, []string{types.QueryUnbondingDelegation}, query)
	require.NoError(t, err)

	var unbondRes stakingtypes.UnbondingDelegation
	errRes = cdc.UnmarshalJSON(res, &unbondRes)
	require.NoError(t, errRes)

	require.Equal(t, unbond, unbondRes)

	// error unknown request
	query.Data = bz[:len(bz)-1]

	_, err = querier(ctx, []string{types.QueryUnbondingDelegation}, query)
	require.Error(t, err)

	// Query Delegator Unbonding Delegations

	query = abci.RequestQuery{
		Path: "/custom/staking/delegatorUnbondingDelegations",
		Data: bz,
	}

	res, err = querier(ctx, []string{types.QueryDelegatorUnbondingDelegations}, query)
	require.NoError(t, err)

	var delegatorUbds []stakingtypes.UnbondingDelegation
	errRes = cdc.UnmarshalJSON(res, &delegatorUbds)
	require.NoError(t, errRes)
	require.Equal(t, unbond, delegatorUbds[0])

	// error unknown request
	query.Data = bz[:len(bz)-1]

	_, err = querier(ctx, []string{types.QueryDelegatorUnbondingDelegations}, query)
	require.Error(t, err)

	// Query redelegation
	redelegationTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	_, err = app.StakingKeeper.BeginRedelegation(ctx, addrAcc2, val1.GetOperator(), val2.GetOperator(), sdk.NewDecFromInt(redelegationTokens))
	require.NoError(t, err)
	redel, found := app.StakingKeeper.GetRedelegation(ctx, addrAcc2, val1.GetOperator(), val2.GetOperator())
	require.True(t, found)

	bz, errRes = cdc.MarshalJSON(types.NewQueryRedelegationParams(addrAcc2, val1.GetOperator(), val2.GetOperator()))
	require.NoError(t, errRes)

	query = abci.RequestQuery{
		Path: "/custom/staking/redelegations",
		Data: bz,
	}

	res, err = querier(ctx, []string{types.QueryRedelegations}, query)
	require.NoError(t, err)

	var redelRes stakingtypes.RedelegationResponses
	errRes = cdc.UnmarshalJSON(res, &redelRes)
	require.NoError(t, errRes)
	require.Len(t, redelRes, 1)
	require.Equal(t, redel.DelegatorAddress, redelRes[0].Redelegation.DelegatorAddress)
	require.Equal(t, redel.ValidatorSrcAddress, redelRes[0].Redelegation.ValidatorSrcAddress)
	require.Equal(t, redel.ValidatorDstAddress, redelRes[0].Redelegation.ValidatorDstAddress)
	require.Len(t, redel.Entries, len(redelRes[0].Entries))
}

func TestQueryGovernorDelegations_Pagination(t *testing.T) {
	cases := []struct {
		page            int
		limit           int
		expectedResults int
	}{
		{
			page:            1,
			limit:           75,
			expectedResults: 75,
		},
		{
			page:            2,
			limit:           75,
			expectedResults: 25,
		},
		{
			page:            1,
			limit:           100,
			expectedResults: 100,
		},
	}

	cdc, app, ctx := createTestInput(t)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := keeper.NewQuerier(app.StakingKeeper, legacyQuerierCdc.LegacyAmino)

	addrs := utils.AddTestAddrs(app, ctx, 100, app.StakingKeeper.TokensFromConsensusPower(ctx, 10000))
	valAddress := sdk.ValAddress(addrs[0])

	val1 := teststaking.NewGovernor(t, valAddress)
	app.StakingKeeper.SetGovernor(ctx, val1)
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, val1)

	// Create Governors and Delegation
	for _, addr := range addrs {
		governor, found := app.StakingKeeper.GetGovernor(ctx, valAddress)
		if !found {
			t.Error("expected governor not found")
		}

		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 20)
		_, err := app.StakingKeeper.Delegate(ctx, addr, delTokens, types.Unbonded, governor, true)
		require.NoError(t, err)
	}

	// apply TM updates
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, -1)

	for _, c := range cases {
		// Query Delegator bonded governors
		queryParams := types.NewQueryDelegatorParams(addrs[0])
		bz, errRes := cdc.MarshalJSON(queryParams)
		require.NoError(t, errRes)

		// Query valAddress delegations
		bz, errRes = cdc.MarshalJSON(types.NewQueryGovernorParams(valAddress, c.page, c.limit))
		require.NoError(t, errRes)

		query := abci.RequestQuery{
			Path: "custom/staking/governorDelegations",
			Data: bz,
		}

		res, err := querier(ctx, []string{types.QueryGovernorDelegations}, query)
		require.NoError(t, err)

		var delegationsRes stakingtypes.DelegationResponses
		errRes = cdc.UnmarshalJSON(res, &delegationsRes)
		require.NoError(t, errRes)
		require.Len(t, delegationsRes, c.expectedResults)
	}

	// Undelegate
	for _, addr := range addrs {
		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 20)
		_, err := app.StakingKeeper.Undelegate(ctx, addr, val1.GetOperator(), sdk.NewDecFromInt(delTokens))
		require.NoError(t, err)
	}

	// apply TM updates
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, -1)

	for _, c := range cases {
		// Query Unbonding delegations with pagination.
		queryParams := types.NewQueryDelegatorParams(addrs[0])
		bz, errRes := cdc.MarshalJSON(queryParams)
		require.NoError(t, errRes)

		bz, errRes = cdc.MarshalJSON(types.NewQueryGovernorParams(valAddress, c.page, c.limit))
		require.NoError(t, errRes)
		query := abci.RequestQuery{
			Data: bz,
		}

		unbondingDelegations := stakingtypes.UnbondingDelegations{}
		res, err := querier(ctx, []string{types.QueryGovernorUnbondingDelegations}, query)
		require.NoError(t, err)

		errRes = cdc.UnmarshalJSON(res, &unbondingDelegations)
		require.NoError(t, errRes)
		require.Len(t, unbondingDelegations, c.expectedResults)
	}
}

func TestQueryRedelegations(t *testing.T) {
	cdc, app, ctx := createTestInput(t)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := keeper.NewQuerier(app.StakingKeeper, legacyQuerierCdc.LegacyAmino)

	addrs := utils.AddTestAddrs(app, ctx, 2, app.StakingKeeper.TokensFromConsensusPower(ctx, 10000))
	addrAcc1, addrAcc2 := addrs[0], addrs[1]
	addrVal1, addrVal2 := sdk.ValAddress(addrAcc1), sdk.ValAddress(addrAcc2)

	// Create Governors and Delegation
	val1 := teststaking.NewGovernor(t, addrVal1)
	val2 := teststaking.NewGovernor(t, addrVal2)
	app.StakingKeeper.SetGovernor(ctx, val1)
	app.StakingKeeper.SetGovernor(ctx, val2)

	delAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 100)
	_, err := app.StakingKeeper.Delegate(ctx, addrAcc2, delAmount, types.Unbonded, val1, true)
	require.NoError(t, err)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, -1)

	rdAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 20)
	_, err = app.StakingKeeper.BeginRedelegation(ctx, addrAcc2, val1.GetOperator(), val2.GetOperator(), sdk.NewDecFromInt(rdAmount))
	require.NoError(t, err)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, -1)

	redel, found := app.StakingKeeper.GetRedelegation(ctx, addrAcc2, val1.GetOperator(), val2.GetOperator())
	require.True(t, found)

	// delegator redelegations
	queryDelegatorParams := types.NewQueryDelegatorParams(addrAcc2)
	bz, errRes := cdc.MarshalJSON(queryDelegatorParams)
	require.NoError(t, errRes)

	query := abci.RequestQuery{
		Path: "/custom/staking/redelegations",
		Data: bz,
	}

	res, err := querier(ctx, []string{types.QueryRedelegations}, query)
	require.NoError(t, err)

	var redelRes stakingtypes.RedelegationResponses
	errRes = cdc.UnmarshalJSON(res, &redelRes)
	require.NoError(t, errRes)
	require.Len(t, redelRes, 1)
	require.Equal(t, redel.DelegatorAddress, redelRes[0].Redelegation.DelegatorAddress)
	require.Equal(t, redel.ValidatorSrcAddress, redelRes[0].Redelegation.ValidatorSrcAddress)
	require.Equal(t, redel.ValidatorDstAddress, redelRes[0].Redelegation.ValidatorDstAddress)
	require.Len(t, redel.Entries, len(redelRes[0].Entries))

	// governor redelegations
	queryGovernorParams := types.NewQueryGovernorParams(val1.GetOperator(), 0, 0)
	bz, errRes = cdc.MarshalJSON(queryGovernorParams)
	require.NoError(t, errRes)

	query = abci.RequestQuery{
		Path: "/custom/staking/redelegations",
		Data: bz,
	}

	res, err = querier(ctx, []string{types.QueryRedelegations}, query)
	require.NoError(t, err)

	errRes = cdc.UnmarshalJSON(res, &redelRes)
	require.NoError(t, errRes)
	require.Len(t, redelRes, 1)
	require.Equal(t, redel.DelegatorAddress, redelRes[0].Redelegation.DelegatorAddress)
	require.Equal(t, redel.ValidatorSrcAddress, redelRes[0].Redelegation.ValidatorSrcAddress)
	require.Equal(t, redel.ValidatorDstAddress, redelRes[0].Redelegation.ValidatorDstAddress)
	require.Len(t, redel.Entries, len(redelRes[0].Entries))
}

func TestQueryUnbondingDelegation(t *testing.T) {
	cdc, app, ctx := createTestInput(t)
	legacyQuerierCdc := codec.NewAminoCodec(app.LegacyAmino())
	querier := keeper.NewQuerier(app.StakingKeeper, legacyQuerierCdc.LegacyAmino)

	addrs := utils.AddTestAddrs(app, ctx, 2, app.StakingKeeper.TokensFromConsensusPower(ctx, 10000))
	addrAcc1, addrAcc2 := addrs[0], addrs[1]
	addrVal1 := sdk.ValAddress(addrAcc1)

	// Create Governors and Delegation
	val1 := teststaking.NewGovernor(t, addrVal1)
	app.StakingKeeper.SetGovernor(ctx, val1)

	// delegate
	delAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 100)
	_, err := app.StakingKeeper.Delegate(ctx, addrAcc1, delAmount, types.Unbonded, val1, true)
	require.NoError(t, err)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, -1)

	// undelegate
	undelAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 20)
	_, err = app.StakingKeeper.Undelegate(ctx, addrAcc1, val1.GetOperator(), sdk.NewDecFromInt(undelAmount))
	require.NoError(t, err)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, -1)

	_, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrAcc1, val1.GetOperator())
	require.True(t, found)

	//
	// found: query unbonding delegation by delegator and governor
	//
	queryGovernorParams := types.QueryDelegatorGovernorRequest{DelegatorAddr: addrAcc1.String(), GovernorAddr: val1.GetOperator().String()}
	bz, errRes := cdc.MarshalJSON(queryGovernorParams)
	require.NoError(t, errRes)
	query := abci.RequestQuery{
		Path: "/custom/staking/unbondingDelegation",
		Data: bz,
	}
	res, err := querier(ctx, []string{types.QueryUnbondingDelegation}, query)
	require.NoError(t, err)
	require.NotNil(t, res)
	var ubDel stakingtypes.UnbondingDelegation
	require.NoError(t, cdc.UnmarshalJSON(res, &ubDel))
	require.Equal(t, addrAcc1.String(), ubDel.DelegatorAddress)
	require.Equal(t, val1.OperatorAddress, ubDel.ValidatorAddress)
	require.Equal(t, 1, len(ubDel.Entries))

	//
	// not found: query unbonding delegation by delegator and governor
	//
	queryGovernorParams = types.QueryDelegatorGovernorRequest{DelegatorAddr: addrAcc2.String(), GovernorAddr: val1.GetOperator().String()}
	bz, errRes = cdc.MarshalJSON(queryGovernorParams)
	require.NoError(t, errRes)
	query = abci.RequestQuery{
		Path: "/custom/staking/unbondingDelegation",
		Data: bz,
	}
	_, err = querier(ctx, []string{types.QueryUnbondingDelegation}, query)
	require.Error(t, err)

	//
	// found: query unbonding delegation by delegator and governor
	//
	queryDelegatorParams := types.NewQueryDelegatorParams(addrAcc1)
	bz, errRes = cdc.MarshalJSON(queryDelegatorParams)
	require.NoError(t, errRes)
	query = abci.RequestQuery{
		Path: "/custom/staking/delegatorUnbondingDelegations",
		Data: bz,
	}
	res, err = querier(ctx, []string{types.QueryDelegatorUnbondingDelegations}, query)
	require.NoError(t, err)
	require.NotNil(t, res)
	var ubDels []stakingtypes.UnbondingDelegation
	require.NoError(t, cdc.UnmarshalJSON(res, &ubDels))
	require.Equal(t, 1, len(ubDels))
	require.Equal(t, addrAcc1.String(), ubDels[0].DelegatorAddress)
	require.Equal(t, val1.OperatorAddress, ubDels[0].ValidatorAddress)

	//
	// not found: query unbonding delegation by delegator and governor
	//
	queryDelegatorParams = types.NewQueryDelegatorParams(addrAcc2)
	bz, errRes = cdc.MarshalJSON(queryDelegatorParams)
	require.NoError(t, errRes)
	query = abci.RequestQuery{
		Path: "/custom/staking/delegatorUnbondingDelegations",
		Data: bz,
	}
	res, err = querier(ctx, []string{types.QueryDelegatorUnbondingDelegations}, query)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NoError(t, cdc.UnmarshalJSON(res, &ubDels))
	require.Equal(t, 0, len(ubDels))
}
