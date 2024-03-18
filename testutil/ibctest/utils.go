package ibctest

import (
	"encoding/json"

	ibctypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

// ChainIDPrefix defines the default chain ID prefix for Evmos test chains
var ChainIDPrefix = "evmos_9000-"

func init() {
	ibctesting.ChainIDPrefix = ChainIDPrefix
}

func testingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	return SetupTestingApp()
}

func SetupTestingApp() (*app.App, app.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := app.MakeEncodingConfig()
	newApp := app.NewRollapp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, app.DefaultNodeHome, 5, encCdc, EmptyAppOptions{})
	return newApp, app.NewDefaultGenesisState(encCdc.Codec)
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

// IBCTestUtilSuite is a testing suite to test keeper functions.
type IBCTestUtilSuite struct {
	suite.Suite

	Coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	HubChain     *ibctesting.TestChain
	RollAppChain *ibctesting.TestChain
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *IBCTestUtilSuite) SetupTest(rollAppDenom string) {
	// set DefaultTestingAppInit function to the original one
	ibctesting.DefaultTestingAppInit = ibctesting.SetupTestingApp

	suite.Coordinator = ibctesting.NewCoordinator(suite.T(), 1) // hubChain
	hubChainId := ibctesting.GetChainID(1)
	suite.HubChain = suite.Coordinator.GetChain(hubChainId)

	// set DefaultTestingAppInit function to our function
	rollAppChainId := ibctesting.GetChainID(2)
	ibctesting.DefaultTestingAppInit = testingApp

	suite.RollAppChain = utils.SetupChain(suite.T(), suite.Coordinator, rollAppChainId, rollAppDenom)

	suite.Coordinator.Chains = map[string]*ibctesting.TestChain{
		rollAppChainId: suite.RollAppChain,
		hubChainId:     suite.HubChain,
	}
}

func (suite *IBCTestUtilSuite) NewTransferPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	path.EndpointA.ChannelConfig.Version = ibctypes.Version
	path.EndpointB.ChannelConfig.Version = ibctypes.Version

	return path
}
