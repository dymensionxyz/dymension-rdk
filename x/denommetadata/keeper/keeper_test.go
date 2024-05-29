package keeper_test

import (
	"fmt"
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/testutils"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

type DenomMetadataKeeperTestSuite struct {
	suite.Suite

	app *app.App
	k   keeper.Keeper
	ctx sdk.Context
}

func TestDenomMetadataKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(DenomMetadataKeeperTestSuite))
}

func (suite *DenomMetadataKeeperTestSuite) setupTest(hooks types.DenomMetadataHooks) {
	suite.app, suite.ctx = testutils.NewTestDenommetadataKeeper(suite.T())
	suite.k = suite.app.DenommetadataKeeper
	suite.k.SetHooks(types.NewMultiDenommetadataHooks(hooks))
	// Set allowed addresses
	initialParams := types.DefaultParams()
	initialParams.AllowedAddresses = []string{senderAddress}
	suite.k.SetParams(suite.ctx, initialParams)
}

const (
	ibcBase       = "ibc/896F0081794734A2DBDF219B7575C569698F872619C43D18CC63C03CFB997257"
	senderAddress = "cosmos1s77x8wr2gzdhq8gt8c085vate0s23xu9u80wtx"
	denomTrace    = "transfer/channel-0/atom"
)

var denomMetadata = banktypes.Metadata{
	Description: "ATOM IBC",
	Base:        ibcBase,
	// NOTE: Denom units MUST be increasing
	DenomUnits: []*banktypes.DenomUnit{
		{Denom: ibcBase, Exponent: 0},
		{Denom: "ATOM", Exponent: 18},
	},
	Name:    "ATOM channel-0",
	Symbol:  "ibcATOM-0",
	Display: ibcBase,
}

func (suite *DenomMetadataKeeperTestSuite) TestCreateDenomMetadata() {
	cases := []struct {
		name             string
		metadata         []types.DenomMetadata
		hooks            *mockERC20Hook
		expectErr        string
		expectHookCalled bool
		malleate         func()
	}{
		{
			name: "success",
			metadata: []types.DenomMetadata{
				{
					TokenMetadata: denomMetadata,
					DenomTrace:    denomTrace,
				},
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: true,
		}, {
			name: "ibc denom does not match",
			metadata: []types.DenomMetadata{
				{
					TokenMetadata: denomMetadata,
					DenomTrace:    "transfer/channel-0/uatom",
				},
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: false,
			expectErr:        "denom parse from denom trace does not match metadata base denom",
		}, {
			name: "denom already exists",
			metadata: []types.DenomMetadata{
				{
					TokenMetadata: denomMetadata,
					DenomTrace:    denomTrace,
				},
			},
			malleate: func() {
				metadatas := []types.DenomMetadata{
					{
						TokenMetadata: denomMetadata,
						DenomTrace:    denomTrace,
					},
				}
				err := suite.k.CreateDenomMetadata(suite.ctx, metadatas...)
				suite.Require().NoError(err, "CreateDenomMetadata() error")
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: true,
			expectErr:        "",
		}, {
			name: "invalid denom units",
			metadata: []types.DenomMetadata{
				{
					TokenMetadata: banktypes.Metadata{
						Description: "ATOM IBC",
						Base:        ibcBase,
						DenomUnits: []*banktypes.DenomUnit{
							{Denom: ibcBase, Exponent: 18},
							{Denom: "ATOM", Exponent: 0},
						},
						Name:    "ATOM channel-0",
						Symbol:  "ibcATOM-0",
						Display: ibcBase,
					},
					DenomTrace: denomTrace,
				},
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: false,
			expectErr:        fmt.Sprintf("the exponent for base denomination unit %s must be 0", ibcBase),
		}, {
			name: "failed to create erc20 contract",
			metadata: []types.DenomMetadata{
				{
					TokenMetadata: denomMetadata,
					DenomTrace:    denomTrace,
				},
			},
			hooks: &mockERC20Hook{
				err: fmt.Errorf("failed to deploy the erc20 contract for the IBC coin"),
			},
			expectHookCalled: false,
			expectErr:        "error in after denom metadata creation hook: failed to deploy the erc20 contract for the IBC coin",
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			// reset the test state
			suite.setupTest(tc.hooks)

			if tc.malleate != nil {
				tc.malleate()
			}

			err := suite.k.CreateDenomMetadata(suite.ctx, tc.metadata...)
			if tc.expectErr != "" {
				suite.Require().ErrorContains(err, tc.expectErr)
				return
			}
			suite.Require().NoError(err, "CreateDenomMetadata() error")

			// check if the denom metadata was added
			_, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, ibcBase)
			suite.Require().True(found, "denom metadata should exist")
			suite.Require().Equal(tc.expectHookCalled, tc.hooks.createCalled, "after denom metadata creation hook should be called")
		})
	}
}

type mockERC20Hook struct {
	createCalled bool
	err          error
	sync.Mutex
}

func (m *mockERC20Hook) AfterDenomMetadataCreation(sdk.Context, banktypes.Metadata) error {
	m.Lock()
	defer m.Unlock()
	m.createCalled = m.err == nil
	return m.err
}

func (m *mockERC20Hook) AfterDenomMetadataUpdate(sdk.Context, banktypes.Metadata) error { return nil }

func TestParams(t *testing.T) {
	// Setup the test environment
	app, ctx := testutils.NewTestDenommetadataKeeper(t) // Assume you have a similar utility function for denommetadata keeper
	k := app.DenommetadataKeeper

	// Set some initial parameters
	initialParams := types.DefaultParams()
	initialParams.AllowedAddresses = []string{"cosmos19crd4fwzm9qtf5ln5l3e2vmquhevjwprk8tgxp", "cosmos1gusne8eh37myphx09hgdsy85zpl2t0kzdvu3en"} // Example addresses
	k.SetParams(ctx, initialParams)

	// Retrieve the parameters
	retrievedParams := k.GetParams(ctx)

	// Assert that the retrieved parameters match the initial ones
	require.Equal(t, initialParams, retrievedParams, "retrieved parameters should match the initial ones")

	// Test setting and getting a different set of parameters
	updatedParams := initialParams
	updatedParams.AllowedAddresses = append(updatedParams.AllowedAddresses, "cosmos1s77x8wr2gzdhq8gt8c085vate0s23xu9u80wtx")
	k.SetParams(ctx, updatedParams)
	retrievedParams = k.GetParams(ctx)
	require.Equal(t, updatedParams, retrievedParams, "retrieved parameters should match the updated ones")
}
