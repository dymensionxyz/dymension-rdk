package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	staking "github.com/dymensionxyz/dymension-rdk/x/governors"
	"github.com/dymensionxyz/dymension-rdk/x/governors/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/governors/teststaking"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

func newMonikerGovernor(t testing.TB, operator sdk.ValAddress, moniker string) types.Governor {
	v, err := types.NewGovernor(operator, types.Description{Moniker: moniker})
	require.NoError(t, err)
	return v
}

func bootstrapGovernorTest(t testing.TB, power int64, numAddrs int) (*app.App, sdk.Context, []sdk.AccAddress, []sdk.ValAddress) {
	_, app, ctx := createTestInput(&testing.T{})

	addrDels, addrVals := generateAddresses(app, ctx, numAddrs)

	amt := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
	totalSupply := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), amt.MulRaw(int64(len(addrDels)))))

	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	// set bonded pool supply
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), totalSupply))

	// unbond genesis governor delegations
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)
	delegation := delegations[0]

	_, err := app.StakingKeeper.Undelegate(ctx, delegation.GetDelegatorAddr(), delegation.GetGovernorAddr(), delegation.Shares)
	require.NoError(t, err)

	// end block to unbond genesis governor
	staking.EndBlocker(ctx, app.StakingKeeper)

	return app, ctx, addrDels, addrVals
}

func initGovernors(t testing.TB, power int64, numAddrs int, powers []int64) (*app.App, sdk.Context, []sdk.AccAddress, []sdk.ValAddress, []types.Governor) {
	app, ctx, addrs, valAddrs := bootstrapGovernorTest(t, power, numAddrs)

	vs := make([]types.Governor, len(powers))
	for i, power := range powers {
		vs[i] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[i]))
		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
		vs[i], _ = vs[i].AddTokensFromDel(tokens)
	}
	return app, ctx, addrs, valAddrs, vs
}

func TestSetGovernor(t *testing.T) {
	app, ctx, _, _ := bootstrapGovernorTest(t, 10, 100)

	valPubKey := PKs[0]
	valAddr := sdk.ValAddress(valPubKey.Address().Bytes())
	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)

	// test how the governor is set from a purely unbonbed pool
	governor := teststaking.NewGovernor(t, valAddr)
	governor, _ = governor.AddTokensFromDel(valTokens)
	require.Equal(t, types.Unbonded, governor.Status)
	assert.Equal(t, valTokens, governor.Tokens)
	assert.Equal(t, valTokens, governor.DelegatorShares.RoundInt())
	app.StakingKeeper.SetGovernor(ctx, governor)
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governor)

	// ensure update
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)
	governor, found := app.StakingKeeper.GetGovernor(ctx, valAddr)
	require.True(t, found)
	// require.Equal(t, governor.ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])

	// after the save the governor should be bonded
	require.Equal(t, types.Bonded, governor.Status)
	assert.Equal(t, valTokens, governor.Tokens)
	assert.Equal(t, valTokens, governor.DelegatorShares.RoundInt())

	// Check each store for being saved
	resVal, found := app.StakingKeeper.GetGovernor(ctx, valAddr)
	assert.True(ValEq(t, governor, resVal))
	require.True(t, found)

	resVals := app.StakingKeeper.GetLastGovernors(ctx)
	require.Equal(t, 1, len(resVals))
	assert.True(ValEq(t, governor, resVals[0]))

	resVals = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, 1, len(resVals))
	require.True(ValEq(t, governor, resVals[0]))

	resVals = app.StakingKeeper.GetGovernors(ctx, 1)
	require.Equal(t, 1, len(resVals))

	resVals = app.StakingKeeper.GetGovernors(ctx, 10)
	require.Equal(t, 2, len(resVals))

	allVals := app.StakingKeeper.GetAllGovernors(ctx)
	require.Equal(t, 2, len(allVals))
}

func TestUpdateGovernorByPowerIndex(t *testing.T) {
	app, ctx, _, _ := bootstrapGovernorTest(t, 0, 100)
	_, addrVals := generateAddresses(app, ctx, 1)

	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)

	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), app.StakingKeeper.TokensFromConsensusPower(ctx, 1234)))))
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), app.StakingKeeper.TokensFromConsensusPower(ctx, 10000)))))

	app.AccountKeeper.SetModuleAccount(ctx, bondedPool)
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	// add a governor
	governor := teststaking.NewGovernor(t, addrVals[0])
	governor, delSharesCreated := governor.AddTokensFromDel(app.StakingKeeper.TokensFromConsensusPower(ctx, 100))
	require.Equal(t, types.Unbonded, governor.Status)
	require.Equal(t, app.StakingKeeper.TokensFromConsensusPower(ctx, 100), governor.Tokens)
	keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	governor, found := app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, app.StakingKeeper.TokensFromConsensusPower(ctx, 100), governor.Tokens)

	power := types.GetGovernorsByPowerIndexKey(governor, app.StakingKeeper.PowerReduction(ctx))
	require.True(t, keeper.GovernorByPowerIndexExists(ctx, app.StakingKeeper, power))

	// burn half the delegator shares
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governor)
	governor, burned := governor.RemoveDelShares(delSharesCreated.Quo(sdk.NewDec(2)))
	require.Equal(t, app.StakingKeeper.TokensFromConsensusPower(ctx, 50), burned)
	keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true) // update the governor, possibly kicking it out
	require.False(t, keeper.GovernorByPowerIndexExists(ctx, app.StakingKeeper, power))

	governor, found = app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)

	power = types.GetGovernorsByPowerIndexKey(governor, app.StakingKeeper.PowerReduction(ctx))
	require.True(t, keeper.GovernorByPowerIndexExists(ctx, app.StakingKeeper, power))
}

