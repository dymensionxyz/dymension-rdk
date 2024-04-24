package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/governors/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/governors/teststaking"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

// tests GetDelegation, GetDelegatorDelegations, SetDelegation, RemoveDelegation, GetDelegatorDelegations
func TestDelegation(t *testing.T) {
	_, app, ctx := createTestInput(t)

	// remove genesis governor delegations
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	require.Len(t, delegations, 1)

	app.StakingKeeper.RemoveDelegation(ctx, stakingtypes.Delegation{
		ValidatorAddress: delegations[0].ValidatorAddress,
		DelegatorAddress: delegations[0].DelegatorAddress,
	})

	addrDels := utils.AddTestAddrs(app, ctx, 3, sdk.NewInt(10000))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrDels)

	// construct the governors
	amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
	var governors [3]types.Governor
	for i, amt := range amts {
		governors[i] = teststaking.NewGovernor(t, valAddrs[i])
		governors[i], _ = governors[i].AddTokensFromDel(amt)
	}

	governors[0] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[0], true)
	governors[1] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[1], true)
	governors[2] = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governors[2], true)

	// first add a governors[0] to delegate too
	bond1to1 := stakingtypes.NewDelegation(addrDels[0], valAddrs[0], sdk.NewDec(9))

	// check the empty keeper first
	_, found := app.StakingKeeper.GetDelegation(ctx, addrDels[0], valAddrs[0])
	require.False(t, found)

	// set and retrieve a record
	app.StakingKeeper.SetDelegation(ctx, bond1to1)
	resBond, found := app.StakingKeeper.GetDelegation(ctx, addrDels[0], valAddrs[0])
	require.True(t, found)
	require.Equal(t, bond1to1, resBond)

	// modify a records, save, and retrieve
	bond1to1.Shares = sdk.NewDec(99)
	app.StakingKeeper.SetDelegation(ctx, bond1to1)
	resBond, found = app.StakingKeeper.GetDelegation(ctx, addrDels[0], valAddrs[0])
	require.True(t, found)
	require.Equal(t, bond1to1, resBond)

	// add some more records
	bond1to2 := stakingtypes.NewDelegation(addrDels[0], valAddrs[1], sdk.NewDec(9))
	bond1to3 := stakingtypes.NewDelegation(addrDels[0], valAddrs[2], sdk.NewDec(9))
	bond2to1 := stakingtypes.NewDelegation(addrDels[1], valAddrs[0], sdk.NewDec(9))
	bond2to2 := stakingtypes.NewDelegation(addrDels[1], valAddrs[1], sdk.NewDec(9))
	bond2to3 := stakingtypes.NewDelegation(addrDels[1], valAddrs[2], sdk.NewDec(9))
	app.StakingKeeper.SetDelegation(ctx, bond1to2)
	app.StakingKeeper.SetDelegation(ctx, bond1to3)
	app.StakingKeeper.SetDelegation(ctx, bond2to1)
	app.StakingKeeper.SetDelegation(ctx, bond2to2)
	app.StakingKeeper.SetDelegation(ctx, bond2to3)

	// test all bond retrieve capabilities
	resBonds := app.StakingKeeper.GetDelegatorDelegations(ctx, addrDels[0], 5)
	require.Equal(t, 3, len(resBonds))
	require.Equal(t, bond1to1, resBonds[0])
	require.Equal(t, bond1to2, resBonds[1])
	require.Equal(t, bond1to3, resBonds[2])
	resBonds = app.StakingKeeper.GetAllDelegatorDelegations(ctx, addrDels[0])
	require.Equal(t, 3, len(resBonds))
	resBonds = app.StakingKeeper.GetDelegatorDelegations(ctx, addrDels[0], 2)
	require.Equal(t, 2, len(resBonds))
	resBonds = app.StakingKeeper.GetDelegatorDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 3, len(resBonds))
	require.Equal(t, bond2to1, resBonds[0])
	require.Equal(t, bond2to2, resBonds[1])
	require.Equal(t, bond2to3, resBonds[2])
	allBonds := app.StakingKeeper.GetAllDelegations(ctx)
	require.Equal(t, 6, len(allBonds))
	require.Equal(t, bond1to1, allBonds[0])
	require.Equal(t, bond1to2, allBonds[1])
	require.Equal(t, bond1to3, allBonds[2])
	require.Equal(t, bond2to1, allBonds[3])
	require.Equal(t, bond2to2, allBonds[4])
	require.Equal(t, bond2to3, allBonds[5])

	resVals := app.StakingKeeper.GetDelegatorGovernors(ctx, addrDels[0], 3)
	require.Equal(t, 3, len(resVals))
	resVals = app.StakingKeeper.GetDelegatorGovernors(ctx, addrDels[1], 4)
	require.Equal(t, 3, len(resVals))

	for i := 0; i < 3; i++ {
		resVal, err := app.StakingKeeper.GetDelegatorGovernor(ctx, addrDels[0], valAddrs[i])
		require.Nil(t, err)
		require.Equal(t, valAddrs[i], resVal.GetOperator())

		resVal, err = app.StakingKeeper.GetDelegatorGovernor(ctx, addrDels[1], valAddrs[i])
		require.Nil(t, err)
		require.Equal(t, valAddrs[i], resVal.GetOperator())

		resDels := app.StakingKeeper.GetGovernorDelegations(ctx, valAddrs[i])
		require.Len(t, resDels, 2)
	}

	// test total bonded for single delegator
	expBonded := bond1to1.Shares.Add(bond2to1.Shares).Add(bond1to3.Shares)
	resDelBond := app.StakingKeeper.GetDelegatorBonded(ctx, addrDels[0])
	require.Equal(t, expBonded, sdk.NewDecFromInt(resDelBond))

	// delete a record
	app.StakingKeeper.RemoveDelegation(ctx, bond2to3)
	_, found = app.StakingKeeper.GetDelegation(ctx, addrDels[1], valAddrs[2])
	require.False(t, found)
	resBonds = app.StakingKeeper.GetDelegatorDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 2, len(resBonds))
	require.Equal(t, bond2to1, resBonds[0])
	require.Equal(t, bond2to2, resBonds[1])

	resBonds = app.StakingKeeper.GetAllDelegatorDelegations(ctx, addrDels[1])
	require.Equal(t, 2, len(resBonds))

	// delete all the records from delegator 2
	app.StakingKeeper.RemoveDelegation(ctx, bond2to1)
	app.StakingKeeper.RemoveDelegation(ctx, bond2to2)
	_, found = app.StakingKeeper.GetDelegation(ctx, addrDels[1], valAddrs[0])
	require.False(t, found)
	_, found = app.StakingKeeper.GetDelegation(ctx, addrDels[1], valAddrs[1])
	require.False(t, found)
	resBonds = app.StakingKeeper.GetDelegatorDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 0, len(resBonds))
}

