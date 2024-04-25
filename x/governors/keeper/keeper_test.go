package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/governors/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *app.App
	ctx         sdk.Context
	addrs       []sdk.AccAddress
	vals        []types.Governor
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	// app := utils.Setup(suite.T(), false)
	_, app, ctx := createTestInput(suite.T())
	// ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	querier := keeper.Querier{Keeper: app.StakingKeeper}

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)
	queryClient := types.NewQueryClient(queryHelper)

	suite.msgServer = keeper.NewMsgServerImpl(app.StakingKeeper)

	addrs, _, governors := createGovernors(suite.T(), ctx, app, []int64{9, 8, 7})

	// sort a copy of the governors, so that original governors does not
	// have its order changed
	sortedVals := make([]types.Governor, len(governors))
	copy(sortedVals, governors)

	suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals = app, ctx, queryClient, addrs, governors
}

func TestParams(t *testing.T) {
	app := utils.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	expParams := types.DefaultParams()

	// check that the empty keeper loads the default
	resParams := app.StakingKeeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))

	// modify a params, save, and retrieve
	expParams.MaxValidators = 777
	app.StakingKeeper.SetParams(ctx, expParams)
	resParams = app.StakingKeeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
