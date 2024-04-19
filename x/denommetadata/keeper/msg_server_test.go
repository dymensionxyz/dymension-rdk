package keeper_test

import (
	"fmt"
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/testutils"
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
				Metadatas: []types.DenomMetadata{
					{
						TokenMetadata: denomMetadata,
						DenomTrace:    denomTrace,
					},
				},
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: true,
		}, {
			name: "permission error",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
				Metadatas: []types.DenomMetadata{
					{
						TokenMetadata: denomMetadata,
						DenomTrace:    denomTrace,
					},
				},
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
			name: "ibc denom does not match",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
				Metadatas: []types.DenomMetadata{
					{
						TokenMetadata: denomMetadata,
						DenomTrace:    "transfer/channel-0/uatom",
					},
				},
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: false,
			expectErr:        "denom parse from denom trace does not match metadata base denom",
		}, {
			name: "denom already exists",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
				Metadatas: []types.DenomMetadata{
					{
						TokenMetadata: denomMetadata,
						DenomTrace:    denomTrace,
					},
				},
			},
			malleate: func() {
				msg := &types.MsgCreateDenomMetadata{
					SenderAddress: senderAddress,
					Metadatas: []types.DenomMetadata{
						{
							TokenMetadata: denomMetadata,
							DenomTrace:    denomTrace,
						},
					},
				}
				_, err := suite.msgServer.CreateDenomMetadata(suite.ctx, msg)
				suite.Require().NoError(err, "CreateDenomMetadata() error")
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: false,
			expectErr:        types.ErrDenomAlreadyExists.Error(),
		}, {
			name: "invalid denom units",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
				Metadatas: []types.DenomMetadata{
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
			},
			hooks:            &mockERC20Hook{},
			expectHookCalled: false,
			expectErr:        fmt.Sprintf("the exponent for base denomination unit %s must be 0", ibcBase),
		}, {
			name: "failed to create erc20 contract",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
				Metadatas: []types.DenomMetadata{
					{
						TokenMetadata: denomMetadata,
						DenomTrace:    denomTrace,
					},
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

			_, err := suite.msgServer.CreateDenomMetadata(suite.ctx, tc.msg)
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