// tests Get/Set/Remove UnbondingDelegation
func TestUnbondingDelegation(t *testing.T) {
	_, app, ctx := createTestInput(t)

	delAddrs := utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))
	valAddrs := simapp.ConvertAddrsToValAddrs(delAddrs)

	ubd := stakingtypes.NewUnbondingDelegation(
		delAddrs[0],
		valAddrs[0],
		0,
		time.Unix(0, 0).UTC(),
		sdk.NewInt(5),
	)

	// set and retrieve a record
	app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
	resUnbond, found := app.StakingKeeper.GetUnbondingDelegation(ctx, delAddrs[0], valAddrs[0])
	require.True(t, found)
	require.Equal(t, ubd, resUnbond)

	// modify a records, save, and retrieve
	expUnbond := sdk.NewInt(21)
	ubd.Entries[0].Balance = expUnbond
	app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)

	resUnbonds := app.StakingKeeper.GetUnbondingDelegations(ctx, delAddrs[0], 5)
	require.Equal(t, 1, len(resUnbonds))

	resUnbonds = app.StakingKeeper.GetAllUnbondingDelegations(ctx, delAddrs[0])
	require.Equal(t, 1, len(resUnbonds))

	resUnbond, found = app.StakingKeeper.GetUnbondingDelegation(ctx, delAddrs[0], valAddrs[0])
	require.True(t, found)
	require.Equal(t, ubd, resUnbond)

	resDelUnbond := app.StakingKeeper.GetDelegatorUnbonding(ctx, delAddrs[0])
	require.Equal(t, expUnbond, resDelUnbond)

	// delete a record
	app.StakingKeeper.RemoveUnbondingDelegation(ctx, ubd)
	_, found = app.StakingKeeper.GetUnbondingDelegation(ctx, delAddrs[0], valAddrs[0])
	require.False(t, found)

	resUnbonds = app.StakingKeeper.GetUnbondingDelegations(ctx, delAddrs[0], 5)
	require.Equal(t, 0, len(resUnbonds))

	resUnbonds = app.StakingKeeper.GetAllUnbondingDelegations(ctx, delAddrs[0])
	require.Equal(t, 0, len(resUnbonds))
}

