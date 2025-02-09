package keeper_test

import (
	"testing"

	"github.com/dymensionxyz/dymension-rdk/testutil/app/apptesting"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	msgServer   keeper.MsgServer
	queryServer keeper.QueryServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.Setup()

	s.msgServer = keeper.NewMsgServer(s.App.DividendsKeeper)
	s.queryServer = keeper.NewQueryServer(s.App.DividendsKeeper)
}

func (s *KeeperTestSuite) CreateGauge(msg types.MsgCreateGauge) {
	handler := s.App.MsgServiceRouter().Handler(&types.MsgCreateGauge{})
	_, err := handler(s.Ctx, &msg)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) GetGauges() []types.Gauge {
	resp, err := s.queryServer.Gauges(s.Ctx, &types.GaugesRequest{})
	s.Require().NoError(err)
	return resp.GetData()
}

func (s *KeeperTestSuite) GetGauge(id uint64) types.Gauge {
	resp, err := s.queryServer.GaugeByID(s.Ctx, &types.GaugeByIDRequest{Id: id})
	s.Require().NoError(err)
	return resp.GetGauge()
}
