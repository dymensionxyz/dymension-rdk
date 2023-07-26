package keeper_test

import (
	"fmt"
	"testing"

	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/rollapp/app"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	seqkeeper "github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	seqtypes "github.com/dymensionxyz/dymension-rdk/x/sequencers/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	PKS = simapp.CreateTestPubKeys(5)

	valConsPk1 = PKS[0]
	valConsPk2 = PKS[1]
	valConsPk3 = PKS[2]

	valConsAddr2 = sdk.ConsAddress(valConsPk2.Address())

	totalFees     = sdk.NewInt(100)
	totalFeesCoin = sdk.NewCoin(sdk.DefaultBondDenom, totalFees)
	totalFeesDec  = sdk.NewDecFromInt(totalFees)
)

//Test multiple sequencers, each propose a block

/* -------------------------------------------------------------------------- */
/*                                    utils                                   */
/* -------------------------------------------------------------------------- */
func assertInitial(t *testing.T, ctx sdk.Context, app *app.App, valAddrs []sdk.ValAddress) {
	// assert initial state: zero outstanding rewards, zero community pool, zero commission, zero current rewards
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[0]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[1]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetFeePool(ctx).CommunityPool.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[0]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[1]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[0]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[1]).Rewards.IsZero())
}

func fundModules(t *testing.T, ctx sdk.Context, app *app.App) {
	fees := sdk.NewCoins(totalFeesCoin)
	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	require.NotNil(t, feeCollector)

	// fund fee collector
	utils.FundModuleAccount(app, ctx, feeCollector.GetName(), fees)
	// require.NoError(t, simapp.FundModuleAccount(app.BankKeeper, ctx, feeCollector.GetName(), fees))
	app.AccountKeeper.SetAccount(ctx, feeCollector)
}

func createSeq(t *testing.T, ctx sdk.Context, app *app.App, valAddr sdk.ValAddress) {
	// create sequencer for dymint
	err := app.SequencersKeeper.SetDymintSequencerByAddr(ctx, sdk.GetConsAddress(valConsPk2), 0)
	require.NoError(t, err)

	// create sequencer
	msgServ := seqkeeper.NewMsgServerImpl(app.SequencersKeeper)
	description := stakingtypes.NewDescription(
		"moniker",
		"identity",
		"website",
		"security",
		"details",
	)

	msg, _ := seqtypes.NewMsgCreateSequencer(
		sdk.ValAddress(valAddr), valConsPk2, description,
	)
	_, err = msgServ.CreateSequencer(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
}

func createValidators(t *testing.T, ctx sdk.Context, app *app.App) []sdk.ValAddress {
	addrs := utils.AddTestAddrs(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// create validator with 6 power and 50% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	coin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(6, sdk.DefaultPowerReduction))
	msg, err := stakingtypes.NewMsgCreateValidator(valAddrs[0], valConsPk1, coin, stakingtypes.Description{}, tstaking.Commission, sdk.OneInt())
	require.NoError(t, err)
	_, err = stakingkeeper.NewMsgServerImpl(app.StakingKeeper).CreateValidator(ctx, msg)
	require.NoError(t, err)

	// create second validator with 4 power and 10% commision
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(1, 1), sdk.NewDec(0))
	coin = sdk.NewCoin(sdk.DefaultBondDenom, sdk.TokensFromConsensusPower(4, sdk.DefaultPowerReduction))
	msg, err = stakingtypes.NewMsgCreateValidator(valAddrs[1], valConsPk2, coin, stakingtypes.Description{}, tstaking.Commission, sdk.OneInt())
	require.NoError(t, err)
	_, err = stakingkeeper.NewMsgServerImpl(app.StakingKeeper).CreateValidator(ctx, msg)
	require.NoError(t, err)
	return valAddrs
}

/* -------------------------------------------------------------------------- */
/*                          stakers only, no proposer                         */
/* -------------------------------------------------------------------------- */
func TestAllocateTokensValidatorsNoProposer(t *testing.T) {
	app := utils.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	//TODO: test with different params
	proposerReward := 0.4
	communityTax := 0.02
	app.DistrKeeper.SetParams(ctx, disttypes.Params{
		CommunityTax:        sdk.MustNewDecFromStr(fmt.Sprintf("%f", communityTax)),
		BaseProposerReward:  sdk.MustNewDecFromStr(fmt.Sprintf("%f", proposerReward)),
		BonusProposerReward: sdk.MustNewDecFromStr("0"),
		WithdrawAddrEnabled: false,
	})

	valAddrs := createValidators(t, ctx, app)
	assertInitial(t, ctx, app, valAddrs)
	fundModules(t, ctx, app)

	// end block to bond validator and start new block
	_ = app.StakingKeeper.BlockValidatorUpdates(ctx)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// allocate tokens as if both had voted and second was proposer
	app.DistrKeeper.AllocateTokens(ctx, valConsAddr2)

	/* ------------------------------ Test stakers ------------------------------ */
	// outstanding rewards: 60% to val1 and 40% to val2
	stakersFees := totalFeesDec.MulTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", (1 - proposerReward - communityTax))))
	val1Coins := stakersFees.MulTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", 0.6)))
	val2Coins := stakersFees.MulTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", 0.4)))

	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val1Coins}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[0]).Rewards)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val2Coins}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[1]).Rewards)

	// Check commissions and delegator rewards val1
	//val1 has 50% commission
	val1Commission := val1Coins.Mul(sdk.MustNewDecFromStr(fmt.Sprintf("%f", 0.5)))
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val1Commission}}, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[0]).Commission)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val1Coins.Sub(val1Commission)}}, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[0]).Rewards)

	// Check commissions and delegator rewards val2
	//val2 has 10% commission
	val2Commission := val2Coins.Mul(sdk.MustNewDecFromStr(fmt.Sprintf("%f", 0.1)))
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val2Commission}}, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[1]).Commission)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val2Coins.Sub(val2Commission)}}, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[1]).Rewards)

	/* ------------------------ Test community pool coins ----------------------- */
	minCommunityFund := totalFeesDec.MulTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", communityTax)))
	communityBalance := app.DistrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(sdk.DefaultBondDenom)
	require.True(t, communityBalance.GTE(minCommunityFund))

	leftoverFees := totalFeesDec.Sub(val1Coins).Sub(val2Coins)
	if leftoverFees.IsPositive() {
		require.Equal(t, leftoverFees, communityBalance)
	}

	if app.DistrKeeper.GetFeePool(ctx).CommunityPool.IsZero() {
		require.True(t, communityTax == 0)
	}
}