func TestUnbondDelegation(t *testing.T) {
	_, app, ctx := createTestInput(t)

	delAddrs := utils.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))
	valAddrs := simapp.ConvertAddrsToValAddrs(delAddrs)

	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)

	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a governor and a delegator to that governor
	// note this governor starts not-bonded
	governor := teststaking.NewGovernor(t, valAddrs[0])

	governor, issuedShares := governor.AddTokensFromDel(startTokens)
	require.Equal(t, startTokens, issuedShares.RoundInt())

	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)

	delegation := stakingtypes.NewDelegation(delAddrs[0], valAddrs[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)

	bondTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 6)
	amount, err := app.StakingKeeper.Unbond(ctx, delAddrs[0], valAddrs[0], sdk.NewDecFromInt(bondTokens))
	require.NoError(t, err)
	require.Equal(t, bondTokens, amount) // shares to be added to an unbonding delegation

	delegation, found := app.StakingKeeper.GetDelegation(ctx, delAddrs[0], valAddrs[0])
	require.True(t, found)
	governor, found = app.StakingKeeper.GetGovernor(ctx, valAddrs[0])
	require.True(t, found)

	remainingTokens := startTokens.Sub(bondTokens)
	require.Equal(t, remainingTokens, delegation.Shares.RoundInt())
	require.Equal(t, remainingTokens, governor.BondedTokens())
}

func TestUnbondingDelegationsMaxEntries(t *testing.T) {
	_, app, ctx := createTestInput(t)

	addrDels := utils.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)

	bondDenom := app.StakingKeeper.BondDenom(ctx)
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)

	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(bondDenom, startTokens))))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a governor and a delegator to that governor
	governor := teststaking.NewGovernor(t, addrVals[0])

	governor, issuedShares := governor.AddTokensFromDel(startTokens)
	require.Equal(t, startTokens, issuedShares.RoundInt())

	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	require.True(sdk.IntEq(t, startTokens, governor.BondedTokens()))
	require.True(t, governor.IsBonded())

	delegation := stakingtypes.NewDelegation(addrDels[0], addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)

	maxEntries := app.StakingKeeper.MaxEntries(ctx)

	oldBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
	oldNotBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount

	// should all pass
	var completionTime time.Time
	for i := uint32(0); i < maxEntries; i++ {
		var err error
		completionTime, err = app.StakingKeeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewDec(1))
		require.NoError(t, err)
	}

	newBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
	newNotBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
	require.True(sdk.IntEq(t, newBonded, oldBonded.SubRaw(int64(maxEntries))))
	require.True(sdk.IntEq(t, newNotBonded, oldNotBonded.AddRaw(int64(maxEntries))))

	oldBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
	oldNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount

	// an additional unbond should fail due to max entries
	_, err := app.StakingKeeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewDec(1))
	require.Error(t, err)

	newBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
	newNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount

	require.True(sdk.IntEq(t, newBonded, oldBonded))
	require.True(sdk.IntEq(t, newNotBonded, oldNotBonded))

	// mature unbonding delegations
	ctx = ctx.WithBlockTime(completionTime)
	_, err = app.StakingKeeper.CompleteUnbonding(ctx, addrDels[0], addrVals[0])
	require.NoError(t, err)

	newBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
	newNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
	require.True(sdk.IntEq(t, newBonded, oldBonded))
	require.True(sdk.IntEq(t, newNotBonded, oldNotBonded.SubRaw(int64(maxEntries))))

	oldNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount

	// unbonding  should work again
	_, err = app.StakingKeeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewDec(1))
	require.NoError(t, err)

	newBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
	newNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
	require.True(sdk.IntEq(t, newBonded, oldBonded.SubRaw(1)))
	require.True(sdk.IntEq(t, newNotBonded, oldNotBonded.AddRaw(1)))
}

