package keeper_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	chain "github.com/dymensionxyz/dymension-rdk/testutil/app"
	testutils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/utils/addressutils"
	"github.com/dymensionxyz/dymension-rdk/utils/sliceutils"
	"github.com/dymensionxyz/dymension-rdk/x/gasless"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
	minttypes "github.com/dymensionxyz/dymension-rdk/x/mint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app            *chain.App
	ctx            sdk.Context
	keeper         keeper.Keeper
	querier        keeper.Querier
	msgServer      types.MsgServer
	encodingConfig simappparams.EncodingConfig
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = testutils.Setup(s.T(), false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.keeper = s.app.GaslessKeeper
	s.querier = keeper.Querier{Keeper: s.keeper}
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
	s.encodingConfig = simapp.MakeTestEncodingConfig()
}

// Below are just shortcuts to frequently-used functions.
func (s *KeeperTestSuite) getBalances(addr sdk.AccAddress) sdk.Coins {
	return s.app.BankKeeper.GetAllBalances(s.ctx, addr)
}

func (s *KeeperTestSuite) getBalance(addr sdk.AccAddress, denom string) sdk.Coin {
	return s.app.BankKeeper.GetBalance(s.ctx, addr, denom)
}

func (s *KeeperTestSuite) sendCoins(fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) {
	s.T().Helper()
	err := s.app.BankKeeper.SendCoins(s.ctx, fromAddr, toAddr, amt)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) nextBlock() {
	gasless.EndBlocker(s.ctx, s.keeper)
	gasless.BeginBlocker(s.ctx, s.keeper)
}

// Below are useful helpers to write test code easily.
func (s *KeeperTestSuite) addr(addrNum int) sdk.AccAddress {
	return addressutils.DeriveAddress(addressutils.AddressType32Bytes, types.ModuleName, fmt.Sprintf("address-%d", addrNum))
}

func (s *KeeperTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coins) {
	// using mint module to mint new test coins, since gasless module is not allowed to mint coins
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, amt)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, addr, amt)
	s.Require().NoError(err)
}

func newInt(i int64) sdkmath.Int {
	return sdkmath.NewInt(i)
}

func newDec(i int64) sdkmath.LegacyDec {
	return sdkmath.LegacyNewDec(i)
}

func coinEq(exp, got sdk.Coin) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func coinsEq(exp, got sdk.Coins) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func intEq(exp, got sdkmath.Int) (bool, string, string, string) {
	return exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func decEq(exp, got sdkmath.LegacyDec) (bool, string, string, string) {
	return exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}

func (s *KeeperTestSuite) CreateNewGasTank(
	provider sdk.AccAddress,
	feeDenom string,
	maxFeeUsagePerTx sdkmath.Int,
	maxFeeUsagePerConsumer sdkmath.Int,
	usageIdentifiers []string,
	deposit string,
) types.GasTank {
	parsedDepositCoin := ParseCoin(deposit)
	s.fundAddr(provider, sdk.NewCoins(parsedDepositCoin))

	usageIdentifiers = sliceutils.RemoveDuplicates(usageIdentifiers)
	tank, err := s.keeper.CreateGasTank(s.ctx, types.NewMsgCreateGasTank(
		provider,
		feeDenom,
		maxFeeUsagePerTx,
		maxFeeUsagePerConsumer,
		usageIdentifiers,
		parsedDepositCoin,
	))
	s.Require().NoError(err)
	s.Require().IsType(types.GasTank{}, tank)
	s.Require().Equal(feeDenom, tank.FeeDenom)
	s.Require().Equal(maxFeeUsagePerTx, tank.MaxFeeUsagePerTx)
	s.Require().Equal(maxFeeUsagePerConsumer, tank.MaxFeeUsagePerConsumer)
	s.Require().Equal(usageIdentifiers, tank.UsageIdentifiers)
	s.Require().Equal(ParseCoin(deposit), s.getBalance(tank.GetGasTankReserveAddress(), feeDenom))

	for _, identifier := range usageIdentifiers {
		uiToGTIDs, found := s.keeper.GetUsageIdentifierToGasTankIds(s.ctx, identifier)
		s.Require().True(found)
		s.Require().IsType(types.UsageIdentifierToGasTankIds{}, uiToGTIDs)
		s.Require().IsType([]uint64{}, uiToGTIDs.GasTankIds)
		s.Require().Equal(uiToGTIDs.UsageIdentifier, identifier)
		s.Require().Equal(tank.Id, uiToGTIDs.GasTankIds[len(uiToGTIDs.GasTankIds)-1])
	}
	return tank
}

// ParseCoin parses and returns sdk.Coin.
func ParseCoin(s string) sdk.Coin {
	coin, err := sdk.ParseCoinNormalized(strings.ReplaceAll(s, "_", ""))
	if err != nil {
		panic(err)
	}
	return coin
}
