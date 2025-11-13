package converter_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"
	"github.com/stretchr/testify/suite"

	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/testutil/ibctest"
	hubtypes "github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

const (
	hubDenom        = "uusdc"   // Source denom on hub
	bridgeDecimals  = 6         // Decimals for the token on hub
	rollappDecimals = 18        // Standard decimals on rollapp
	testAmount      = "1000000" // 1 USDC with 6 decimals
)

// MiddlewareTestSuite tests the decimal conversion middleware using actual IBC flows
type MiddlewareTestSuite struct {
	ibctest.IBCTestUtilSuite

	path       *ibctesting.Path
	hubApp     *ibctesting.TestingApp
	rollappApp *app.App
}

// TestMiddlewareTestSuite runs the test suite
func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

// SetupTest initializes the test suite with two chains and an IBC channel
func (suite *MiddlewareTestSuite) SetupTest() {
	// Setup chains using the IBCTestUtilSuite
	suite.IBCTestUtilSuite.SetupTest("stake") // rollapp native denom

	// Get the rollapp application
	var ok bool
	suite.rollappApp, ok = suite.RollAppChain.App.(*app.App)
	suite.Require().True(ok, "failed to cast rollapp to *app.App")

	// Store hub app reference
	suite.hubApp = &suite.HubChain.App

	// Create and setup IBC path between hub and rollapp
	suite.path = suite.NewTransferPath(suite.HubChain, suite.RollAppChain)
	suite.Coordinator.Setup(suite.path)
}

// setupDecimalConversion configures the decimal conversion on the rollapp for a given IBC denom
func (suite *MiddlewareTestSuite) setupDecimalConversion(ibcDenom string) {
	ctx := suite.RollAppChain.GetContext()

	// Register the IBC denom metadata with 18 decimals (rollapp standard)
	metadata := banktypes.Metadata{
		Base: ibcDenom,
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: ibcDenom, Exponent: 0},
			{Denom: "usdc", Exponent: rollappDecimals},
		},
		Display: "usdc",
	}
	suite.rollappApp.BankKeeper.SetDenomMetaData(ctx, metadata)

	// Set up the decimal conversion pair in hub keeper
	pair := hubtypes.DecimalConversionPair{
		FromToken:    ibcDenom,
		FromDecimals: bridgeDecimals,
	}
	err := suite.rollappApp.HubKeeper.SetDecimalConversionPair(ctx, pair)
	suite.Require().NoError(err)
}

// enableOutboundTransfers enables outbound IBC transfers on the rollapp
func (suite *MiddlewareTestSuite) enableOutboundTransfers() {
	ctx := suite.RollAppChain.GetContext()
	state := suite.rollappApp.HubGenesisKeeper.GetState(ctx)
	state.OutboundTransfersEnabled = true
	suite.rollappApp.HubGenesisKeeper.SetState(ctx, state)
}

// getIBCDenom returns the IBC denom hash for the hub denom as it appears on rollapp
func (suite *MiddlewareTestSuite) getIBCDenom() string {
	// The IBC denom is a hash of "transfer/{channel-id}/{denom}"
	sourcePrefix := transfertypes.GetDenomPrefix(suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID)
	prefixedDenom := sourcePrefix + hubDenom
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)
	return denomTrace.IBCDenom()
}