// // test undelegating self delegation from a governor pushing it below MinSelfDelegation
// // shift it from the bonded to unbonding state and jailed
func TestUndelegateSelfDelegationBelowMinSelfDelegation(t *testing.T) {
	_, app, ctx := createTestInput(t)

	addrDels := utils.AddTestAddrs(app, ctx, 1, sdk.NewInt(10000))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	delCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), delTokens))

	// create a governor with a self-delegation
	governor := teststaking.NewGovernor(t, addrVals[0])

	governor.MinSelfDelegation = delTokens
	governor, issuedShares := governor.AddTokensFromDel(delTokens)
	require.Equal(t, delTokens, issuedShares.RoundInt())

	// add bonded tokens to pool for delegations
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), delCoins))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	require.True(t, governor.IsBonded())

	selfDelegation := stakingtypes.NewDelegation(sdk.AccAddress(addrVals[0].Bytes()), addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, selfDelegation)

	// add bonded tokens to pool for delegations
	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
	app.AccountKeeper.SetModuleAccount(ctx, bondedPool)

	// create a second delegation to this governor
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governor)
	governor, issuedShares = governor.AddTokensFromDel(delTokens)
	require.True(t, governor.IsBonded())
	require.Equal(t, delTokens, issuedShares.RoundInt())

	// add bonded tokens to pool for delegations
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
	app.AccountKeeper.SetModuleAccount(ctx, bondedPool)

	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	delegation := stakingtypes.NewDelegation(addrDels[0], addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)

	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
	_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(app.StakingKeeper.TokensFromConsensusPower(ctx, 6)))
	require.NoError(t, err)

	// end block
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)

	governor, found := app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, app.StakingKeeper.TokensFromConsensusPower(ctx, 14), governor.Tokens)
	require.Equal(t, types.Unbonding, governor.Status)
}

func TestUndelegateFromUnbondingGovernor(t *testing.T) {
	_, app, ctx := createTestInput(t)
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	delCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), delTokens))

	addrDels := utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(0))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	// create a governor with a self-delegation
	governor := teststaking.NewGovernor(t, addrVals[0])

	governor, issuedShares := governor.AddTokensFromDel(delTokens)
	require.Equal(t, delTokens, issuedShares.RoundInt())

	// add bonded tokens to pool for delegations
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), delCoins))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	require.True(t, governor.IsBonded())

	selfDelegation := stakingtypes.NewDelegation(addrVals[0].Bytes(), addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, selfDelegation)

	// add bonded tokens to pool for delegations
	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
	app.AccountKeeper.SetModuleAccount(ctx, bondedPool)

	// create a second delegation to this governor
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governor)

	governor, issuedShares = governor.AddTokensFromDel(delTokens)
	require.Equal(t, delTokens, issuedShares.RoundInt())

	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
	app.AccountKeeper.SetModuleAccount(ctx, bondedPool)

	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	delegation := stakingtypes.NewDelegation(addrDels[1], addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)

	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
	app.AccountKeeper.SetModuleAccount(ctx, bondedPool)

	header := ctx.BlockHeader()
	blockHeight := int64(10)
	header.Height = blockHeight
	blockTime := time.Unix(333, 0)
	header.Time = blockTime
	ctx = ctx.WithBlockHeader(header)

	// unbond the all self-delegation to put governor in unbonding state
	val0AccAddr := sdk.AccAddress(addrVals[0])
	_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(delTokens))
	require.NoError(t, err)

	// end block
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)

	governor, found := app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, blockHeight, governor.UnbondingHeight)
	params := app.StakingKeeper.GetParams(ctx)
	require.True(t, blockTime.Add(params.UnbondingTime).Equal(governor.UnbondingTime))

	blockHeight2 := int64(20)
	blockTime2 := time.Unix(444, 0).UTC()
	ctx = ctx.WithBlockHeight(blockHeight2)
	ctx = ctx.WithBlockTime(blockTime2)

	// unbond some of the other delegation's shares
	_, err = app.StakingKeeper.Undelegate(ctx, addrDels[1], addrVals[0], sdk.NewDec(6))
	require.NoError(t, err)

	// retrieve the unbonding delegation
	ubd, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrDels[1], addrVals[0])
	require.True(t, found)
	require.Len(t, ubd.Entries, 1)
	require.True(t, ubd.Entries[0].Balance.Equal(sdk.NewInt(6)))
	assert.Equal(t, blockHeight2, ubd.Entries[0].CreationHeight)
	assert.True(t, blockTime2.Add(params.UnbondingTime).Equal(ubd.Entries[0].CompletionTime))
}