func TestUpdateBondedGovernorsDecreaseCliff(t *testing.T) {
	numVals := 10
	maxVals := 5

	// create context, keeper, and pool for tests
	app, ctx, _, valAddrs := bootstrapGovernorTest(t, 0, 100)

	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)

	// create keeper parameters
	params := app.StakingKeeper.GetParams(ctx)
	params.MaxValidators = uint32(maxVals)
	app.StakingKeeper.SetParams(ctx, params)

	// create a random pool
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), app.StakingKeeper.TokensFromConsensusPower(ctx, 1234)))))
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), app.StakingKeeper.TokensFromConsensusPower(ctx, 10000)))))

	app.AccountKeeper.SetModuleAccount(ctx, bondedPool)
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	governors := make([]types.Governor, numVals)
	for i := 0; i < len(governors); i++ {
		moniker := fmt.Sprintf("val#%d", int64(i))
		val := newMonikerGovernor(t, valAddrs[i], moniker)
		delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, int64((i+1)*10))
		val, _ = val.AddTokensFromDel(delTokens)

		val = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, val, true)
		governors[i] = val
	}

	nextCliffVal := governors[numVals-maxVals+1]

	// remove enough tokens to kick out the governor below the current cliff
	// governor and next in line cliff governor
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, nextCliffVal)
	shares := app.StakingKeeper.TokensFromConsensusPower(ctx, 21)
	nextCliffVal, _ = nextCliffVal.RemoveDelShares(sdk.NewDecFromInt(shares))
	nextCliffVal = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, nextCliffVal, true)

	expectedValStatus := map[int]types.BondStatus{
		9: types.Bonded, 8: types.Bonded, 7: types.Bonded, 5: types.Bonded, 4: types.Bonded,
		0: types.Unbonding, 1: types.Unbonding, 2: types.Unbonding, 3: types.Unbonding, 6: types.Unbonding,
	}

	// require all the governors have their respective statuses
	for valIdx, status := range expectedValStatus {
		valAddr := governors[valIdx].OperatorAddress
		addr, err := sdk.ValAddressFromBech32(valAddr)
		assert.NoError(t, err)
		val, _ := app.StakingKeeper.GetGovernor(ctx, addr)

		assert.Equal(
			t, status, val.GetStatus(),
			fmt.Sprintf("expected governor at index %v to have status: %s", valIdx, status),
		)
	}
}

// This function tests UpdateGovernor, GetGovernor, GetLastGovernors, RemoveGovernor
func TestGovernorBasics(t *testing.T) {
	app, ctx, _, addrVals := bootstrapGovernorTest(t, 1000, 20)

	// construct the governors
	var governors [3]types.Governor
	powers := []int64{9, 8, 7}
	for i, power := range powers {
		governors[i] = teststaking.NewGovernor(t, addrVals[i])
		governors[i].Status = types.Unbonded
		governors[i].Tokens = sdk.ZeroInt()
		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)

		governors[i], _ = governors[i].AddTokensFromDel(tokens)
	}
	assert.Equal(t, app.StakingKeeper.TokensFromConsensusPower(ctx, 9), governors[0].Tokens)
	assert.Equal(t, app.StakingKeeper.TokensFromConsensusPower(ctx, 8), governors[1].Tokens)
	assert.Equal(t, app.StakingKeeper.TokensFromConsensusPower(ctx, 7), governors[2].Tokens)

	// check the empty keeper first
	_, found := app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.False(t, found)
	resVals := app.StakingKeeper.GetLastGovernors(ctx)
	require.Zero(t, len(resVals))

	resVals = app.StakingKeeper.GetGovernors(ctx, 2)
	require.Len(t, resVals, 1)

	// set and retrieve a record
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], true)
	resVal, found := app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)
	assert.True(ValEq(t, governors[0], resVal))

	resVals = app.StakingKeeper.GetLastGovernors(ctx)
	require.Equal(t, 1, len(resVals))
	assert.True(ValEq(t, governors[0], resVals[0]))
	assert.Equal(t, types.Bonded, governors[0].Status)
	assert.True(sdk.IntEq(t, app.StakingKeeper.TokensFromConsensusPower(ctx, 9), governors[0].BondedTokens()))

	// modify a records, save, and retrieve
	governors[0].Status = types.Bonded
	governors[0].Tokens = app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	governors[0].DelegatorShares = sdk.NewDecFromInt(governors[0].Tokens)
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], true)
	resVal, found = app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)
	assert.True(ValEq(t, governors[0], resVal))

	resVals = app.StakingKeeper.GetLastGovernors(ctx)
	require.Equal(t, 1, len(resVals))
	assert.True(ValEq(t, governors[0], resVals[0]))

	// add other governors
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], true)
	governors[2] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[2], true)
	resVal, found = app.StakingKeeper.GetGovernor(ctx, addrVals[1])
	require.True(t, found)
	assert.True(ValEq(t, governors[1], resVal))
	resVal, found = app.StakingKeeper.GetGovernor(ctx, addrVals[2])
	require.True(t, found)
	assert.True(ValEq(t, governors[2], resVal))

	resVals = app.StakingKeeper.GetLastGovernors(ctx)
	require.Equal(t, 3, len(resVals))
	assert.True(ValEq(t, governors[0], resVals[0])) // order doesn't matter here
	assert.True(ValEq(t, governors[1], resVals[1]))
	assert.True(ValEq(t, governors[2], resVals[2]))

	// remove a record

	// shouldn't be able to remove if status is not unbonded
	assert.PanicsWithValue(t,
		"cannot call RemoveGovernor on bonded or unbonding governors",
		func() { app.StakingKeeper.RemoveGovernor(ctx, governors[1].GetOperator()) })

	// shouldn't be able to remove if there are still tokens left
	governors[1].Status = types.Unbonded
	app.StakingKeeper.SetGovernor(ctx, governors[1])
	assert.PanicsWithValue(t,
		"attempting to remove a governor which still contains tokens",
		func() { app.StakingKeeper.RemoveGovernor(ctx, governors[1].GetOperator()) })

	governors[1].Tokens = sdk.ZeroInt()                               // ...remove all tokens
	app.StakingKeeper.SetGovernor(ctx, governors[1])                  // ...set the governor
	app.StakingKeeper.RemoveGovernor(ctx, governors[1].GetOperator()) // Now it can be removed.
	_, found = app.StakingKeeper.GetGovernor(ctx, addrVals[1])
	require.False(t, found)
}

