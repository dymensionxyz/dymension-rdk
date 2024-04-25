package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/governors/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
	"github.com/stretchr/testify/require"
)

func TestCancelUnbondingDelegation(t *testing.T) {
	// setup the app
	_, app, ctx := createTestInput(t)
	msgServer := keeper.NewMsgServerImpl(app.StakingKeeper)
	bondDenom := app.StakingKeeper.BondDenom(ctx)

	// set the not bonded pool module account
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	startTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 5)

	require.NoError(t, testutil.FundModuleAccount(app.BankKeeper, ctx, notBondedPool.GetName(), sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens))))
	app.AccountKeeper.SetModuleAccount(ctx, notBondedPool)

	moduleBalance := app.BankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), app.StakingKeeper.BondDenom(ctx))
	require.Equal(t, sdk.NewInt64Coin(bondDenom, startTokens.Int64()), moduleBalance)

	// accounts
	delAddrs := utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(10000))
	governors := app.StakingKeeper.GetGovernors(ctx, 10)
	require.Equal(t, len(governors), 1)

	governorAddr, err := sdk.ValAddressFromBech32(governors[0].OperatorAddress)
	require.NoError(t, err)
	delegatorAddr := delAddrs[0]

	// setting the ubd entry
	unbondingAmount := sdk.NewInt64Coin(app.StakingKeeper.BondDenom(ctx), 5)
	ubd := stakingtypes.NewUnbondingDelegation(
		delegatorAddr, governorAddr, 10,
		ctx.BlockTime().Add(time.Minute*10),
		unbondingAmount.Amount,
	)

	// set and retrieve a record
	app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
	resUnbond, found := app.StakingKeeper.GetUnbondingDelegation(ctx, delegatorAddr, governorAddr)
	require.True(t, found)
	require.Equal(t, ubd, resUnbond)

	testCases := []struct {
		Name      string
		ExceptErr bool
		req       types.MsgCancelUnbondingDelegation
	}{
		{
			Name:      "invalid height",
			ExceptErr: true,
			req: types.MsgCancelUnbondingDelegation{
				DelegatorAddress: resUnbond.DelegatorAddress,
				GovernorAddress:  resUnbond.ValidatorAddress,
				Amount:           sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), sdk.NewInt(4)),
				CreationHeight:   0,
			},
		},
		{
			Name:      "invalid coin",
			ExceptErr: true,
			req: types.MsgCancelUnbondingDelegation{
				DelegatorAddress: resUnbond.DelegatorAddress,
				GovernorAddress:  resUnbond.ValidatorAddress,
				Amount:           sdk.NewCoin("dump_coin", sdk.NewInt(4)),
				CreationHeight:   0,
			},
		},
		{
			Name:      "governor not exists",
			ExceptErr: true,
			req: types.MsgCancelUnbondingDelegation{
				DelegatorAddress: resUnbond.DelegatorAddress,
				GovernorAddress:  sdk.ValAddress(sdk.AccAddress("asdsad")).String(),
				Amount:           unbondingAmount,
				CreationHeight:   0,
			},
		},
		{
			Name:      "invalid delegator address",
			ExceptErr: true,
			req: types.MsgCancelUnbondingDelegation{
				DelegatorAddress: "invalid_delegator_addrtess",
				GovernorAddress:  resUnbond.ValidatorAddress,
				Amount:           unbondingAmount,
				CreationHeight:   0,
			},
		},
		{
			Name:      "invalid amount",
			ExceptErr: true,
			req: types.MsgCancelUnbondingDelegation{
				DelegatorAddress: resUnbond.DelegatorAddress,
				GovernorAddress:  resUnbond.ValidatorAddress,
				Amount:           unbondingAmount.Add(sdk.NewInt64Coin(bondDenom, 10)),
				CreationHeight:   10,
			},
		},
		{
			Name:      "success",
			ExceptErr: false,
			req: types.MsgCancelUnbondingDelegation{
				DelegatorAddress: resUnbond.DelegatorAddress,
				GovernorAddress:  resUnbond.ValidatorAddress,
				Amount:           unbondingAmount.Sub(sdk.NewInt64Coin(bondDenom, 1)),
				CreationHeight:   10,
			},
		},
		{
			Name:      "success",
			ExceptErr: false,
			req: types.MsgCancelUnbondingDelegation{
				DelegatorAddress: resUnbond.DelegatorAddress,
				GovernorAddress:  resUnbond.ValidatorAddress,
				Amount:           unbondingAmount.Sub(unbondingAmount.Sub(sdk.NewInt64Coin(bondDenom, 1))),
				CreationHeight:   10,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := msgServer.CancelUnbondingDelegation(ctx, &testCase.req)
			if testCase.ExceptErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				balanceForNotBondedPool := app.BankKeeper.GetBalance(ctx, sdk.AccAddress(notBondedPool.GetAddress()), bondDenom)
				require.Equal(t, balanceForNotBondedPool, moduleBalance.Sub(testCase.req.Amount))
				moduleBalance = moduleBalance.Sub(testCase.req.Amount)
			}
		})
	}
}