func TestUndelegateFromUnbondedGovernor(t *testing.T) {
	_, app, ctx := createTestInput(t)
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	delCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), delTokens))

	addrDels := utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(0))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	// add bonded tokens to pool for delegations
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), delCoins))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a governor with a self-delegation
	governor := teststaking.NewGovernor(t, addrVals[0])

	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	governor, issuedShares := governor.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())
	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	require.True(t, governor.IsBonded())

	val0AccAddr := sdk.AccAddress(addrVals[0])
	selfDelegation := stakingtypes.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, selfDelegation)

	// add bonded tokens to pool for delegations
	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
	app.AccountKeeper.SetModuleAccount(ctx, bondedPool)

	// create a second delegation to this governor
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governor)
	governor, issuedShares = governor.AddTokensFromDel(delTokens)
	require.Equal(t, delTokens, issuedShares.RoundInt())
	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	require.True(t, governor.IsBonded())
	delegation := stakingtypes.NewDelegation(addrDels[1], addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)

	ctx = ctx.WithBlockHeight(10)
	ctx = ctx.WithBlockTime(time.Unix(333, 0))

	// unbond the all self-delegation to put governor in unbonding state
	_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(valTokens))
	require.NoError(t, err)

	// end block
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)

	governor, found := app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, ctx.BlockHeight(), governor.UnbondingHeight)
	params := app.StakingKeeper.GetParams(ctx)
	require.True(t, ctx.BlockHeader().Time.Add(params.UnbondingTime).Equal(governor.UnbondingTime))

	// unbond the governor
	ctx = ctx.WithBlockTime(governor.UnbondingTime)
	app.StakingKeeper.UnbondAllMatureGovernors(ctx)

	// Make sure governor is still in state because there is still an outstanding delegation
	governor, found = app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, governor.Status, types.Unbonded)

	// unbond some of the other delegation's shares
	unbondTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 6)
	_, err = app.StakingKeeper.Undelegate(ctx, addrDels[1], addrVals[0], sdk.NewDecFromInt(unbondTokens))
	require.NoError(t, err)

	// unbond rest of the other delegation's shares
	remainingTokens := delTokens.Sub(unbondTokens)
	_, err = app.StakingKeeper.Undelegate(ctx, addrDels[1], addrVals[0], sdk.NewDecFromInt(remainingTokens))
	require.NoError(t, err)

	//  now governor should be deleted from state
	governor, found = app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.False(t, found, "%v", governor)
}

func TestUnbondingAllDelegationFromGovernor(t *testing.T) {
	_, app, ctx := createTestInput(t)
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	delCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), delTokens))

	addrDels := utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(0))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	// add bonded tokens to pool for delegations
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), delCoins))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a governor with a self-delegation
	governor := teststaking.NewGovernor(t, addrVals[0])

	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	governor, issuedShares := governor.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())

	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	require.True(t, governor.IsBonded())
	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())

	selfDelegation := stakingtypes.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, selfDelegation)

	// create a second delegation to this governor
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governor)
	governor, issuedShares = governor.AddTokensFromDel(delTokens)
	require.Equal(t, delTokens, issuedShares.RoundInt())

	// add bonded tokens to pool for delegations
	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, bondedPool.GetName(), delCoins))
	app.AccountKeeper.SetModuleAccount(ctx, bondedPool)

	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	require.True(t, governor.IsBonded())

	delegation := stakingtypes.NewDelegation(addrDels[1], addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)

	ctx = ctx.WithBlockHeight(10)
	ctx = ctx.WithBlockTime(time.Unix(333, 0))

	// unbond the all self-delegation to put governor in unbonding state
	_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(valTokens))
	require.NoError(t, err)

	// end block
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)

	// unbond all the remaining delegation
	_, err = app.StakingKeeper.Undelegate(ctx, addrDels[1], addrVals[0], sdk.NewDecFromInt(delTokens))
	require.NoError(t, err)

	// governor should still be in state and still be in unbonding state
	governor, found := app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, governor.Status, types.Unbonding)

	// unbond the governor
	ctx = ctx.WithBlockTime(governor.UnbondingTime)
	app.StakingKeeper.UnbondAllMatureGovernors(ctx)

	// governor should now be deleted from state
	_, found = app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.False(t, found)
}