// test how the governors are sorted, tests GetBondedGovernorsByPower
func TestGetGovernorSortingUnmixed(t *testing.T) {
	app, ctx, addrs, _ := bootstrapGovernorTest(t, 1000, 20)

	// initialize some governors into the state
	amts := []sdk.Int{
		sdk.NewIntFromUint64(0),
		app.StakingKeeper.PowerReduction(ctx).MulRaw(100),
		app.StakingKeeper.PowerReduction(ctx),
		app.StakingKeeper.PowerReduction(ctx).MulRaw(400),
		app.StakingKeeper.PowerReduction(ctx).MulRaw(200),
	}
	n := len(amts)
	var governors [5]types.Governor
	for i, amt := range amts {
		governors[i] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[i]))
		governors[i].Status = types.Bonded
		governors[i].Tokens = amt
		governors[i].DelegatorShares = sdk.NewDecFromInt(amt)
		keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[i], true)
	}

	// first make sure everything made it in to the gotGovernor group
	resGovernors := app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	assert.Equal(t, n, len(resGovernors))
	assert.Equal(t, sdk.NewInt(400).Mul(app.StakingKeeper.PowerReduction(ctx)), resGovernors[0].BondedTokens(), "%v", resGovernors)
	assert.Equal(t, sdk.NewInt(200).Mul(app.StakingKeeper.PowerReduction(ctx)), resGovernors[1].BondedTokens(), "%v", resGovernors)
	assert.Equal(t, sdk.NewInt(100).Mul(app.StakingKeeper.PowerReduction(ctx)), resGovernors[2].BondedTokens(), "%v", resGovernors)
	assert.Equal(t, sdk.NewInt(1).Mul(app.StakingKeeper.PowerReduction(ctx)), resGovernors[3].BondedTokens(), "%v", resGovernors)
	assert.Equal(t, sdk.NewInt(0), resGovernors[4].BondedTokens(), "%v", resGovernors)
	assert.Equal(t, governors[3].OperatorAddress, resGovernors[0].OperatorAddress, "%v", resGovernors)
	assert.Equal(t, governors[4].OperatorAddress, resGovernors[1].OperatorAddress, "%v", resGovernors)
	assert.Equal(t, governors[1].OperatorAddress, resGovernors[2].OperatorAddress, "%v", resGovernors)
	assert.Equal(t, governors[2].OperatorAddress, resGovernors[3].OperatorAddress, "%v", resGovernors)
	assert.Equal(t, governors[0].OperatorAddress, resGovernors[4].OperatorAddress, "%v", resGovernors)

	// test a basic increase in voting power
	governors[3].Tokens = sdk.NewInt(500).Mul(app.StakingKeeper.PowerReduction(ctx))
	keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[3], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, len(resGovernors), n)
	assert.True(ValEq(t, governors[3], resGovernors[0]))

	// test a decrease in voting power
	governors[3].Tokens = sdk.NewInt(300).Mul(app.StakingKeeper.PowerReduction(ctx))
	keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[3], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, len(resGovernors), n)
	assert.True(ValEq(t, governors[3], resGovernors[0]))
	assert.True(ValEq(t, governors[4], resGovernors[1]))

	// test equal voting power, different age
	governors[3].Tokens = sdk.NewInt(200).Mul(app.StakingKeeper.PowerReduction(ctx))
	ctx = ctx.WithBlockHeight(10)
	keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[3], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, len(resGovernors), n)
	assert.True(ValEq(t, governors[3], resGovernors[0]))
	assert.True(ValEq(t, governors[4], resGovernors[1]))

	// no change in voting power - no change in sort
	ctx = ctx.WithBlockHeight(20)
	keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[4], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, len(resGovernors), n)
	assert.True(ValEq(t, governors[3], resGovernors[0]))
	assert.True(ValEq(t, governors[4], resGovernors[1]))

	// change in voting power of both governors, both still in v-set, no age change
	governors[3].Tokens = sdk.NewInt(300).Mul(app.StakingKeeper.PowerReduction(ctx))
	governors[4].Tokens = sdk.NewInt(300).Mul(app.StakingKeeper.PowerReduction(ctx))
	keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[3], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, len(resGovernors), n)
	ctx = ctx.WithBlockHeight(30)
	keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[4], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, len(resGovernors), n, "%v", resGovernors)
	assert.True(ValEq(t, governors[3], resGovernors[0]))
	assert.True(ValEq(t, governors[4], resGovernors[1]))
}

