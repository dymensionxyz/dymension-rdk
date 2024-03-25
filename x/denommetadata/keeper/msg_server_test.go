package keeper_test

import (
	"fmt"
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

type DenomMetadataMsgServerTestSuite struct {
	suite.Suite

	app       *app.App
	k         keeper.Keeper
	msgServer types.MsgServer
	ctx       sdk.Context
}

func TestDenomMetadataMsgServerTestSuite(t *testing.T) {
	suite.Run(t, new(DenomMetadataMsgServerTestSuite))
}

func (suite *DenomMetadataMsgServerTestSuite) setupTest(hooks types.DenomMetadataHooks) {
	suite.app, suite.ctx = testutils.NewTestDenommetadataKeeper(suite.T())
	suite.k = suite.app.DenommetadataKeeper
	suite.k.SetHooks(types.NewMultiDenommetadataHooks(hooks))
	suite.msgServer = keeper.NewMsgServerImpl(suite.k)
	// Set allowed addresses
	initialParams := types.DefaultParams()
	initialParams.AllowedAddresses = []string{senderAddress}
	suite.k.SetParams(suite.ctx, initialParams)
}

const (
	ibcBase       = "ibc/7B2A4F6E798182988D77B6B884919AF617A73503FDAC27C916CD7A69A69013CF"
	senderAddress = "cosmos1s77x8wr2gzdhq8gt8c085vate0s23xu9u80wtx"
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

func (suite *DenomMetadataMsgServerTestSuite) TestCreateDenomMetadata() {
	cases := []struct {
		name             string
		msg              *types.MsgCreateDenomMetadata
		hooks            *mockERC20Hook
		expectErr        string
		expectHookCalled bool
		malleate         func()
	}{
		{
			name: "success",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
				TokenMetadata: denomMetadata,
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: true,
		}, {
			name: "permission error",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
				TokenMetadata: denomMetadata,
			},
			malleate: func() {
				initialParams := types.DefaultParams()
				initialParams.AllowedAddresses = []string{}
				suite.k.SetParams(suite.ctx, initialParams)
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: false,
			expectErr:        types.ErrNoPermission.Error(),
		}, {
			name: "denom already exists",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
				TokenMetadata: denomMetadata,
			},
			malleate: func() {
				msg := &types.MsgCreateDenomMetadata{
					SenderAddress: senderAddress,
					TokenMetadata: denomMetadata,
				}
				_, err := suite.msgServer.CreateDenomMetadata(suite.ctx, msg)
				require.NoError(suite.T(), err, "CreateDenomMetadata() error")
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: false,
			expectErr:        types.ErrDenomAlreadyExists.Error(),
		}, {
			name: "invalid denom units",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
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
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: false,
			expectErr:        fmt.Sprintf("the exponent for base denomination unit %s must be 0", ibcBase),
		}, {
			name: "failed to create erc20 contract",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
				TokenMetadata: denomMetadata,
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

			_, err := suite.msgServer.CreateDenomMetadata(suite.ctx, tc.msg)
			if tc.expectErr != "" {
				require.ErrorContains(suite.T(), err, tc.expectErr)
				return
			}
			require.NoError(suite.T(), err, "CreateDenomMetadata() error")

			// check if the denom metadata was added
			_, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, ibcBase)
			require.True(suite.T(), found, "denom metadata should exist")
			require.Equal(suite.T(), tc.expectHookCalled, tc.hooks.createCalled, "after denom metadata creation hook should be called")
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