// Make sure that that the retrieving the delegations doesn't affect the state
func TestGetRedelegationsFromSrcGovernor(t *testing.T) {
	_, app, ctx := createTestInput(t)

	addrDels := utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(0))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	rd := stakingtypes.NewRedelegation(addrDels[0], addrVals[0], addrVals[1], 0,
		time.Unix(0, 0), sdk.NewInt(5),
		sdk.NewDec(5))

	// set and retrieve a record
	app.StakingKeeper.SetRedelegation(ctx, rd)
	resBond, found := app.StakingKeeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
	require.True(t, found)

	// get the redelegations one time
	redelegations := app.StakingKeeper.GetRedelegationsFromSrcGovernor(ctx, addrVals[0])
	require.Equal(t, 1, len(redelegations))
	require.Equal(t, redelegations[0], resBond)

	// get the redelegations a second time, should be exactly the same
	redelegations = app.StakingKeeper.GetRedelegationsFromSrcGovernor(ctx, addrVals[0])
	require.Equal(t, 1, len(redelegations))
	require.Equal(t, redelegations[0], resBond)
}

// tests Get/Set/Remove/Has UnbondingDelegation
func TestRedelegation(t *testing.T) {
	_, app, ctx := createTestInput(t)

	addrDels := utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(0))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	rd := stakingtypes.NewRedelegation(addrDels[0], addrVals[0], addrVals[1], 0,
		time.Unix(0, 0).UTC(), sdk.NewInt(5),
		sdk.NewDec(5))

	// test shouldn't have and redelegations
	has := app.StakingKeeper.HasReceivingRedelegation(ctx, addrDels[0], addrVals[1])
	require.False(t, has)

	// set and retrieve a record
	app.StakingKeeper.SetRedelegation(ctx, rd)
	resRed, found := app.StakingKeeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
	require.True(t, found)

	redelegations := app.StakingKeeper.GetRedelegationsFromSrcGovernor(ctx, addrVals[0])
	require.Equal(t, 1, len(redelegations))
	require.Equal(t, redelegations[0], resRed)

	redelegations = app.StakingKeeper.GetRedelegations(ctx, addrDels[0], 5)
	require.Equal(t, 1, len(redelegations))
	require.Equal(t, redelegations[0], resRed)

	redelegations = app.StakingKeeper.GetAllRedelegations(ctx, addrDels[0], nil, nil)
	require.Equal(t, 1, len(redelegations))
	require.Equal(t, redelegations[0], resRed)

	// check if has the redelegation
	has = app.StakingKeeper.HasReceivingRedelegation(ctx, addrDels[0], addrVals[1])
	require.True(t, has)

	// modify a records, save, and retrieve
	rd.Entries[0].SharesDst = sdk.NewDec(21)
	app.StakingKeeper.SetRedelegation(ctx, rd)

	resRed, found = app.StakingKeeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
	require.True(t, found)
	require.Equal(t, rd, resRed)

	redelegations = app.StakingKeeper.GetRedelegationsFromSrcGovernor(ctx, addrVals[0])
	require.Equal(t, 1, len(redelegations))
	require.Equal(t, redelegations[0], resRed)

	redelegations = app.StakingKeeper.GetRedelegations(ctx, addrDels[0], 5)
	require.Equal(t, 1, len(redelegations))
	require.Equal(t, redelegations[0], resRed)

	// delete a record
	app.StakingKeeper.RemoveRedelegation(ctx, rd)
	_, found = app.StakingKeeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
	require.False(t, found)

	redelegations = app.StakingKeeper.GetRedelegations(ctx, addrDels[0], 5)
	require.Equal(t, 0, len(redelegations))

	redelegations = app.StakingKeeper.GetAllRedelegations(ctx, addrDels[0], nil, nil)
	require.Equal(t, 0, len(redelegations))
}