func TestGetGovernorSortingMixed(t *testing.T) {
	app, ctx, addrs, _ := bootstrapGovernorTest(t, 1000, 20)
	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)

	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), app.StakingKeeper.TokensFromConsensusPower(ctx, 501)))))
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), app.StakingKeeper.TokensFromConsensusPower(ctx, 0)))))

	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
	app.AccountKeeper.SetModuleAccount(ctx, bondedPool)

	// now 2 max resGovernors
	params := app.StakingKeeper.GetParams(ctx)
	params.MaxValidators = 2
	app.StakingKeeper.SetParams(ctx, params)

	// initialize some governors into the state
	amts := []sdk.Int{
		sdk.NewIntFromUint64(0),
		app.StakingKeeper.PowerReduction(ctx).MulRaw(100),
		app.StakingKeeper.PowerReduction(ctx),
		app.StakingKeeper.PowerReduction(ctx).MulRaw(400),
		app.StakingKeeper.PowerReduction(ctx).MulRaw(200),
	}

	var governors [5]types.Governor
	for i, amt := range amts {
		governors[i] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[i]))
		governors[i].DelegatorShares = sdk.NewDecFromInt(amt)
		governors[i].Status = types.Bonded
		governors[i].Tokens = amt
		keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[i], true)
	}

	val0, found := app.StakingKeeper.GetGovernor(ctx, sdk.ValAddress(addrs[0]))
	require.True(t, found)
	val1, found := app.StakingKeeper.GetGovernor(ctx, sdk.ValAddress(addrs[1]))
	require.True(t, found)
	val2, found := app.StakingKeeper.GetGovernor(ctx, sdk.ValAddress(addrs[2]))
	require.True(t, found)
	val3, found := app.StakingKeeper.GetGovernor(ctx, sdk.ValAddress(addrs[3]))
	require.True(t, found)
	val4, found := app.StakingKeeper.GetGovernor(ctx, sdk.ValAddress(addrs[4]))
	require.True(t, found)
	require.Equal(t, types.Bonded, val0.Status)
	require.Equal(t, types.Unbonding, val1.Status)
	require.Equal(t, types.Unbonding, val2.Status)
	require.Equal(t, types.Bonded, val3.Status)
	require.Equal(t, types.Bonded, val4.Status)

	// first make sure everything made it in to the gotGovernor group
	resGovernors := app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	// The governors returned should match the max governors
	assert.Equal(t, 2, len(resGovernors))
	assert.Equal(t, sdk.NewInt(400).Mul(app.StakingKeeper.PowerReduction(ctx)), resGovernors[0].BondedTokens(), "%v", resGovernors)
	assert.Equal(t, sdk.NewInt(200).Mul(app.StakingKeeper.PowerReduction(ctx)), resGovernors[1].BondedTokens(), "%v", resGovernors)
	assert.Equal(t, governors[3].OperatorAddress, resGovernors[0].OperatorAddress, "%v", resGovernors)
	assert.Equal(t, governors[4].OperatorAddress, resGovernors[1].OperatorAddress, "%v", resGovernors)
}

// TODO separate out into multiple tests
func TestGetGovernorsEdgeCases(t *testing.T) {
	app, ctx, addrs, _ := bootstrapGovernorTest(t, 1000, 20)

	// set max governors to 2
	params := app.StakingKeeper.GetParams(ctx)
	nMax := uint32(2)
	params.MaxValidators = nMax
	app.StakingKeeper.SetParams(ctx, params)

	// initialize some governors into the state
	powers := []int64{0, 100, 400, 400}
	var governors [4]types.Governor
	for i, power := range powers {
		moniker := fmt.Sprintf("val#%d", int64(i))
		governors[i] = newMonikerGovernor(t, sdk.ValAddress(addrs[i]), moniker)

		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
		governors[i], _ = governors[i].AddTokensFromDel(tokens)

		notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
		require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(params.BondDenom, tokens))))
		app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)
		governors[i] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[i], true)
	}

	// ensure that the first two bonded governors are the largest governors
	resGovernors := app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, nMax, uint32(len(resGovernors)))
	assert.True(ValEq(t, governors[2], resGovernors[0]))
	assert.True(ValEq(t, governors[3], resGovernors[1]))

	// delegate 500 tokens to governor 0
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governors[0])
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 500)
	governors[0], _ = governors[0].AddTokensFromDel(delTokens)
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)

	newTokens := sdk.NewCoins()

	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), newTokens))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	// test that the two largest governors are
	//   a) governor 0 with 500 tokens
	//   b) governor 2 with 400 tokens (delegated before governor 3)
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, nMax, uint32(len(resGovernors)))
	assert.True(ValEq(t, governors[0], resGovernors[0]))
	assert.True(ValEq(t, governors[2], resGovernors[1]))

	// A governor which leaves the bonded governor set due to a decrease in voting power,
	// then increases to the original voting power, does not get its spot back in the
	// case of a tie.
	//
	// Order of operations for this test:
	//  - governor 3 enter governor set with 1 new token
	//  - governor 3 removed governor set by removing 201 tokens (governor 2 enters)
	//  - governor 3 adds 200 tokens (equal to governor 2 now) and does not get its spot back

	// governor 3 enters bonded governor set
	ctx = ctx.WithBlockHeight(40)

	var found bool
	governors[3], found = app.StakingKeeper.GetGovernor(ctx, governors[3].GetOperator())
	assert.True(t, found)
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governors[3])
	governors[3], _ = governors[3].AddTokensFromDel(app.StakingKeeper.TokensFromConsensusPower(ctx, 1))

	notBondedPool = app.StakingKeeper.GetNotBondedPool(ctx)
	newTokens = sdk.NewCoins(sdk.NewCoin(params.BondDenom, app.StakingKeeper.TokensFromConsensusPower(ctx, 1)))
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), newTokens))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	governors[3] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[3], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, nMax, uint32(len(resGovernors)))
	assert.True(ValEq(t, governors[0], resGovernors[0]))
	assert.True(ValEq(t, governors[3], resGovernors[1]))

	// governor 3 kicked out temporarily
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governors[3])
	rmTokens := governors[3].TokensFromShares(sdk.NewDec(201)).TruncateInt()
	governors[3], _ = governors[3].RemoveDelShares(sdk.NewDec(201))

	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(params.BondDenom, rmTokens))))
	app.AccountKeeper.SetModuleAccount(ctx, bondedPool)

	governors[3] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[3], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, nMax, uint32(len(resGovernors)))
	assert.True(ValEq(t, governors[0], resGovernors[0]))
	assert.True(ValEq(t, governors[2], resGovernors[1]))

	// governor 3 does not get spot back
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governors[3])
	governors[3], _ = governors[3].AddTokensFromDel(sdk.NewInt(200))

	notBondedPool = app.StakingKeeper.GetNotBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdk.NewInt(200)))))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	governors[3] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[3], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, nMax, uint32(len(resGovernors)))
	assert.True(ValEq(t, governors[0], resGovernors[0]))
	assert.True(ValEq(t, governors[2], resGovernors[1]))
	_, exists := app.StakingKeeper.GetGovernor(ctx, governors[3].GetOperator())
	require.True(t, exists)
}

