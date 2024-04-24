package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/governors/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/governors/teststaking"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryGovernors() {
	queryClient, vals := suite.queryClient, suite.vals
	var req *types.QueryGovernorsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		numVals  int
		hasNext  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryGovernorsRequest{}
			},
			true,

			len(vals) + 1, // +1 governor from genesis state
			false,
		},
		{
			"empty status returns all the governors",
			func() {
				req = &types.QueryGovernorsRequest{Status: ""}
			},
			true,
			len(vals) + 1, // +1 governor from genesis state
			false,
		},
		{
			"invalid request",
			func() {
				req = &types.QueryGovernorsRequest{Status: "test"}
			},
			false,
			0,
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryGovernorsRequest{
					Status:     types.Bonded.String(),
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			1,
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			valsResp, err := queryClient.Governors(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.NotNil(valsResp)
				suite.Equal(tc.numVals, len(valsResp.Governors))
				suite.Equal(uint64(len(vals))+1, valsResp.Pagination.Total) // +1 governor from genesis state

				if tc.hasNext {
					suite.NotNil(valsResp.Pagination.NextKey)
				} else {
					suite.Nil(valsResp.Pagination.NextKey)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryGovernor() {
	app, ctx, queryClient, vals := suite.app, suite.ctx, suite.queryClient, suite.vals
	governor, found := app.StakingKeeper.GetGovernor(ctx, vals[0].GetOperator())
	suite.True(found)
	var req *types.QueryGovernorRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryGovernorRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryGovernorRequest{GovernorAddr: vals[0].OperatorAddress}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Governor(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.True(governor.Equal(&res.Governor))
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegatorGovernors() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs
	params := app.StakingKeeper.GetParams(ctx)
	delGovernors := app.StakingKeeper.GetDelegatorGovernors(ctx, addrs[0], params.MaxValidators)
	var req *types.QueryDelegatorGovernorsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorGovernorsRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorGovernorsRequest{
					DelegatorAddr: addrs[0].String(),
					Pagination:    &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorGovernors(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.Equal(1, len(res.Governors))
				suite.NotNil(res.Pagination.NextKey)
				suite.Equal(uint64(len(delGovernors)), res.Pagination.Total)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegatorGovernor() {
	queryClient, addrs, vals := suite.queryClient, suite.addrs, suite.vals
	addr := addrs[1]
	addrVal, addrVal1 := vals[0].OperatorAddress, vals[1].OperatorAddress
	var req *types.QueryDelegatorGovernorRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorGovernorRequest{}
			},
			false,
		},
		{
			"invalid delegator, governor pair",
			func() {
				req = &types.QueryDelegatorGovernorRequest{
					DelegatorAddr: addr.String(),
					GovernorAddr:  addrVal,
				}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorGovernorRequest{
					DelegatorAddr: addr.String(),
					GovernorAddr:  addrVal1,
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorGovernor(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.Equal(addrVal1, res.Governor.OperatorAddress)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegation() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc, addrAcc1 := addrs[0], addrs[1]
	addrVal := vals[0].OperatorAddress
	valAddr, err := sdk.ValAddressFromBech32(addrVal)
	suite.NoError(err)
	delegation, found := app.StakingKeeper.GetDelegation(ctx, addrAcc, valAddr)
	suite.True(found)
	var req *types.QueryDelegationRequest

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegationRequest{}
			},
			false,
		},
		{
			"invalid governor, delegator pair",
			func() {
				req = &types.QueryDelegationRequest{
					DelegatorAddr: addrAcc1.String(),
					GovernorAddr:  addrVal,
				}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegationRequest{DelegatorAddr: addrAcc.String(), GovernorAddr: addrVal}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Delegation(gocontext.Background(), req)
			if tc.expPass {
				suite.Equal(delegation.ValidatorAddress, res.DelegationResponse.Delegation.ValidatorAddress)
				suite.Equal(delegation.DelegatorAddress, res.DelegationResponse.Delegation.DelegatorAddress)
				suite.Equal(sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), res.DelegationResponse.Balance)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegatorDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc := addrs[0]
	addrVal1 := vals[0].OperatorAddress
	valAddr, err := sdk.ValAddressFromBech32(addrVal1)
	suite.NoError(err)
	delegation, found := app.StakingKeeper.GetDelegation(ctx, addrAcc, valAddr)
	suite.True(found)
	var req *types.QueryDelegatorDelegationsRequest

	testCases := []struct {
		msg       string
		malleate  func()
		onSuccess func(suite *KeeperTestSuite, response *types.QueryDelegatorDelegationsResponse)
		expErr    bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorDelegationsRequest{}
			},
			func(suite *KeeperTestSuite, response *types.QueryDelegatorDelegationsResponse) {},
			true,
		},
		{
			"valid request with no delegations",
			func() {
				req = &types.QueryDelegatorDelegationsRequest{DelegatorAddr: addrs[4].String()}
			},
			func(suite *KeeperTestSuite, response *types.QueryDelegatorDelegationsResponse) {
				suite.Equal(uint64(0), response.Pagination.Total)
				suite.Len(response.DelegationResponses, 0)
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorDelegationsRequest{
					DelegatorAddr: addrAcc.String(),
					Pagination:    &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			func(suite *KeeperTestSuite, response *types.QueryDelegatorDelegationsResponse) {
				suite.Equal(uint64(2), response.Pagination.Total)
				suite.Len(response.DelegationResponses, 1)
				suite.Equal(sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), response.DelegationResponses[0].Balance)
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorDelegations(gocontext.Background(), req)
			if tc.expErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
				tc.onSuccess(suite, res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryGovernorDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc := addrs[0]
	addrVal1 := vals[1].OperatorAddress
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	addrVal2 := valAddrs[4]
	valAddr, err := sdk.ValAddressFromBech32(addrVal1)
	suite.NoError(err)
	delegation, found := app.StakingKeeper.GetDelegation(ctx, addrAcc, valAddr)
	suite.True(found)

	var req *types.QueryGovernorDelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryGovernorDelegationsRequest{}
			},
			false,
			true,
		},
		{
			"invalid governor delegator pair",
			func() {
				req = &types.QueryGovernorDelegationsRequest{GovernorAddr: addrVal2.String()}
			},
			false,
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryGovernorDelegationsRequest{
					GovernorAddr: addrVal1,
					Pagination:   &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.GovernorDelegations(gocontext.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.DelegationResponses, 1)
				suite.NotNil(res.Pagination.NextKey)
				suite.Equal(uint64(2), res.Pagination.Total)
				suite.Equal(addrVal1, res.DelegationResponses[0].Delegation.ValidatorAddress)
				suite.Equal(sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), res.DelegationResponses[0].Balance)
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.DelegationResponses)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryUnbondingDelegation() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc2 := addrs[1]
	addrVal2 := vals[1].OperatorAddress

	unbondingTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	valAddr, err1 := sdk.ValAddressFromBech32(addrVal2)
	suite.NoError(err1)
	_, err := app.StakingKeeper.Undelegate(ctx, addrAcc2, valAddr, sdk.NewDecFromInt(unbondingTokens))
	suite.NoError(err)

	unbond, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrAcc2, valAddr)
	suite.True(found)
	var req *types.QueryUnbondingDelegationRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryUnbondingDelegationRequest{}
			},
			false,
		},
		{
			"invalid request",
			func() {
				req = &types.QueryUnbondingDelegationRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryUnbondingDelegationRequest{
					DelegatorAddr: addrAcc2.String(), GovernorAddr: addrVal2,
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.UnbondingDelegation(gocontext.Background(), req)
			if tc.expPass {
				suite.NotNil(res)
				suite.Equal(unbond, res.Unbond)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegatorUnbondingDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc, addrAcc1 := addrs[0], addrs[1]
	addrVal, addrVal2 := vals[0].OperatorAddress, vals[1].OperatorAddress

	unbondingTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	valAddr1, err1 := sdk.ValAddressFromBech32(addrVal)
	suite.NoError(err1)
	_, err := app.StakingKeeper.Undelegate(ctx, addrAcc, valAddr1, sdk.NewDecFromInt(unbondingTokens))
	suite.NoError(err)
	valAddr2, err1 := sdk.ValAddressFromBech32(addrVal2)
	suite.NoError(err1)
	_, err = app.StakingKeeper.Undelegate(ctx, addrAcc, valAddr2, sdk.NewDecFromInt(unbondingTokens))
	suite.NoError(err)

	unbond, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrAcc, valAddr1)
	suite.True(found)
	var req *types.QueryDelegatorUnbondingDelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorUnbondingDelegationsRequest{}
			},
			false,
			true,
		},
		{
			"invalid request",
			func() {
				req = &types.QueryDelegatorUnbondingDelegationsRequest{DelegatorAddr: addrAcc1.String()}
			},
			false,
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorUnbondingDelegationsRequest{
					DelegatorAddr: addrAcc.String(),
					Pagination:    &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorUnbondingDelegations(gocontext.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.NotNil(res.Pagination.NextKey)
				suite.Equal(uint64(2), res.Pagination.Total)
				suite.Len(res.UnbondingResponses, 1)
				suite.Equal(unbond, res.UnbondingResponses[0])
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.UnbondingResponses)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryPoolParameters() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	bondDenom := sdk.DefaultBondDenom

	// Query pool
	res, err := queryClient.Pool(gocontext.Background(), &types.QueryPoolRequest{})
	suite.NoError(err)
	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	suite.Equal(app.BankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount, res.Pool.NotBondedTokens)
	suite.Equal(app.BankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount, res.Pool.BondedTokens)

	// Query Params
	resp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.NoError(err)
	suite.Equal(app.StakingKeeper.GetParams(ctx), resp.Params)
}

func (suite *KeeperTestSuite) TestGRPCQueryRedelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals

	addrAcc, addrAcc1 := addrs[0], addrs[1]
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	val1, val2, val3, val4 := vals[0], vals[1], valAddrs[3], valAddrs[4]
	delAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 1)
	_, err := app.StakingKeeper.Delegate(ctx, addrAcc1, delAmount, types.Unbonded, val1, true)
	suite.NoError(err)
	applyGovernorsSetUpdates(suite.T(), ctx, app.StakingKeeper, -1)

	rdAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 1)
	_, err = app.StakingKeeper.BeginRedelegation(ctx, addrAcc1, val1.GetOperator(), val2.GetOperator(), sdk.NewDecFromInt(rdAmount))
	suite.NoError(err)
	applyGovernorsSetUpdates(suite.T(), ctx, app.StakingKeeper, -1)

	redel, found := app.StakingKeeper.GetRedelegation(ctx, addrAcc1, val1.GetOperator(), val2.GetOperator())
	suite.True(found)

	var req *types.QueryRedelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"request redelegations for non existent addr",
			func() {
				req = &types.QueryRedelegationsRequest{DelegatorAddr: addrAcc.String()}
			},
			false,
			false,
		},
		{
			"request redelegations with non existent pairs",
			func() {
				req = &types.QueryRedelegationsRequest{
					DelegatorAddr: addrAcc.String(), SrcGovernorAddr: val3.String(),
					DstGovernorAddr: val4.String(),
				}
			},
			false,
			true,
		},
		{
			"request redelegations with delegatoraddr, sourceValAddr, destValAddr",
			func() {
				req = &types.QueryRedelegationsRequest{
					DelegatorAddr: addrAcc1.String(), SrcGovernorAddr: val1.OperatorAddress,
					DstGovernorAddr: val2.OperatorAddress, Pagination: &query.PageRequest{},
				}
			},
			true,
			false,
		},
		{
			"request redelegations with delegatoraddr and sourceValAddr",
			func() {
				req = &types.QueryRedelegationsRequest{
					DelegatorAddr: addrAcc1.String(), SrcGovernorAddr: val1.OperatorAddress,
					Pagination: &query.PageRequest{},
				}
			},
			true,
			false,
		},
		{
			"query redelegations with sourceValAddr only",
			func() {
				req = &types.QueryRedelegationsRequest{
					SrcGovernorAddr: val1.GetOperator().String(),
					Pagination:      &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Redelegations(gocontext.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.RedelegationResponses, len(redel.Entries))
				suite.Equal(redel.DelegatorAddress, res.RedelegationResponses[0].Redelegation.DelegatorAddress)
				suite.Equal(redel.ValidatorSrcAddress, res.RedelegationResponses[0].Redelegation.ValidatorSrcAddress)
				suite.Equal(redel.ValidatorDstAddress, res.RedelegationResponses[0].Redelegation.ValidatorDstAddress)
				suite.Len(redel.Entries, len(res.RedelegationResponses[0].Entries))
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.RedelegationResponses)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryGovernorUnbondingDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc1, _ := addrs[0], addrs[1]
	val1 := vals[0]

	// undelegate
	undelAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	_, err := app.StakingKeeper.Undelegate(ctx, addrAcc1, val1.GetOperator(), sdk.NewDecFromInt(undelAmount))
	suite.NoError(err)
	applyGovernorsSetUpdates(suite.T(), ctx, app.StakingKeeper, -1)

	var req *types.QueryGovernorUnbondingDelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryGovernorUnbondingDelegationsRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryGovernorUnbondingDelegationsRequest{
					GovernorAddr: val1.GetOperator().String(),
					Pagination:   &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.GovernorUnbondingDelegations(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.Equal(uint64(1), res.Pagination.Total)
				suite.Equal(1, len(res.UnbondingResponses))
				suite.Equal(res.UnbondingResponses[0].ValidatorAddress, val1.OperatorAddress)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func createGovernors(t *testing.T, ctx sdk.Context, app *app.App, powers []int64) ([]sdk.AccAddress, []sdk.ValAddress, []types.Governor) {
	addrs := utils.AddTestAddrs(app, ctx, 5, app.StakingKeeper.TokensFromConsensusPower(ctx, 300))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	cdc := simapp.MakeTestEncodingConfig().Codec
	app.StakingKeeper = keeper.NewKeeper(
		cdc,
		app.GetKey(types.StoreKey),
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(types.ModuleName),
	)

	val1 := teststaking.NewGovernor(t, valAddrs[0])
	val2 := teststaking.NewGovernor(t, valAddrs[1])
	vals := []types.Governor{val1, val2}

	app.StakingKeeper.SetGovernor(ctx, val1)
	app.StakingKeeper.SetGovernor(ctx, val2)
	app.StakingKeeper.SetNewGovernorByPowerIndex(ctx, val1)
	app.StakingKeeper.SetNewGovernorByPowerIndex(ctx, val2)

	_, err := app.StakingKeeper.Delegate(ctx, addrs[0], app.StakingKeeper.TokensFromConsensusPower(ctx, powers[0]), types.Unbonded, val1, true)
	require.NoError(t, err)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[1], app.StakingKeeper.TokensFromConsensusPower(ctx, powers[1]), types.Unbonded, val2, true)
	require.NoError(t, err)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[0], app.StakingKeeper.TokensFromConsensusPower(ctx, powers[2]), types.Unbonded, val2, true)
	require.NoError(t, err)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, -1)

	return addrs, valAddrs, vals
}