func TestRedelegateToSameGovernor(t *testing.T) {
	_, app, ctx := createTestInput(t)

	addrDels := utils.AddTestAddrs(app, ctx, 1, sdk.NewInt(0))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	startCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), valTokens))

	// add bonded tokens to pool for delegations
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), startCoins))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a governor with a self-delegation
	governor := teststaking.NewGovernor(t, addrVals[0])
	governor, issuedShares := governor.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())
	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	require.True(t, governor.IsBonded())

	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
	selfDelegation := stakingtypes.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, selfDelegation)

	_, err := app.StakingKeeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[0], sdk.NewDec(5))
	require.Error(t, err)
}

func TestRedelegationMaxEntries(t *testing.T) {
	_, app, ctx := createTestInput(t)

	addrDels := utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(0))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 20)
	startCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))

	// add bonded tokens to pool for delegations
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), startCoins))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a governor with a self-delegation
	governor := teststaking.NewGovernor(t, addrVals[0])
	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	governor, issuedShares := governor.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())
	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
	selfDelegation := stakingtypes.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, selfDelegation)

	// create a second governor
	governor2 := teststaking.NewGovernor(t, addrVals[1])
	governor2, issuedShares = governor2.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())

	governor2 = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor2, true)
	require.Equal(t, types.Bonded, governor2.Status)

	maxEntries := app.StakingKeeper.MaxEntries(ctx)

	// redelegations should pass
	var completionTime time.Time
	for i := uint32(0); i < maxEntries; i++ {
		var err error
		completionTime, err = app.StakingKeeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], sdk.NewDec(1))
		require.NoError(t, err)
	}

	// an additional redelegation should fail due to max entries
	_, err := app.StakingKeeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], sdk.NewDec(1))
	require.Error(t, err)

	// mature redelegations
	ctx = ctx.WithBlockTime(completionTime)
	_, err = app.StakingKeeper.CompleteRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1])
	require.NoError(t, err)

	// redelegation should work again
	_, err = app.StakingKeeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], sdk.NewDec(1))
	require.NoError(t, err)
}

func TestRedelegateSelfDelegation(t *testing.T) {
	_, app, ctx := createTestInput(t)

	addrDels := utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(0))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
	startCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))

	// add bonded tokens to pool for delegations
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), startCoins))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a governor with a self-delegation
	governor := teststaking.NewGovernor(t, addrVals[0])

	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	governor, issuedShares := governor.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())

	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)

	val0AccAddr := sdk.AccAddress(addrVals[0])
	selfDelegation := stakingtypes.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, selfDelegation)

	// create a second governor
	governor2 := teststaking.NewGovernor(t, addrVals[1])
	governor2, issuedShares = governor2.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())
	governor2 = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor2, true)
	require.Equal(t, types.Bonded, governor2.Status)

	// create a second delegation to governor 1
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	governor, issuedShares = governor.AddTokensFromDel(delTokens)
	require.Equal(t, delTokens, issuedShares.RoundInt())
	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)

	delegation := stakingtypes.NewDelegation(addrDels[0], addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)

	_, err := app.StakingKeeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], sdk.NewDecFromInt(delTokens))
	require.NoError(t, err)

	// end block
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 2)

	governor, found := app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, valTokens, governor.Tokens)
	require.Equal(t, types.Unbonding, governor.Status)
}