func TestGovernorBondHeight(t *testing.T) {
	app, ctx, addrs, _ := bootstrapGovernorTest(t, 1000, 20)

	// now 2 max resGovernors
	params := app.StakingKeeper.GetParams(ctx)
	params.MaxValidators = 2
	app.StakingKeeper.SetParams(ctx, params)

	// initialize some governors into the state
	var governors [3]types.Governor
	governors[0] = teststaking.NewGovernor(t, sdk.ValAddress(PKs[0].Address().Bytes()))
	governors[1] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[1]))
	governors[2] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[2]))

	tokens0 := app.StakingKeeper.TokensFromConsensusPower(ctx, 200)
	tokens1 := app.StakingKeeper.TokensFromConsensusPower(ctx, 100)
	tokens2 := app.StakingKeeper.TokensFromConsensusPower(ctx, 100)
	governors[0], _ = governors[0].AddTokensFromDel(tokens0)
	governors[1], _ = governors[1].AddTokensFromDel(tokens1)
	governors[2], _ = governors[2].AddTokensFromDel(tokens2)

	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], true)

	////////////////////////////////////////
	// If two governors both increase to the same voting power in the same block,
	// the one with the first transaction should become bonded
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], true)
	governors[2] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[2], true)

	resGovernors := app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, uint32(len(resGovernors)), params.MaxValidators)

	assert.True(ValEq(t, governors[0], resGovernors[0]))
	assert.True(ValEq(t, governors[1], resGovernors[1]))
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governors[1])
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governors[2])
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 50)
	governors[1], _ = governors[1].AddTokensFromDel(delTokens)
	governors[2], _ = governors[2].AddTokensFromDel(delTokens)
	governors[2] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[2], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	require.Equal(t, params.MaxValidators, uint32(len(resGovernors)))
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], true)
	assert.True(ValEq(t, governors[0], resGovernors[0]))
	assert.True(ValEq(t, governors[2], resGovernors[1]))
}

func TestFullGovernorSetPowerChange(t *testing.T) {
	app, ctx, addrs, _ := bootstrapGovernorTest(t, 1000, 20)
	params := app.StakingKeeper.GetParams(ctx)
	max := 2
	params.MaxValidators = uint32(2)
	app.StakingKeeper.SetParams(ctx, params)

	// initialize some governors into the state
	powers := []int64{0, 100, 400, 400, 200}
	var governors [5]types.Governor
	for i, power := range powers {
		governors[i] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[i]))
		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
		governors[i], _ = governors[i].AddTokensFromDel(tokens)
		keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[i], true)
	}
	for i := range powers {
		var found bool
		governors[i], found = app.StakingKeeper.GetGovernor(ctx, governors[i].GetOperator())
		require.True(t, found)
	}
	assert.Equal(t, types.Unbonded, governors[0].Status)
	assert.Equal(t, types.Unbonding, governors[1].Status)
	assert.Equal(t, types.Bonded, governors[2].Status)
	assert.Equal(t, types.Bonded, governors[3].Status)
	assert.Equal(t, types.Unbonded, governors[4].Status)
	resGovernors := app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	assert.Equal(t, max, len(resGovernors))
	assert.True(ValEq(t, governors[2], resGovernors[0])) // in the order of txs
	assert.True(ValEq(t, governors[3], resGovernors[1]))

	// test a swap in voting power

	tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 600)
	governors[0], _ = governors[0].AddTokensFromDel(tokens)
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], true)
	resGovernors = app.StakingKeeper.GetBondedGovernorsByPower(ctx)
	assert.Equal(t, max, len(resGovernors))
	assert.True(ValEq(t, governors[0], resGovernors[0]))
	assert.True(ValEq(t, governors[2], resGovernors[1]))
}

