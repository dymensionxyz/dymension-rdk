package keeper_test

import (
	"testing"

	"github.com/dymensionxyz/dymension-rdk/testutil/app/apptesting"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/keeper"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	querier keeper.Querier
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// SetupTest sets incentives parameters from the suite's context
func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()

	suite.querier = keeper.NewQuerier(suite.App.DividendsKeeper)
}