/* -------------------------------------------------------------------------- */
/*                          proposer only, no stakers                         */
/* -------------------------------------------------------------------------- */
func TestAllocateTokensToProposerNoValidators(t *testing.T) {
	app := utils.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs := utils.AddTestAddrs(app, ctx, 2, sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)

	fundModules(t, ctx, app)

	// create sequencer
	createSeq(t, ctx, app, valAddrs[1])

	proposerReward := 0.4
	communityTax := 0.02
	app.DistrKeeper.SetParams(ctx, disttypes.Params{
		CommunityTax:        sdk.MustNewDecFromStr(fmt.Sprintf("%f", communityTax)),
		BaseProposerReward:  sdk.MustNewDecFromStr(fmt.Sprintf("%f", proposerReward)),
		BonusProposerReward: sdk.MustNewDecFromStr("0"),
		WithdrawAddrEnabled: false,
	})
	// allocate tokens as if both had voted and second was proposer
	app.DistrKeeper.AllocateTokens(ctx, valConsAddr2)

	/* ------------------------- Test proposer rewards ------------------------ */
	proposerFees := totalFeesDec.MulTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", proposerReward)))

	initialBalance := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	currentBalance := app.BankKeeper.GetAllBalances(ctx, sdk.AccAddress(valAddrs[1]))
	//expected = initial + proposer fees
	expectedBalance := initialBalance.Add(proposerFees.RoundInt())
	expectedCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, expectedBalance))

	require.Equal(t, expectedCoins, currentBalance)

	/* ------------------------ Test community pool coins ----------------------- */
	minCommunityFund := totalFeesDec.MulTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", communityTax)))
	communityBalance := app.DistrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(sdk.DefaultBondDenom)
	require.True(t, communityBalance.GTE(minCommunityFund))

	leftoverFees := totalFeesDec.Sub(proposerFees)
	if leftoverFees.IsPositive() {
		require.Equal(t, leftoverFees, communityBalance)
	}

	if app.DistrKeeper.GetFeePool(ctx).CommunityPool.IsZero() {
		require.True(t, communityTax == 0)
	}
}