func TestApplyAndReturnGovernorSetUpdatesAllNone(t *testing.T) {
	app, ctx, _, _ := bootstrapGovernorTest(t, 1000, 20)

	powers := []int64{10, 20}
	var governors [2]types.Governor
	for i, power := range powers {
		valPubKey := PKs[i+1]
		valAddr := sdk.ValAddress(valPubKey.Address().Bytes())

		governors[i] = teststaking.NewGovernor(t, valAddr)
		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
		governors[i], _ = governors[i].AddTokensFromDel(tokens)
	}

	// test from nothing to something
	//  tendermintUpdate set: {} -> {c1, c3}
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 0)
	app.StakingKeeper.SetGovernor(ctx, governors[0])
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[0])
	app.StakingKeeper.SetGovernor(ctx, governors[1])
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[1])

	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)
	governors[0], _ = app.StakingKeeper.GetGovernor(ctx, governors[0].GetOperator())
	governors[1], _ = app.StakingKeeper.GetGovernor(ctx, governors[1].GetOperator())
	// assert.Equal(t, governors[0].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[1])
	// assert.Equal(t, governors[1].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])
}

func TestApplyAndReturnGovernorSetUpdatesIdentical(t *testing.T) {
	app, ctx, addrs, _ := bootstrapGovernorTest(t, 1000, 20)

	powers := []int64{10, 20}
	var governors [2]types.Governor
	for i, power := range powers {
		governors[i] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[i]))

		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
		governors[i], _ = governors[i].AddTokensFromDel(tokens)

	}
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], false)
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], false)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)

	// test identical,
	//  tendermintUpdate set: {} -> {}
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], false)
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], false)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 0)
}

func TestApplyAndReturnGovernorSetUpdatesSingleValueChange(t *testing.T) {
	app, ctx, addrs, _ := bootstrapGovernorTest(t, 1000, 20)

	powers := []int64{10, 20}
	var governors [2]types.Governor
	for i, power := range powers {
		governors[i] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[i]))

		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
		governors[i], _ = governors[i].AddTokensFromDel(tokens)

	}
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], false)
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], false)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)

	// test single value change
	//  tendermintUpdate set: {} -> {c1'}
	governors[0].Status = types.Bonded
	governors[0].Tokens = app.StakingKeeper.TokensFromConsensusPower(ctx, 600)
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], false)

	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)
	// require.Equal(t, governors[0].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])
}

func TestApplyAndReturnGovernorSetUpdatesMultipleValueChange(t *testing.T) {
	powers := []int64{10, 20}
	// TODO: use it in other places
	app, ctx, _, _, governors := initGovernors(t, 1000, 20, powers)

	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], false)
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], false)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)

	// test multiple value change
	//  tendermintUpdate set: {c1, c3} -> {c1', c3'}
	delTokens1 := app.StakingKeeper.TokensFromConsensusPower(ctx, 190)
	delTokens2 := app.StakingKeeper.TokensFromConsensusPower(ctx, 80)
	governors[0], _ = governors[0].AddTokensFromDel(delTokens1)
	governors[1], _ = governors[1].AddTokensFromDel(delTokens2)
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], false)
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], false)

	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)
	// require.Equal(t, governors[0].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])
	// require.Equal(t, governors[1].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[1])
}

func TestApplyAndReturnGovernorSetUpdatesInserted(t *testing.T) {
	powers := []int64{10, 20, 5, 15, 25}
	app, ctx, _, _, governors := initGovernors(t, 1000, 20, powers)

	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], false)
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], false)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)

	// test validtor added at the beginning
	//  tendermintUpdate set: {} -> {c0}
	app.StakingKeeper.SetGovernor(ctx, governors[2])
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[2])
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)
	governors[2], _ = app.StakingKeeper.GetGovernor(ctx, governors[2].GetOperator())
	// require.Equal(t, governors[2].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])

	// test validtor added at the beginning
	//  tendermintUpdate set: {} -> {c0}
	app.StakingKeeper.SetGovernor(ctx, governors[3])
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[3])
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)
	governors[3], _ = app.StakingKeeper.GetGovernor(ctx, governors[3].GetOperator())
	// require.Equal(t, governors[3].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])

	// test validtor added at the end
	//  tendermintUpdate set: {} -> {c0}
	app.StakingKeeper.SetGovernor(ctx, governors[4])
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[4])
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)
	governors[4], _ = app.StakingKeeper.GetGovernor(ctx, governors[4].GetOperator())
	// require.Equal(t, governors[4].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])
}

func TestApplyAndReturnGovernorSetUpdatesWithCliffGovernor(t *testing.T) {
	app, ctx, addrs, _ := bootstrapGovernorTest(t, 1000, 20)
	params := types.DefaultParams()
	params.MaxValidators = 2
	app.StakingKeeper.SetParams(ctx, params)

	powers := []int64{10, 20, 5}
	var governors [5]types.Governor
	for i, power := range powers {
		governors[i] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[i]))
		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
		governors[i], _ = governors[i].AddTokensFromDel(tokens)
	}
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], false)
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], false)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)

	// test governor added at the end but not inserted in the valset
	//  tendermintUpdate set: {} -> {}
	keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[2], false)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 0)

	// test governor change its power and become a gotGovernor (pushing out an existing)
	//  tendermintUpdate set: {}     -> {c0, c4}
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 0)

	tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	governors[2], _ = governors[2].AddTokensFromDel(tokens)
	app.StakingKeeper.SetGovernor(ctx, governors[2])
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[2])
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)
	governors[2], _ = app.StakingKeeper.GetGovernor(ctx, governors[2].GetOperator())
	// require.Equal(t, governors[0].ABCIGovernorUpdateZero(), updates[1])
	// require.Equal(t, governors[2].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])
}