func TestRedelegateFromUnbondingGovernor(t *testing.T) {
	_, app, ctx := createTestInput(t)

	addrDels := utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(0))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
	startCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))

	// add bonded tokens to pool for delegations
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), startCoins))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a governor with a self-delegation
	governor := teststaking.NewGovernor(t, addrVals[0])

	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	governor, issuedShares := governor.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())
	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
	selfDelegation := stakingtypes.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, selfDelegation)

	// create a second delegation to this governor
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governor)
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	governor, issuedShares = governor.AddTokensFromDel(delTokens)
	require.Equal(t, delTokens, issuedShares.RoundInt())
	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	delegation := stakingtypes.NewDelegation(addrDels[1], addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)

	// create a second governor
	governor2 := teststaking.NewGovernor(t, addrVals[1])
	governor2, issuedShares = governor2.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())
	governor2 = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor2, true)

	header := ctx.BlockHeader()
	blockHeight := int64(10)
	header.Height = blockHeight
	blockTime := time.Unix(333, 0)
	header.Time = blockTime
	ctx = ctx.WithBlockHeader(header)

	// unbond the all self-delegation to put governor in unbonding state
	_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(delTokens))
	require.NoError(t, err)

	// end block
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)

	governor, found := app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, blockHeight, governor.UnbondingHeight)
	params := app.StakingKeeper.GetParams(ctx)
	require.True(t, blockTime.Add(params.UnbondingTime).Equal(governor.UnbondingTime))

	// change the context
	header = ctx.BlockHeader()
	blockHeight2 := int64(20)
	header.Height = blockHeight2
	blockTime2 := time.Unix(444, 0)
	header.Time = blockTime2
	ctx = ctx.WithBlockHeader(header)

	// unbond some of the other delegation's shares
	redelegateTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 6)
	_, err = app.StakingKeeper.BeginRedelegation(ctx, addrDels[1], addrVals[0], addrVals[1], sdk.NewDecFromInt(redelegateTokens))
	require.NoError(t, err)

	// retrieve the unbonding delegation
	ubd, found := app.StakingKeeper.GetRedelegation(ctx, addrDels[1], addrVals[0], addrVals[1])
	require.True(t, found)
	require.Len(t, ubd.Entries, 1)
	assert.Equal(t, blockHeight, ubd.Entries[0].CreationHeight)
	assert.True(t, blockTime.Add(params.UnbondingTime).Equal(ubd.Entries[0].CompletionTime))
}

func TestRedelegateFromUnbondedGovernor(t *testing.T) {
	_, app, ctx := createTestInput(t)

	addrDels := utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(0))
	addrVals := simapp.ConvertAddrsToValAddrs(addrDels)

	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 30)
	startCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))

	// add bonded tokens to pool for delegations
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), startCoins))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a governor with a self-delegation
	governor := teststaking.NewGovernor(t, addrVals[0])

	valTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	governor, issuedShares := governor.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())
	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
	selfDelegation := stakingtypes.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, selfDelegation)

	// create a second delegation to this governor
	app.StakingKeeper.DeleteGovernorByPowerIndex(ctx, governor)
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 10)
	governor, issuedShares = governor.AddTokensFromDel(delTokens)
	require.Equal(t, delTokens, issuedShares.RoundInt())
	governor = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor, true)
	delegation := stakingtypes.NewDelegation(addrDels[1], addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)

	// create a second governor
	governor2 := teststaking.NewGovernor(t, addrVals[1])
	governor2, issuedShares = governor2.AddTokensFromDel(valTokens)
	require.Equal(t, valTokens, issuedShares.RoundInt())
	governor2 = keeper.TestingUpdateGovernor(app.StakingKeeper, ctx, governor2, true)
	require.Equal(t, types.Bonded, governor2.Status)

	ctx = ctx.WithBlockHeight(10)
	ctx = ctx.WithBlockTime(time.Unix(333, 0))

	// unbond the all self-delegation to put governor in unbonding state
	_, err := app.StakingKeeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.NewDecFromInt(delTokens))
	require.NoError(t, err)

	// end block
	applyGovernorsSetUpdates(t, ctx, app.StakingKeeper, 1)

	governor, found := app.StakingKeeper.GetGovernor(ctx, addrVals[0])
	require.True(t, found)
	require.Equal(t, ctx.BlockHeight(), governor.UnbondingHeight)
	params := app.StakingKeeper.GetParams(ctx)
	require.True(t, ctx.BlockHeader().Time.Add(params.UnbondingTime).Equal(governor.UnbondingTime))

	// unbond the governor
	app.StakingKeeper.UnbondingToUnbonded(ctx, governor)

	// redelegate some of the delegation's shares
	redelegationTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 6)
	_, err = app.StakingKeeper.BeginRedelegation(ctx, addrDels[1], addrVals[0], addrVals[1], sdk.NewDecFromInt(redelegationTokens))
	require.NoError(t, err)

	// no red should have been found
	red, found := app.StakingKeeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
	require.False(t, found, "%v", red)
}