// TestOnRecvPacket_ConversionWithActualIBCFlow tests that receiving an IBC transfer
// from the hub correctly converts the amount from 6 to 18 decimals
func (suite *MiddlewareTestSuite) TestOnRecvPacket_ConversionWithActualIBCFlow() {
	// Get receiver address on rollapp
	receiver := suite.RollAppChain.SenderAccount.GetAddress()

	// Get the IBC denom that will be created on rollapp
	ibcDenom := suite.getIBCDenom()

	// Setup decimal conversion for this IBC denom
	suite.setupDecimalConversion(ibcDenom)

	// Fund hub sender with tokens
	hubSender := suite.HubChain.SenderAccount.GetAddress()
	hubCoins := sdk.NewCoins(sdk.NewCoin(hubDenom, math.NewInt(1000000))) // 1 USDC with 6 decimals
	err := suite.HubChain.GetSimApp().BankKeeper.MintCoins(suite.HubChain.GetContext(), transfertypes.ModuleName, hubCoins)
	suite.Require().NoError(err)
	err = suite.HubChain.GetSimApp().BankKeeper.SendCoinsFromModuleToAccount(suite.HubChain.GetContext(), transfertypes.ModuleName, hubSender, hubCoins)
	suite.Require().NoError(err)

	// Send IBC transfer from hub to rollapp
	msg := transfertypes.NewMsgTransfer(
		suite.path.EndpointA.ChannelConfig.PortID,
		suite.path.EndpointA.ChannelID,
		hubCoins[0],
		hubSender.String(),
		receiver.String(),
		clienttypes.NewHeight(10, 1000),
		0,
		"",
	)

	res, err := suite.HubChain.SendMsgs(msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	// Relay the packet to rollapp
	packet, err := ibctesting.ParsePacketFromEvents(res.GetEvents())
	suite.Require().NoError(err)

	err = suite.path.RelayPacket(packet)
	suite.Require().NoError(err)

	// Verify the receiver got the converted amount (18 decimals)
	balance := suite.rollappApp.BankKeeper.GetBalance(suite.RollAppChain.GetContext(), receiver, ibcDenom)
	expectedAmount, ok := math.NewIntFromString("1000000000000000000") // 1 USDC with 18 decimals
	suite.Require().True(ok)
	suite.Require().Equal(expectedAmount.String(), balance.Amount.String(),
		"receiver should have 1 USDC with 18 decimals")
}

// TestOnRecvPacket_NoConversionRequired tests that tokens without conversion
// configured pass through unchanged
func (suite *MiddlewareTestSuite) TestOnRecvPacket_NoConversionRequired() {
	// Use a different token that doesn't have conversion configured
	otherDenom := "uatom"

	// Get receiver address on rollapp
	receiver := suite.RollAppChain.SenderAccount.GetAddress()

	// Fund hub sender with tokens
	hubSender := suite.HubChain.SenderAccount.GetAddress()
	hubCoins := sdk.NewCoins(sdk.NewCoin(otherDenom, math.NewInt(5000000)))
	err := suite.HubChain.GetSimApp().BankKeeper.MintCoins(suite.HubChain.GetContext(), transfertypes.ModuleName, hubCoins)
	suite.Require().NoError(err)
	err = suite.HubChain.GetSimApp().BankKeeper.SendCoinsFromModuleToAccount(suite.HubChain.GetContext(), transfertypes.ModuleName, hubSender, hubCoins)
	suite.Require().NoError(err)

	// Send IBC transfer from hub to rollapp
	msg := transfertypes.NewMsgTransfer(
		suite.path.EndpointA.ChannelConfig.PortID,
		suite.path.EndpointA.ChannelID,
		hubCoins[0],
		hubSender.String(),
		receiver.String(),
		clienttypes.NewHeight(10, 1000),
		0,
		"",
	)

	res, err := suite.HubChain.SendMsgs(msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	// Relay the packet to rollapp
	packet, err := ibctesting.ParsePacketFromEvents(res.GetEvents())
	suite.Require().NoError(err)

	err = suite.path.RelayPacket(packet)
	suite.Require().NoError(err)

	// Get the IBC denom for the other token
	sourcePrefix := transfertypes.GetDenomPrefix(suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID)
	prefixedDenom := sourcePrefix + otherDenom
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)
	ibcDenom := denomTrace.IBCDenom()

	// Verify the receiver got the exact same amount (no conversion)
	balance := suite.rollappApp.BankKeeper.GetBalance(suite.RollAppChain.GetContext(), receiver, ibcDenom)
	suite.Require().Equal(hubCoins[0].Amount.String(), balance.Amount.String(),
		"receiver should have exact amount without conversion")
}

// TestMultipleTransfers tests multiple sequential transfers with conversion
func (suite *MiddlewareTestSuite) TestMultipleTransfers() {
	// Get the IBC denom
	ibcDenom := suite.getIBCDenom()

	// Setup decimal conversion
	suite.setupDecimalConversion(ibcDenom)

	// Get addresses
	receiver := suite.RollAppChain.SenderAccount.GetAddress()
	hubSender := suite.HubChain.SenderAccount.GetAddress()

	// Test multiple transfers with different amounts
	testAmounts := []math.Int{
		math.NewInt(1000000),   // 1 USDC
		math.NewInt(5000000),   // 5 USDC
		math.NewInt(100000000), // 100 USDC
	}

	expectedTotal := math.ZeroInt()

	for i, amount := range testAmounts {
		// Fund hub sender
		hubCoins := sdk.NewCoins(sdk.NewCoin(hubDenom, amount))
		err := suite.HubChain.GetSimApp().BankKeeper.MintCoins(suite.HubChain.GetContext(), transfertypes.ModuleName, hubCoins)
		suite.Require().NoError(err, "transfer %d", i)
		err = suite.HubChain.GetSimApp().BankKeeper.SendCoinsFromModuleToAccount(suite.HubChain.GetContext(), transfertypes.ModuleName, hubSender, hubCoins)
		suite.Require().NoError(err, "transfer %d", i)

		// Send IBC transfer
		msg := transfertypes.NewMsgTransfer(
			suite.path.EndpointA.ChannelConfig.PortID,
			suite.path.EndpointA.ChannelID,
			hubCoins[0],
			hubSender.String(),
			receiver.String(),
			clienttypes.NewHeight(10, 1000),
			0,
			"",
		)

		res, err := suite.HubChain.SendMsgs(msg)
		suite.Require().NoError(err, "transfer %d", i)

		packet, err := ibctesting.ParsePacketFromEvents(res.GetEvents())
		suite.Require().NoError(err, "transfer %d", i)

		err = suite.path.RelayPacket(packet)
		suite.Require().NoError(err, "transfer %d", i)

		// Calculate expected amount with conversion (multiply by 10^12)
		convertedAmount := amount.Mul(math.NewInt(1_000_000_000_000))
		expectedTotal = expectedTotal.Add(convertedAmount)
	}

	// Verify total balance
	finalBalance := suite.rollappApp.BankKeeper.GetBalance(suite.RollAppChain.GetContext(), receiver, ibcDenom)
	suite.Require().Equal(expectedTotal.String(), finalBalance.Amount.String(),
		"receiver should have sum of all converted amounts")
}

// TestOnSendPacket_ConversionWithActualIBCFlow tests that sending an IBC transfer
// back to the hub correctly converts the amount from 18 to 6 decimals
func (suite *MiddlewareTestSuite) TestOnSendPacket_ConversionWithActualIBCFlow() {
	// First, receive tokens from hub with conversion
	ibcDenom := suite.getIBCDenom()
	suite.setupDecimalConversion(ibcDenom)
	suite.enableOutboundTransfers()

	// Get addresses
	rollappSender := suite.RollAppChain.SenderAccount.GetAddress()
	hubReceiver := suite.HubChain.SenderAccount.GetAddress()

	// Fund hub and send to rollapp first
	hubCoins := sdk.NewCoins(sdk.NewCoin(hubDenom, math.NewInt(10000000))) // 10 USDC with 6 decimals
	err := suite.HubChain.GetSimApp().BankKeeper.MintCoins(suite.HubChain.GetContext(), transfertypes.ModuleName, hubCoins)
	suite.Require().NoError(err)
	err = suite.HubChain.GetSimApp().BankKeeper.SendCoinsFromModuleToAccount(suite.HubChain.GetContext(), transfertypes.ModuleName, suite.HubChain.SenderAccount.GetAddress(), hubCoins)
	suite.Require().NoError(err)

	// Send from hub to rollapp
	msg := transfertypes.NewMsgTransfer(
		suite.path.EndpointA.ChannelConfig.PortID,
		suite.path.EndpointA.ChannelID,
		hubCoins[0],
		suite.HubChain.SenderAccount.GetAddress().String(),
		rollappSender.String(),
		clienttypes.NewHeight(10, 1000),
		0,
		"",
	)

	res, err := suite.HubChain.SendMsgs(msg)
	suite.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(res.GetEvents())
	suite.Require().NoError(err)

	err = suite.path.RelayPacket(packet)
	suite.Require().NoError(err)

	// Verify rollapp user received converted amount (18 decimals)
	balance := suite.rollappApp.BankKeeper.GetBalance(suite.RollAppChain.GetContext(), rollappSender, ibcDenom)
	expectedRollappAmount, ok := math.NewIntFromString("10000000000000000000") // 10 USDC with 18 decimals
	suite.Require().True(ok)
	suite.Require().Equal(expectedRollappAmount.String(), balance.Amount.String())

	// Now send back to hub with conversion (18 decimals -> 6 decimals)
	// Send 5 USDC back (5 * 10^18)
	sendBackAmount, ok := math.NewIntFromString("5000000000000000000") // 5 USDC with 18 decimals
	suite.Require().True(ok)
	sendBackMsg := transfertypes.NewMsgTransfer(
		suite.path.EndpointB.ChannelConfig.PortID,
		suite.path.EndpointB.ChannelID,
		sdk.NewCoin(ibcDenom, sendBackAmount),
		rollappSender.String(),
		hubReceiver.String(),
		clienttypes.NewHeight(10, 1000),
		0,
		"",
	)

	res2, err := suite.RollAppChain.SendMsgs(sendBackMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res2)

	// Relay the packet back to hub
	packet2, err := ibctesting.ParsePacketFromEvents(res2.GetEvents())
	suite.Require().NoError(err)

	err = suite.path.RelayPacket(packet2)
	suite.Require().NoError(err)

	// Verify hub receiver got the converted amount (6 decimals)
	// The packet should contain 5000000 (5 USDC with 6 decimals)
	hubBalance := suite.HubChain.GetSimApp().BankKeeper.GetBalance(suite.HubChain.GetContext(), hubReceiver, hubDenom)
	expectedHubAmount := math.NewInt(5000000) // 5 USDC with 6 decimals
	suite.Require().Equal(expectedHubAmount.String(), hubBalance.Amount.String(),
		"hub receiver should have 5 USDC with 6 decimals")

	// Verify rollapp sender has the remaining amount
	remainingBalance := suite.rollappApp.BankKeeper.GetBalance(suite.RollAppChain.GetContext(), rollappSender, ibcDenom)
	expectedRemaining := expectedRollappAmount.Sub(sendBackAmount) // 5 USDC with 18 decimals
	suite.Require().Equal(expectedRemaining.String(), remainingBalance.Amount.String(),
		"rollapp sender should have 5 USDC with 18 decimals remaining")
}

// TestOnSendPacket_NoConversionRequired tests that sending tokens without conversion
// configured pass through unchanged
func (suite *MiddlewareTestSuite) TestOnSendPacket_NoConversionRequired() {
	// Use a different token that doesn't have conversion configured
	otherDenom := "uatom"
	suite.enableOutboundTransfers()

	// Get addresses
	rollappSender := suite.RollAppChain.SenderAccount.GetAddress()
	hubReceiver := suite.HubChain.SenderAccount.GetAddress()

	// Fund hub sender and send to rollapp
	hubCoins := sdk.NewCoins(sdk.NewCoin(otherDenom, math.NewInt(7000000)))
	err := suite.HubChain.GetSimApp().BankKeeper.MintCoins(suite.HubChain.GetContext(), transfertypes.ModuleName, hubCoins)
	suite.Require().NoError(err)
	err = suite.HubChain.GetSimApp().BankKeeper.SendCoinsFromModuleToAccount(suite.HubChain.GetContext(), transfertypes.ModuleName, suite.HubChain.SenderAccount.GetAddress(), hubCoins)
	suite.Require().NoError(err)

	// Send from hub to rollapp
	msg := transfertypes.NewMsgTransfer(
		suite.path.EndpointA.ChannelConfig.PortID,
		suite.path.EndpointA.ChannelID,
		hubCoins[0],
		suite.HubChain.SenderAccount.GetAddress().String(),
		rollappSender.String(),
		clienttypes.NewHeight(10, 1000),
		0,
		"",
	)

	res, err := suite.HubChain.SendMsgs(msg)
	suite.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(res.GetEvents())
	suite.Require().NoError(err)

	err = suite.path.RelayPacket(packet)
	suite.Require().NoError(err)

	// Get the IBC denom on rollapp
	sourcePrefix := transfertypes.GetDenomPrefix(suite.path.EndpointB.ChannelConfig.PortID, suite.path.EndpointB.ChannelID)
	prefixedDenom := sourcePrefix + otherDenom
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)
	ibcDenom := denomTrace.IBCDenom()

	// Verify rollapp user received exact same amount
	balance := suite.rollappApp.BankKeeper.GetBalance(suite.RollAppChain.GetContext(), rollappSender, ibcDenom)
	suite.Require().Equal(hubCoins[0].Amount.String(), balance.Amount.String())

	// Now send back to hub (should be no conversion)
	sendBackAmount := math.NewInt(3000000)
	sendBackMsg := transfertypes.NewMsgTransfer(
		suite.path.EndpointB.ChannelConfig.PortID,
		suite.path.EndpointB.ChannelID,
		sdk.NewCoin(ibcDenom, sendBackAmount),
		rollappSender.String(),
		hubReceiver.String(),
		clienttypes.NewHeight(10, 1000),
		0,
		"",
	)

	res2, err := suite.RollAppChain.SendMsgs(sendBackMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res2)

	// Relay the packet back to hub
	packet2, err := ibctesting.ParsePacketFromEvents(res2.GetEvents())
	suite.Require().NoError(err)

	err = suite.path.RelayPacket(packet2)
	suite.Require().NoError(err)

	// Verify hub receiver got the exact same amount (no conversion)
	hubBalance := suite.HubChain.GetSimApp().BankKeeper.GetBalance(suite.HubChain.GetContext(), hubReceiver, otherDenom)
	suite.Require().Equal(sendBackAmount.String(), hubBalance.Amount.String(),
		"hub receiver should have exact amount without conversion")

	// Verify rollapp sender has the remaining amount
	remainingBalance := suite.rollappApp.BankKeeper.GetBalance(suite.RollAppChain.GetContext(), rollappSender, ibcDenom)
	expectedRemaining := hubCoins[0].Amount.Sub(sendBackAmount)
	suite.Require().Equal(expectedRemaining.String(), remainingBalance.Amount.String(),
		"rollapp sender should have remaining tokens")
}