func TestApplyAndReturnGovernorSetUpdatesPowerDecrease(t *testing.T) {
	app, ctx, addrs, _ := bootstrapGovernorTest(t, 1000, 20)

	powers := []int64{100, 100}
	var governors [2]types.Governor
	for i, power := range powers {
		governors[i] = teststaking.NewGovernor(t, sdk.ValAddress(addrs[i]))
		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
		governors[i], _ = governors[i].AddTokensFromDel(tokens)
	}
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], false)
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], false)
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)

	// check initial power
	require.Equal(t, int64(100), governors[0].GetConsensusPower(app.StakingKeeper.PowerReduction(ctx)))
	require.Equal(t, int64(100), governors[1].GetConsensusPower(app.StakingKeeper.PowerReduction(ctx)))

	// test multiple value change
	//  tendermintUpdate set: {c1, c3} -> {c1', c3'}
	delTokens1 := app.StakingKeeper.TokensFromConsensusPower(ctx, 20)
	delTokens2 := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
	governors[0], _ = governors[0].RemoveDelShares(sdk.NewDecFromInt(delTokens1))
	governors[1], _ = governors[1].RemoveDelShares(sdk.NewDecFromInt(delTokens2))
	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], false)
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], false)

	// power has changed
	require.Equal(t, int64(80), governors[0].GetConsensusPower(app.StakingKeeper.PowerReduction(ctx)))
	require.Equal(t, int64(70), governors[1].GetConsensusPower(app.StakingKeeper.PowerReduction(ctx)))

	// Tendermint updates should reflect power change
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)
	// require.Equal(t, governors[0].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])
	// require.Equal(t, governors[1].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[1])
}

func TestApplyAndReturnGovernorSetUpdatesNewGovernor(t *testing.T) {
	app, ctx, _, _ := bootstrapGovernorTest(t, 1000, 20)
	params := app.StakingKeeper.GetParams(ctx)
	params.MaxValidators = uint32(3)

	app.StakingKeeper.SetParams(ctx, params)

	powers := []int64{100, 100}
	var governors [2]types.Governor

	// initialize some governors into the state
	for i, power := range powers {
		valPubKey := PKs[i+1]
		valAddr := sdk.ValAddress(valPubKey.Address().Bytes())

		governors[i] = teststaking.NewGovernor(t, valAddr)
		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
		governors[i], _ = governors[i].AddTokensFromDel(tokens)

		app.StakingKeeper.SetGovernor(ctx, governors[i])
		app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[i])
	}

	// verify initial Tendermint updates are correct
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, len(governors))
	governors[0], _ = app.StakingKeeper.GetGovernor(ctx, governors[0].GetOperator())
	governors[1], _ = app.StakingKeeper.GetGovernor(ctx, governors[1].GetOperator())
	// require.Equal(t, governors[0].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])
	// require.Equal(t, governors[1].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[1])

	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 0)

	// update initial governor set
	for i, power := range powers {

		app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governors[i])
		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
		governors[i], _ = governors[i].AddTokensFromDel(tokens)

		app.StakingKeeper.SetGovernor(ctx, governors[i])
		app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[i])
	}

	// add a new governor that goes from zero power, to non-zero power, back to
	// zero power
	valPubKey := PKs[len(governors)+1]
	valAddr := sdk.ValAddress(valPubKey.Address().Bytes())
	amt := sdk.NewInt(100)

	governor := teststaking.NewGovernor(t, valAddr)
	governor, _ = governor.AddTokensFromDel(amt)

	app.StakingKeeper.SetGovernor(ctx, governor)

	governor, _ = governor.RemoveDelShares(sdk.NewDecFromInt(amt))
	app.StakingKeeper.SetGovernor(ctx, governor)
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governor)

	// add a new governor that increases in power
	valAddr = sdk.ValAddress(valPubKey.Address().Bytes())

	governor = teststaking.NewGovernor(t, valAddr)
	tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 500)
	governor, _ = governor.AddTokensFromDel(tokens)
	app.StakingKeeper.SetGovernor(ctx, governor)
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governor)

	// verify initial Tendermint updates are correct
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, len(governors)+1)
	governor, _ = app.StakingKeeper.GetGovernor(ctx, governor.GetOperator())
	// governors[0], _ = app.StakingKeeper.GetGovernor(ctx, governors[0].GetOperator())
	governors[1], _ = app.StakingKeeper.GetGovernor(ctx, governors[1].GetOperator())
	// require.Equal(t, governor.ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])
	// require.Equal(t, governors[0].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[1])
	// require.Equal(t, governors[1].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[2])
}