/* -------------------------------------------------------------------------- */
/*                          both proposer and agents                          */
/* -------------------------------------------------------------------------- */
func TestAllocateTokensValidatorsAndProposer(t *testing.T) {
	app := utils.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	valAddrs := createValidators(t, ctx, app)
	assertInitial(t, ctx, app, valAddrs)
	fundModules(t, ctx, app)

	// end block to bond validator and start new block
	_ = app.StakingKeeper.BlockValidatorUpdates(ctx)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// create sequencer
	createSeq(t, ctx, app, valAddrs[1])

	proposerReward := 0.4
	communityTax := 0.02
	app.DistrKeeper.SetParams(ctx, disttypes.Params{
		CommunityTax:        sdk.MustNewDecFromStr(fmt.Sprintf("%f", communityTax)),
		BaseProposerReward:  sdk.MustNewDecFromStr(fmt.Sprintf("%f", proposerReward)),
		BonusProposerReward: sdk.MustNewDecFromStr("0"),
		WithdrawAddrEnabled: false,
	})
	// allocate tokens as if both had voted and second was proposer
	app.DistrKeeper.AllocateTokens(ctx, valConsAddr2)

	/* ------------------------- Test proposer rewards ------------------------ */
	proposerFees := totalFeesDec.MulTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", proposerReward)))

	initialBalance := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	currentBalance := app.BankKeeper.GetAllBalances(ctx, sdk.AccAddress(valAddrs[1]))
	//expected = initial + proposer fees - staked amount
	expectedBalance := initialBalance.Add(proposerFees.RoundInt()).Sub(sdk.TokensFromConsensusPower(4, sdk.DefaultPowerReduction))
	expectedCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, expectedBalance))

	require.Equal(t, expectedCoins, currentBalance)

	/* ------------------------------ Test stakers ------------------------------ */
	// outstanding rewards: 60% to val1 and 40% to val2
	//val1 has 50% commission as well
	stakersFees := totalFeesDec.MulTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", (1 - proposerReward - communityTax))))
	val1Coins := stakersFees.MulTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", 0.6)))
	val2Coins := stakersFees.MulTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", 0.4)))

	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val1Coins}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[0]).Rewards)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val2Coins}}, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[1]).Rewards)

	// Check commissions and delegator rewards val1
	val1Commission := val1Coins.Mul(sdk.MustNewDecFromStr(fmt.Sprintf("%f", 0.5)))
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val1Commission}}, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[0]).Commission)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val1Coins.Sub(val1Commission)}}, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[0]).Rewards)

	// Check commissions and delegator rewards val2
	val2Commission := val2Coins.Mul(sdk.MustNewDecFromStr(fmt.Sprintf("%f", 0.1)))
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val2Commission}}, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[1]).Commission)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: val2Coins.Sub(val2Commission)}}, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[1]).Rewards)

	/* ------------------------ Test community pool coins ----------------------- */
	minCommunityFund := totalFeesDec.MulTruncate(sdk.MustNewDecFromStr(fmt.Sprintf("%f", communityTax)))
	communityBalance := app.DistrKeeper.GetFeePool(ctx).CommunityPool.AmountOf(sdk.DefaultBondDenom)
	require.True(t, communityBalance.GTE(minCommunityFund))

	leftoverFees := totalFeesDec.Sub(proposerFees).Sub(val1Coins).Sub(val2Coins)
	if leftoverFees.IsPositive() {
		require.Equal(t, leftoverFees, communityBalance)
	}

	if app.DistrKeeper.GetFeePool(ctx).CommunityPool.IsZero() {
		require.True(t, communityTax == 0)
	}
}

/* -------------------------------------------------------------------------- */
/*                               original tests                               */
/* -------------------------------------------------------------------------- */
func TestAllocateTokensTruncation(t *testing.T) {
	app := utils.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs := utils.AddTestAddrs(app, ctx, 3, sdk.NewInt(1234))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// create validator with 10% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(1, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[0], valConsPk1, sdk.NewInt(110), true)

	// create second validator with 10% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(1, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[1], valConsPk2, sdk.NewInt(100), true)

	// create third validator with 10% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(1, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[2], valConsPk3, sdk.NewInt(100), true)

	abciValA := abci.Validator{
		Address: valConsPk1.Address(),
		Power:   11,
	}
	abciValB := abci.Validator{
		Address: valConsPk2.Address(),
		Power:   10,
	}
	abciValС := abci.Validator{
		Address: valConsPk3.Address(),
		Power:   10,
	}

	// assert initial state: zero outstanding rewards, zero community pool, zero commission, zero current rewards
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[0]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[1]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[1]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetFeePool(ctx).CommunityPool.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[0]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, valAddrs[1]).Commission.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[0]).Rewards.IsZero())
	require.True(t, app.DistrKeeper.GetValidatorCurrentRewards(ctx, valAddrs[1]).Rewards.IsZero())

	// allocate tokens as if both had voted and second was proposer
	fees := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(634195840)))

	feeCollector := app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
	require.NotNil(t, feeCollector)

	utils.FundModuleAccount(app, ctx, feeCollector.GetName(), fees)

	app.AccountKeeper.SetAccount(ctx, feeCollector)

	_ = []abci.VoteInfo{
		{
			Validator:       abciValA,
			SignedLastBlock: true,
		},
		{
			Validator:       abciValB,
			SignedLastBlock: true,
		},
		{
			Validator:       abciValС,
			SignedLastBlock: true,
		},
	}
	app.DistrKeeper.AllocateTokens(ctx, sdk.ConsAddress(valConsPk2.Address()))

	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[0]).Rewards.IsValid())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[1]).Rewards.IsValid())
	require.True(t, app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[2]).Rewards.IsValid())
}

func TestAllocateTokensToValidatorWithCommission(t *testing.T) {
	app := utils.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs := utils.AddTestAddrs(app, ctx, 3, sdk.NewInt(1234))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrs)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	// create validator with 50% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(sdk.ValAddress(addrs[0]), valConsPk1, sdk.NewInt(100), true)
	val := app.StakingKeeper.Validator(ctx, valAddrs[0])

	// allocate tokens
	tokens := sdk.DecCoins{
		{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(10)},
	}
	app.DistrKeeper.AllocateTokensToValidator(ctx, val, tokens)

	// check commission
	expected := sdk.DecCoins{
		{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(5)},
	}
	require.Equal(t, expected, app.DistrKeeper.GetValidatorAccumulatedCommission(ctx, val.GetOperator()).Commission)

	// check current rewards
	require.Equal(t, expected, app.DistrKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator()).Rewards)
}
