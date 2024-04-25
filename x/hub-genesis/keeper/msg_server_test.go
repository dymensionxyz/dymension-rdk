package keeper_test

import (
	"fmt"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/testutil/ibctest"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

const (
	rollappDenom = "arax"
)

var (
	authorisedAddress     = utils.AccAddress()
	initialRollappBalance = func() sdk.Coin {
		n, _ := sdk.NewIntFromString("100000000000000000000")
		return sdk.NewCoin(rollappDenom, n)
	}()
)

type HubGenesisMsgServerTestSuite struct {
	ibctest.IBCTestUtilSuite

	app       *app.App
	k         *keeper.Keeper
	msgServer types.MsgServer
	ctx       sdk.Context
}

func TestHubGenesisMsgServerTestSuite(t *testing.T) {
	suite.Run(t, new(HubGenesisMsgServerTestSuite))
}

func (suite *HubGenesisMsgServerTestSuite) setupTest() {
	suite.IBCTestUtilSuite.SetupTest(rollappDenom)
	suite.app = suite.RollAppChain.App.(*app.App)
	suite.k, suite.ctx = keepers.NewTestHubGenesisKeeperFromApp(suite.app)
	suite.msgServer = keeper.NewMsgServerImpl(*suite.k)
}

func (suite *HubGenesisMsgServerTestSuite) TestTriggerGenesisEvent() {
	suite.setupTest()
	path := suite.NewTransferPath(suite.RollAppChain, suite.HubChain)
	suite.Coordinator.Setup(path)

	cases := []struct {
		name                 string
		genesisState         *types.GenesisState
		msg                  *types.MsgHubGenesisEvent
		rollappBalanceBefore sdk.Coin
		expErr               error
		runBefore            func()
	}{
		{
			name: "successful hub genesis event",
			genesisState: &types.GenesisState{
				Params: types.Params{
					GenesisTriggererAllowlist: []types.GenesisTriggererParams{{Address: authorisedAddress.String()}},
				},
				State: types.State{GenesisTokens: sdk.NewCoins(initialRollappBalance)},
			},
			msg: &types.MsgHubGenesisEvent{
				Address:   authorisedAddress.String(),
				ChannelId: path.EndpointA.ChannelID,
				HubId:     path.EndpointB.Chain.ChainID,
			},
			rollappBalanceBefore: initialRollappBalance,
			expErr:               nil,
		}, {
			name: "invalid rollapp genesis event - genesis event already triggered",
			genesisState: &types.GenesisState{
				Params: types.Params{
					GenesisTriggererAllowlist: []types.GenesisTriggererParams{{Address: authorisedAddress.String()}},
				},
				State: types.State{IsLocked: true, GenesisTokens: sdk.NewCoins(initialRollappBalance)},
			},
			msg: &types.MsgHubGenesisEvent{
				Address:   authorisedAddress.String(),
				ChannelId: path.EndpointA.ChannelID,
				HubId:     path.EndpointB.Chain.ChainID,
			},
			rollappBalanceBefore: initialRollappBalance,
			expErr:               types.ErrGenesisEventAlreadyTriggered,
		}, {
			name: "invalid rollapp genesis event - address not in Allowlist",
			genesisState: &types.GenesisState{
				Params: types.Params{
					GenesisTriggererAllowlist: []types.GenesisTriggererParams{{Address: utils.AccAddress().String()}},
				},
				State: types.State{GenesisTokens: sdk.NewCoins(initialRollappBalance)},
			},
			msg: &types.MsgHubGenesisEvent{
				Address:   authorisedAddress.String(),
				ChannelId: path.EndpointA.ChannelID,
				HubId:     path.EndpointB.Chain.ChainID,
			},
			rollappBalanceBefore: initialRollappBalance,
			expErr:               sdkerrors.ErrUnauthorized,
		}, {
			name: "invalid rollapp genesis event - invalid channel id",
			genesisState: &types.GenesisState{
				Params: types.Params{
					GenesisTriggererAllowlist: []types.GenesisTriggererParams{{Address: authorisedAddress.String()}},
				},
				State: types.State{GenesisTokens: sdk.NewCoins(initialRollappBalance)},
			},
			msg: &types.MsgHubGenesisEvent{
				Address:   authorisedAddress.String(),
				ChannelId: "invalid-channel",
				HubId:     path.EndpointB.Chain.ChainID,
			},
			rollappBalanceBefore: initialRollappBalance,
			expErr:               types.ErrFailedGetClientState,
		}, {
			name: "invalid rollapp genesis event - invalid chain id",
			genesisState: &types.GenesisState{
				Params: types.Params{
					GenesisTriggererAllowlist: []types.GenesisTriggererParams{{Address: authorisedAddress.String()}},
				},
				State: types.State{GenesisTokens: sdk.NewCoins(initialRollappBalance)},
			},
			msg: &types.MsgHubGenesisEvent{
				Address:   authorisedAddress.String(),
				ChannelId: path.EndpointA.ChannelID,
				HubId:     "invalid-chain-id",
			},
			rollappBalanceBefore: initialRollappBalance,
			expErr:               types.ErrChainIDMismatch,
		}, {
			name: "invalid rollapp genesis event - module account has no coins",
			genesisState: &types.GenesisState{
				Params: types.Params{
					GenesisTriggererAllowlist: []types.GenesisTriggererParams{{Address: authorisedAddress.String()}},
				},
				State: types.State{GenesisTokens: sdk.NewCoins(initialRollappBalance)},
			},
			msg: &types.MsgHubGenesisEvent{
				Address:   authorisedAddress.String(),
				ChannelId: path.EndpointA.ChannelID,
				HubId:     path.EndpointB.Chain.ChainID,
			},
			rollappBalanceBefore: sdk.NewCoin(rollappDenom, sdk.NewInt(0)),
			runBefore: func() {
				// remove all coins from the module account
				err := suite.app.BankKeeper.BurnCoins(suite.ctx, types.ModuleName, sdk.Coins{initialRollappBalance})
				suite.Require().NoError(err)
			},
			expErr: types.ErrLockingGenesisTokens,
		},

		// TODO: add test of not enough tokens
	}

	for _, tc := range cases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			defer func() {
				suite.setupTest()
				path = suite.NewTransferPath(suite.RollAppChain, suite.HubChain)
				suite.Coordinator.Setup(path)
			}()

			if tc.runBefore != nil {
				tc.runBefore()
			}

			suite.k.SetState(suite.ctx, tc.genesisState.State)
			suite.k.SetParams(suite.ctx, tc.genesisState.Params)
			moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)

			// check the initial module balance
			rollappBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, rollappDenom)
			suite.Require().Equal(tc.rollappBalanceBefore, rollappBalanceBefore, tc.name)

			// trigger the genesis event
			_, err := suite.msgServer.TriggerGenesisEvent(suite.ctx, tc.msg)
			suite.Require().ErrorIs(err, tc.expErr, tc.name)

			var (
				expectedBalance       sdk.Coin
				expectedEscrowBalance sdk.Coin
				expectedState         bool
			)

			if tc.expErr == nil {
				expectedBalance = sdk.NewCoin(rollappDenom, sdk.NewInt(0))
				expectedEscrowBalance = tc.rollappBalanceBefore
				expectedState = true
			} else {
				expectedBalance = tc.rollappBalanceBefore
				expectedEscrowBalance = sdk.NewCoin(rollappDenom, sdk.NewInt(0))
				expectedState = false || tc.genesisState.State.IsLocked
			}

			// check the hub genesis state
			stateState := suite.k.GetState(suite.ctx)
			suite.Require().Equal(expectedState, stateState.IsLocked, tc.name)

			// check the module balance after the genesis event
			rollappBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, rollappDenom)
			suite.Require().True(expectedBalance.IsEqual(rollappBalanceAfter), tc.name)

			// check the escrow balance after the genesis event
			sourceChannel := path.EndpointA.ChannelID
			escrowAddress := ibctypes.GetEscrowAddress("transfer", sourceChannel)
			escrowBalance := suite.app.BankKeeper.GetBalance(suite.ctx, escrowAddress, rollappDenom)
			suite.Require().True(expectedEscrowBalance.IsEqual(escrowBalance), tc.name)
		})
	}
}
