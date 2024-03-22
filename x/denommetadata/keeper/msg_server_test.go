package keeper_test

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	erc20keeper "github.com/evmos/evmos/v12/x/erc20/keeper"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
	"github.com/evmos/evmos/v12/x/evm/statedb"
	evmtypes "github.com/evmos/evmos/v12/x/evm/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/testutils"
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

func (suite *DenomMetadataMsgServerTestSuite) setupTest() {
	suite.app, suite.ctx = setupAppWithMockEVMKeeper(suite.T())
	suite.k = suite.app.DenommetadataKeeper
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
	suite.setupTest()

	cases := []struct {
		name      string
		msg       *types.MsgCreateDenomMetadata
		expectErr string
		malleate  func()
	}{
		{
			name: "success",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
				TokenMetadata: denomMetadata,
			},
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
			expectErr: types.ErrNoPermission.Error(),
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
			expectErr: types.ErrDenomAlreadyExists.Error(),
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
			expectErr: fmt.Sprintf("the exponent for base denomination unit %s must be 0", ibcBase),
		}, {
			name: "failed to create erc20 contract",
			msg: &types.MsgCreateDenomMetadata{
				SenderAddress: senderAddress,
				TokenMetadata: denomMetadata,
			},
			malleate: func() {
				// disable erc20
				params := suite.app.Erc20Keeper.GetParams(suite.ctx)
				params.EnableErc20 = false
				err := suite.app.Erc20Keeper.SetParams(suite.ctx, params)
				require.NoError(suite.T(), err, "SetParams() error")
			},
			expectErr: "error in after denom metadata creation hook: failed to deploy the erc20 contract for the IBC coin",
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			// reset the test state
			suite.setupTest()

			if tc.malleate != nil {
				tc.malleate()
			}

			_, err := suite.msgServer.CreateDenomMetadata(suite.ctx, tc.msg)
			if tc.expectErr != "" {
				require.ErrorContains(suite.T(), err, tc.expectErr)
				return
			}
			require.NoError(suite.T(), err, "CreateDenomMetadata() error")

			// check if the erc20 contract was created
			token, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, ibcBase)
			require.True(suite.T(), found, "denom metadata should exist")

			pairID := suite.app.Erc20Keeper.GetTokenPairID(suite.ctx, token.Base)
			require.True(suite.T(), len(pairID) > 0, "token pair id should exist")

			pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, pairID)
			require.True(suite.T(), found, "token pair should exist")

			require.True(suite.T(), common.IsHexAddress(pair.Erc20Address), "erc20 address should be a valid hex address")
			address := common.HexToAddress(pair.Erc20Address)

			isERC20Registered := suite.app.Erc20Keeper.IsERC20Registered(suite.ctx, address)
			require.True(suite.T(), isERC20Registered, "erc20 contract should be registered")
		})
	}
}

func setupAppWithMockEVMKeeper(t *testing.T) (*app.App, sdk.Context) {
	t.Helper()

	tapp, ctx := testutils.NewTestDenommetadataKeeper(t)

	// sneak in a mock EVM keeper
	evmKeeper := &mockEVMKeeper{}

	tapp.Erc20Keeper = erc20keeper.NewKeeper(
		tapp.GetKey(erc20types.StoreKey), tapp.AppCodec(), authtypes.NewModuleAddress(govtypes.ModuleName),
		tapp.AccountKeeper, tapp.BankKeeper, evmKeeper, tapp.StakingKeeper,
	)

	tapp.DenommetadataKeeper = keeper.NewKeeper(
		tapp.AppCodec(),
		tapp.GetKey(types.StoreKey),
		tapp.BankKeeper,
		types.NewMultiDenommetadataHooks(
			erc20keeper.NewERC20ContractRegistrationHook(tapp.Erc20Keeper),
		),
		tapp.GetSubspace(types.ModuleName),
	)
	return tapp, ctx
}

type mockEVMKeeper struct{}

func (m *mockEVMKeeper) GetParams(sdk.Context) evmtypes.Params { return evmtypes.DefaultParams() }

func (m *mockEVMKeeper) GetAccountWithoutBalance(sdk.Context, common.Address) *statedb.Account {
	return &statedb.Account{Nonce: 0, CodeHash: []byte{}}
}

func (m *mockEVMKeeper) EstimateGas(context.Context, *evmtypes.EthCallRequest) (*evmtypes.EstimateGasResponse, error) {
	return &evmtypes.EstimateGasResponse{Gas: 1000}, nil
}

func (m *mockEVMKeeper) ApplyMessage(sdk.Context, core.Message, vm.EVMLogger, bool) (*evmtypes.MsgEthereumTxResponse, error) {
	return &evmtypes.MsgEthereumTxResponse{}, nil
}
