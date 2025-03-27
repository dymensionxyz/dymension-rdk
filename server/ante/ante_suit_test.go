package ante_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	minttypes "github.com/dymensionxyz/dymension-rdk/x/mint/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type AnteTestSuite struct {
	suite.Suite

	app *app.App
	ctx sdk.Context
}

func (s *AnteTestSuite) SetupTest() {
	app := utils.Setup(s.T(), false)
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
	_, err := app.Erc20Keeper.RegisterCoin(ctx, fooMetadata)
	s.NoError(err)

	// set as staking denom
	params := app.StakingKeeper.GetParams(ctx)
	params.BondDenom = fooDenom
	app.StakingKeeper.SetParams(ctx, params)

	s.app = app
	s.ctx = ctx
}

func (s *AnteTestSuite) FundAccount(addr sdk.AccAddress, coin sdk.Coin) {
	err := s.app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, sdk.NewCoins(coin))
	s.NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, addr, sdk.NewCoins(coin))
	s.NoError(err)
}
