package ante_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	"github.com/dymensionxyz/dymension-rdk/server/ante"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/utils/erc20"
)

var terminatorAnteHandler = func(ctx sdk.Context, _ sdk.Tx, simulate bool) (sdk.Context, error) {
	return ctx, nil
}

// var terminatorAnteHandler sdk.AnteHandler

func (s *AnteTestSuite) TestERC20ConvertDecorator_Staking_ConvertFromERC20IfNeeded(t *testing.T) {
	stakeAmount := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	fooDenom := "foo"

	cleanCtx := s.ctx
	tc := []struct {
		name   string
		setup  func(sdk.Context, sdk.AccAddress)
		expErr bool
	}{
		{
			name: "create validator success - enough bank balance",
			setup: func(ctx sdk.Context, addr sdk.AccAddress) {
				s.FundAccount(addr, sdk.NewCoin("foo", stakeAmount))
			},
		},
		{
			name: "create validator success - not enough bank balance, pulled from ERC20",
			setup: func(ctx sdk.Context, addr sdk.AccAddress) {
				// fund the account with ERC20 tokens
				s.FundAccount(addr, sdk.NewCoin("foo", stakeAmount))
				err := erc20.ConvertCoin(ctx, s.app.Erc20Keeper, sdk.NewCoin("foo", stakeAmount), addr, addr)
				s.NoError(err)

				balance := s.app.BankKeeper.GetBalance(ctx, addr, fooDenom)
				s.True(balance.IsZero())
			},
		},
		{
			name:   "create validator fails - not enough balance",
			expErr: true,
		},
	}

	for _, tc := range tc {
		s.ctx = cleanCtx
		tstaking := teststaking.NewHelper(t, s.ctx, s.app.StakingKeeper.Keeper)
		tstaking.Denom = fooDenom
		t.Run(tc.name, func(t *testing.T) {
			addr := utils.AccAddress()
			if tc.setup != nil {
				tc.setup(s.ctx, addr)
			}

			builder := s.app.GetTxConfig().NewTxBuilder()
			msg := tstaking.CreateValidatorMsg(sdk.ValAddress(addr), ed25519.GenPrivKey().PubKey(), stakeAmount)
			err := builder.SetMsgs(msg)
			s.NoError(err)
			tx := builder.GetTx()

			decorator := ante.NewERC20ConversionDecorator(s.app.Erc20Keeper, s.app.BankKeeper)
			_, err = decorator.AnteHandle(s.ctx, tx, false, terminatorAnteHandler)
			if tc.expErr {
				s.Error(err)
				return
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *AnteTestSuite) TestERC20ConvertPostDecorator(t *testing.T) {
	stakeAmount := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	addr := utils.AccAddress()

	// Set fees to the fee account
	fees := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(100)))
	s.FundFees(fees)

	// mint some balance on the user bank account. we expect it to be converted to ERC20 as well
	s.FundAccount(addr, sdk.NewCoin("foo", stakeAmount))

	// Generate drawRewards tx
	drawRewardsMsg := disttypes.NewMsgWithdrawDelegatorReward(addr, sdk.ValAddress(addr))
	builder := s.app.GetTxConfig().NewTxBuilder()
	err := builder.SetMsgs(drawRewardsMsg)
	s.NoError(err)
	tx := builder.GetTx()

	// Call post handler
	postDecorator := ante.NewERC20ConversionPostHandlerDecorator(s.app.Erc20Keeper, s.app.BankKeeper)
	_, err = postDecorator.AnteHandle(s.ctx, tx, false, terminatorAnteHandler)
	s.NoError(err)

	// Check that the balance has been converted to ERC20
	balance := s.app.BankKeeper.GetBalance(s.ctx, addr, "foo")
	s.True(balance.IsZero())
}

func (s *AnteTestSuite) TestERC20ConvertPostDecorator_VestingAccount(t *testing.T) {
	stakeAmount := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)

	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	baseAcc := authtypes.NewBaseAccount(addr, pubkey, 0, 0)

	vestingCoin := sdk.NewCoin("foo", stakeAmount)
	vestingAcc := vestingtypes.NewContinuousVestingAccount(
		baseAcc,
		sdk.NewCoins(vestingCoin),
		100,
		200,
	)
	s.app.AccountKeeper.SetAccount(s.ctx, vestingAcc)
	// we fund the vesting account with 2x the amount of tokens
	// half of the tokens will be vested
	s.FundAccount(addr, sdk.NewCoin("foo", stakeAmount.MulRaw(2)))

	// Set fees to the fee account
	fees := sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(100)))
	s.FundFees(fees)

	// Generate drawRewards tx
	drawRewardsMsg := disttypes.NewMsgWithdrawDelegatorReward(addr, sdk.ValAddress(addr))
	builder := s.app.GetTxConfig().NewTxBuilder()
	err := builder.SetMsgs(drawRewardsMsg)
	s.NoError(err)
	tx := builder.GetTx()

	// Call post handler
	postDecorator := ante.NewERC20ConversionPostHandlerDecorator(s.app.Erc20Keeper, s.app.BankKeeper)
	_, err = postDecorator.AnteHandle(s.ctx, tx, false, terminatorAnteHandler)
	s.NoError(err)

	// Check that the balance has been converted to ERC20
	balance := s.app.BankKeeper.GetBalance(s.ctx, addr, "foo")
	s.Equal(stakeAmount, balance.Amount)
}