// TestOnSendPacket_NativeRollappTokenNoConversion tests sending native rollapp tokens
// to the hub without any conversion configured
func (suite *MiddlewareTestSuite) TestOnSendPacket_NativeRollappTokenNoConversion() {
	suite.enableOutboundTransfers()

	// Get addresses
	rollappSender := suite.RollAppChain.SenderAccount.GetAddress()
	hubReceiver := suite.HubChain.SenderAccount.GetAddress()

	// Use the native rollapp token (stake) - no conversion configured
	nativeDenom := "stake"
	sendAmount := math.NewInt(10000000) // 10 stake tokens

	// Fund rollapp sender with native tokens (they should already have some from chain setup)
	rollappCoins := sdk.NewCoins(sdk.NewCoin(nativeDenom, sendAmount))
	err := suite.rollappApp.BankKeeper.MintCoins(suite.RollAppChain.GetContext(), transfertypes.ModuleName, rollappCoins)
	suite.Require().NoError(err)
	err = suite.rollappApp.BankKeeper.SendCoinsFromModuleToAccount(suite.RollAppChain.GetContext(), transfertypes.ModuleName, rollappSender, rollappCoins)
	suite.Require().NoError(err)

	// Record initial balance
	initialBalance := suite.rollappApp.BankKeeper.GetBalance(suite.RollAppChain.GetContext(), rollappSender, nativeDenom)

	// Send native rollapp tokens to hub (should be no conversion)
	transferAmount := math.NewInt(3000000) // Send 3 stake tokens
	sendMsg := transfertypes.NewMsgTransfer(
		suite.path.EndpointB.ChannelConfig.PortID,
		suite.path.EndpointB.ChannelID,
		sdk.NewCoin(nativeDenom, transferAmount),
		rollappSender.String(),
		hubReceiver.String(),
		clienttypes.NewHeight(10, 1000),
		0,
		"",
	)

	res, err := suite.RollAppChain.SendMsgs(sendMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	// Relay the packet to hub
	packet, err := ibctesting.ParsePacketFromEvents(res.GetEvents())
	suite.Require().NoError(err)

	err = suite.path.RelayPacket(packet)
	suite.Require().NoError(err)

	// Calculate the IBC denom on hub (transfer/channel-0/stake)
	sourcePrefix := transfertypes.GetDenomPrefix(suite.path.EndpointA.ChannelConfig.PortID, suite.path.EndpointA.ChannelID)
	prefixedDenom := sourcePrefix + nativeDenom
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)
	ibcDenomOnHub := denomTrace.IBCDenom()

	// Verify hub receiver got the exact same amount (no conversion)
	hubBalance := suite.HubChain.GetSimApp().BankKeeper.GetBalance(suite.HubChain.GetContext(), hubReceiver, ibcDenomOnHub)
	suite.Require().Equal(transferAmount.String(), hubBalance.Amount.String(),
		"hub receiver should have exact amount of IBC-wrapped rollapp tokens without conversion")

	// Verify rollapp sender balance decreased by transfer amount
	finalBalance := suite.rollappApp.BankKeeper.GetBalance(suite.RollAppChain.GetContext(), rollappSender, nativeDenom)
	expectedFinal := initialBalance.Amount.Sub(transferAmount)
	suite.Require().Equal(expectedFinal.String(), finalBalance.Amount.String(),
		"rollapp sender balance should decrease by transfer amount")
}