func TestApplyAndReturnGovernorSetUpdatesBondTransition(t *testing.T) {
	app, ctx, _, _ := bootstrapGovernorTest(t, 1000, 20)
	params := app.StakingKeeper.GetParams(ctx)
	params.MaxValidators = uint32(2)

	app.StakingKeeper.SetParams(ctx, params)

	powers := []int64{100, 200, 300}
	var governors [3]types.Governor

	// initialize some governors into the state
	for i, power := range powers {
		moniker := fmt.Sprintf("%d", i)
		valPubKey := PKs[i+1]
		valAddr := sdk.ValAddress(valPubKey.Address().Bytes())

		governors[i] = newMonikerGovernor(t, valAddr, moniker)
		tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, power)
		governors[i], _ = governors[i].AddTokensFromDel(tokens)
		app.StakingKeeper.SetGovernor(ctx, governors[i])
		app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[i])
	}

	// verify initial Tendermint updates are correct
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)
	governors[2], _ = app.StakingKeeper.GetGovernor(ctx, governors[2].GetOperator())
	governors[1], _ = app.StakingKeeper.GetGovernor(ctx, governors[1].GetOperator())
	// require.Equal(t, governors[2].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])
	// require.Equal(t, governors[1].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[1])

	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 0)

	// delegate to governor with lowest power but not enough to bond
	ctx = ctx.WithBlockHeight(1)

	var found bool
	governors[0], found = app.StakingKeeper.GetGovernor(ctx, governors[0].GetOperator())
	require.True(t, found)

	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governors[0])
	tokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 1)
	governors[0], _ = governors[0].AddTokensFromDel(tokens)
	app.StakingKeeper.SetGovernor(ctx, governors[0])
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[0])

	// verify initial Tendermint updates are correct
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 0)

	// create a series of events that will bond and unbond the governor with
	// lowest power in a single block context (height)
	ctx = ctx.WithBlockHeight(2)

	governors[1], found = app.StakingKeeper.GetGovernor(ctx, governors[1].GetOperator())
	require.True(t, found)

	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governors[0])
	governors[0], _ = governors[0].RemoveDelShares(governors[0].DelegatorShares)
	app.StakingKeeper.SetGovernor(ctx, governors[0])
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[0])
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 0)

	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governors[1])
	tokens = app.StakingKeeper.TokensFromConsensusPower(ctx, 250)
	governors[1], _ = governors[1].AddTokensFromDel(tokens)
	app.StakingKeeper.SetGovernor(ctx, governors[1])
	app.StakingKeeper.SetGovernorByPowerIndex(ctx, governors[1])

	// verify initial Tendermint updates are correct
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)
	// require.Equal(t, governors[1].ABCIGovernorUpdate(app.StakingKeeper.PowerReduction(ctx)), updates[0])
	//
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 0)
}

func TestUpdateGovernorCommission(t *testing.T) {
	app, ctx, _, addrVals := bootstrapGovernorTest(t, 1000, 20)
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: time.Now().UTC()})

	// Set MinCommissionRate to 0.05
	params := app.StakingKeeper.GetParams(ctx)
	params.MinCommissionRate = sdk.NewDecWithPrec(5, 2)
	app.StakingKeeper.SetParams(ctx, params)

	commission1 := types.NewCommissionWithTime(
		sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(3, 1),
		sdk.NewDecWithPrec(1, 1), time.Now().UTC().Add(time.Duration(-1)*time.Hour),
	)
	commission2 := types.NewCommission(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(3, 1), sdk.NewDecWithPrec(1, 1))

	val1 := teststaking.NewGovernor(t, addrVals[0])
	val2 := teststaking.NewGovernor(t, addrVals[1])

	val1, _ = val1.SetInitialCommission(commission1)
	val2, _ = val2.SetInitialCommission(commission2)

	app.StakingKeeper.SetGovernor(ctx, val1)
	app.StakingKeeper.SetGovernor(ctx, val2)

	testCases := []struct {
		governor    types.Governor
		newRate     sdk.Dec
		expectedErr bool
	}{
		{val1, sdk.ZeroDec(), true},
		{val2, sdk.NewDecWithPrec(-1, 1), true},
		{val2, sdk.NewDecWithPrec(4, 1), true},
		{val2, sdk.NewDecWithPrec(3, 1), true},
		{val2, sdk.NewDecWithPrec(1, 2), true},
		{val2, sdk.NewDecWithPrec(2, 1), false},
	}

	for i, tc := range testCases {
		commission, err := app.StakingKeeper.UpdateGovernorCommission(ctx, tc.governor, tc.newRate)

		if tc.expectedErr {
			require.Error(t, err, "expected error for test case #%d with rate: %s", i, tc.newRate)
		} else {
			tc.governor.Commission = commission
			app.StakingKeeper.SetGovernor(ctx, tc.governor)
			val, found := app.StakingKeeper.GetGovernor(ctx, tc.governor.GetOperator())

			require.True(t, found,
				"expected to find governor for test case #%d with rate: %s", i, tc.newRate,
			)
			require.NoError(t, err,
				"unexpected error for test case #%d with rate: %s", i, tc.newRate,
			)
			require.Equal(t, tc.newRate, val.Commission.Rate,
				"expected new governor commission rate for test case #%d with rate: %s", i, tc.newRate,
			)
			require.Equal(t, ctx.BlockHeader().Time, val.Commission.UpdateTime,
				"expected new governor commission update time for test case #%d with rate: %s", i, tc.newRate,
			)
		}
	}
}

func applyGovernorsSetUpdates(t *testing.T, ctx sdk.Context, k keeper.Keeper, expectedUpdatesLen int) {
	err := k.ApplyGovernorSetUpdates(ctx)
	require.NoError(t, err)
}
