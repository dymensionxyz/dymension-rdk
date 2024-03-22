package keeper_test

import (
	"context"
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

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/testutils"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

const ibcBase = "ibc/7B2A4F6E798182988D77B6B884919AF617A73503FDAC27C916CD7A69A69013CF"

func TestCreateDenomMetadata(t *testing.T) {
	tapp, ctx := setupAppWithMockEVMKeeper(t)
	k := tapp.DenommetadataKeeper

	// Prepare the test message for creating denom metadata
	createMsg := &types.MsgCreateDenomMetadata{
		SenderAddress: "cosmos1s77x8wr2gzdhq8gt8c085vate0s23xu9u80wtx",
		TokenMetadata: banktypes.Metadata{
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
		},
	}

	// Test permission error
	_, err := k.CreateDenomMetadata(sdk.WrapSDKContext(ctx), createMsg)
	require.ErrorIs(t, err, types.ErrNoPermission, "should return permission error")

	// Set allowed addresses
	initialParams := types.DefaultParams()
	initialParams.AllowedAddresses = []string{"cosmos1s77x8wr2gzdhq8gt8c085vate0s23xu9u80wtx", "cosmos1gusne8eh37myphx09hgdsy85zpl2t0kzdvu3en"}
	k.SetParams(ctx, initialParams)

	// Test creating denom metadata successfully
	_, err = k.CreateDenomMetadata(sdk.WrapSDKContext(ctx), createMsg)
	require.NoError(t, err, "creating denom metadata with allowed address should not error")

	// Test creating duplicate denom metadata
	_, err = k.CreateDenomMetadata(sdk.WrapSDKContext(ctx), createMsg)
	require.ErrorIs(t, err, types.ErrDenomAlreadyExists, "creating duplicate denom metadata should fail")

	// check if the erc20 contract was created
	token, found := tapp.BankKeeper.GetDenomMetaData(ctx, ibcBase)
	require.True(t, found, "denom metadata should exist")

	pairID := tapp.Erc20Keeper.GetTokenPairID(ctx, token.Base)
	require.True(t, len(pairID) > 0, "token pair id should exist")

	pair, found := tapp.Erc20Keeper.GetTokenPair(ctx, pairID)
	require.True(t, found, "token pair should exist")

	require.True(t, common.IsHexAddress(pair.Erc20Address), "erc20 address should be a valid hex address")
	address := common.HexToAddress(pair.Erc20Address)

	isERC20Registered := tapp.Erc20Keeper.IsERC20Registered(ctx, address)
	require.True(t, isERC20Registered, "erc20 contract should be registered")
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
