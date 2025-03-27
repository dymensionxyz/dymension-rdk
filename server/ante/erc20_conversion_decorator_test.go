package ante_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	"github.com/dymensionxyz/dymension-rdk/server/ante"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/utils/erc20"
	"github.com/stretchr/testify/require"
)

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
				err := erc20.ConvertCoin(ctx, s.app.Erc20Keeper, sdk.NewCoin("foo", stakeAmount), addr)
				require.NoError(t, err)

				balance := s.app.BankKeeper.GetBalance(ctx, addr, fooDenom)
				require.True(t, balance.IsZero())
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

			msg := tstaking.CreateValidatorMsg(sdk.ValAddress(addr), ed25519.GenPrivKey().PubKey(), stakeAmount)
			decorator := ante.NewERC20ConversionDecorator(s.app.Erc20Keeper, s.app.BankKeeper)

			builder := s.app.GetTxConfig().NewTxBuilder()
			err := builder.SetMsgs(msg)
			require.NoError(t, err)

			tx := builder.GetTx()

			var terminatorAnteHandler sdk.AnteHandler
			terminatorAnteHandler = func(ctx sdk.Context, _ sdk.Tx, simulate bool) (sdk.Context, error) {
				return ctx, nil
			}

			_, err = decorator.AnteHandle(s.ctx, tx, false, terminatorAnteHandler)
			if tc.expErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestERC20ConvertPostDecorator(t *testing.T) {

}
