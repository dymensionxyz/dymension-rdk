package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/dist/keeper"
	"github.com/evmos/evmos/v12/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/ethereum/go-ethereum/common"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func TestHooks(t *testing.T) {
	app := utils.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{
		ChainID: "test_100-1",
	}).WithChainID("test_100-1")
	ctx = ctx.WithConsensusParams(utils.DefaultConsensusParams)

	app.SequencersKeeper.SetSequencer(ctx, utils.Proposer)
	app.SequencersKeeper.SetRewardAddr(ctx, utils.Proposer, utils.OperatorAcc())
	ctx = ctx.WithProposer(utils.ProposerCons())

	// Create native denom "foo" and register it as ERC20
	fooDenom := "foo"
	fooMetadata := banktypes.Metadata{
		Base:        fooDenom,
		Name:        "Foo",
		Symbol:      "FOO",
		Description: "fdsfds",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "FOO",
				Exponent: 0,
			},
			{
				Denom:    fooDenom,
				Exponent: 18,
			},
		},
	}
	app.BankKeeper.SetDenomMetaData(ctx, fooMetadata)
	pair, err := app.Erc20Keeper.RegisterCoin(ctx, fooMetadata)
	require.NoError(t, err)

	// create validators
	valAddrs := createValidators(t, ctx, app)
	app.StakingKeeper.BlockValidatorUpdates(ctx)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	assertInitial(t, ctx, app, valAddrs)

	// fund the fee collector - both foo (erc20) and non-erc20
	fees := sdk.NewCoins(sdk.NewCoin(fooDenom, sdk.NewInt(1000000000000000000)), sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000000000000000)))
	fundModules(t, ctx, app, fees)

	// allocate tokens - all fees goes to stakers
	app.DistrKeeper.SetParams(ctx, types.Params{
		CommunityTax:        sdk.ZeroDec(),
		BaseProposerReward:  sdk.ZeroDec(),
		BonusProposerReward: sdk.ZeroDec(),
		WithdrawAddrEnabled: false,
	})
	app.DistrKeeper.AllocateTokens(ctx, utils.ProposerCons())

	// trigger the hook, make sure the balance is converted
	hooks := keeper.Hooks{
		Hooks:      app.DistrKeeper.Keeper.Hooks(),
		DistKeeper: app.DistrKeeper,
	}
	delegation, found := app.StakingKeeper.GetDelegation(ctx, sdk.AccAddress(valAddrs[0]), valAddrs[0])
	require.True(t, found)
	err = hooks.BeforeDelegationSharesModified(ctx, delegation.GetDelegatorAddr(), delegation.GetValidatorAddr())
	require.NoError(t, err)

	// get balance
	balance := app.BankKeeper.GetBalance(ctx, delegation.GetDelegatorAddr(), fooDenom)
	assert.True(t, balance.IsZero())

	erc20 := contracts.ERC20MinterBurnerDecimalsContract.ABI
	ercBalance := app.Erc20Keeper.BalanceOf(ctx, erc20, pair.GetERC20Contract(), common.BytesToAddress(sdk.AccAddress(valAddrs[0]).Bytes()))
	assert.NotNil(t, ercBalance)
	assert.True(t, math.NewIntFromBigInt(ercBalance).IsPositive())
}
