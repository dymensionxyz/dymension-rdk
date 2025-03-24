package ante_test

import (
	"context"
	"errors"
	"testing"

	"cosmossdk.io/math"
	"github.com/tendermint/tendermint/libs/log"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/server/ante"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock ERC20Keeper for testing
type mockERC20Keeper struct {
	mock.Mock
}

// TryConvertErc20Sdk implements ante.ERC20Keeper.
func (m *mockERC20Keeper) TryConvertErc20Sdk(ctx sdk.Context, sender sdk.AccAddress, receiver sdk.AccAddress, denom string, amount math.Int) error {
	args := m.Called(ctx, sender, receiver, denom, amount)
	return args.Error(0)
}

func (m *mockERC20Keeper) IsDenomRegistered(ctx sdk.Context, denom string) bool {
	args := m.Called(ctx, denom)
	return args.Bool(0)
}

func (m *mockERC20Keeper) ConvertCoin(ctx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error) {
	args := m.Called(ctx, msg)
	return args.Get(0).(*erc20types.MsgConvertCoinResponse), args.Error(1)
}

// Mock BankKeeper for testing
type mockBankKeeper struct {
	mock.Mock
}

func (m *mockBankKeeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	args := m.Called(ctx, addr)
	return args.Get(0).(sdk.Coins)
}

func (m *mockBankKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	args := m.Called(ctx, addr, denom)
	return args.Get(0).(sdk.Coin)
}

// Mock AnteHandler for testing
type mockAnteHandler struct {
	mock.Mock
}

func (m *mockAnteHandler) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
	args := m.Called(ctx, tx, simulate)
	return args.Get(0).(sdk.Context), args.Error(1)
}

func TestERC20ConversionDecorator(t *testing.T) {
	// Setup
	mockERC20 := new(mockERC20Keeper)
	mockBank := new(mockBankKeeper)
	mockNext := new(mockAnteHandler)

	decorator := ante.NewERC20ConversionDecorator(mockERC20, mockBank)

	// Create a context
	ctx := sdk.Context{}.WithLogger(log.NewNopLogger())

	accAddr := utils.AccAddress()

	// Create test cases
	testCases := []struct {
		name        string
		msgs        []sdk.Msg
		setupMocks  func()
		expectedErr bool
	}{
		{
			name: "MsgCreateValidator with registered denom",
			msgs: []sdk.Msg{
				&stakingtypes.MsgCreateValidator{
					DelegatorAddress: accAddr.String(),
					ValidatorAddress: sdk.ValAddress(accAddr).String(),
					Value:            sdk.NewCoin("registered", sdk.NewInt(100)),
				},
			},
			setupMocks: func() {
				// Setup expectations
				mockERC20.On("IsDenomRegistered", mock.Anything, "registered").Return(true)
				// Mock that the balance is not sufficient, so conversion is needed
				mockBank.On("GetBalance", mock.Anything, mock.Anything, "registered").Return(sdk.NewCoin("registered", sdk.NewInt(0)))
				mockERC20.On("ConvertCoin", mock.Anything, mock.Anything).Return(&erc20types.MsgConvertCoinResponse{}, nil)
				mockNext.On("AnteHandle", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
			},
			expectedErr: false,
		},
		{
			name: "MsgCreateValidator with unregistered denom",
			msgs: []sdk.Msg{
				&stakingtypes.MsgCreateValidator{
					DelegatorAddress: accAddr.String(),
					ValidatorAddress: sdk.ValAddress(accAddr).String(),
					Value:            sdk.NewCoin("unregistered", sdk.NewInt(100)),
				},
			},
			setupMocks: func() {
				// Setup expectations
				mockERC20.On("IsDenomRegistered", mock.Anything, "unregistered").Return(false)
				mockNext.On("AnteHandle", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
			},
			expectedErr: false,
		},
		{
			name: "MsgDelegate with registered denom",
			msgs: []sdk.Msg{
				&stakingtypes.MsgDelegate{
					DelegatorAddress: accAddr.String(),
					ValidatorAddress: sdk.ValAddress(accAddr).String(),
					Amount:           sdk.NewCoin("registered", sdk.NewInt(100)),
				},
			},
			setupMocks: func() {
				// Setup expectations
				mockERC20.On("IsDenomRegistered", mock.Anything, "registered").Return(true)
				// Mock that the balance is not sufficient, so conversion is needed
				mockBank.On("GetBalance", mock.Anything, mock.Anything, "registered").Return(sdk.NewCoin("registered", sdk.NewInt(0)))
				mockERC20.On("ConvertCoin", mock.Anything, mock.Anything).Return(&erc20types.MsgConvertCoinResponse{}, nil)
				mockNext.On("AnteHandle", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
			},
			expectedErr: false,
		},
		{
			name: "MsgDelegate with unregistered denom",
			msgs: []sdk.Msg{
				&stakingtypes.MsgDelegate{
					DelegatorAddress: accAddr.String(),
					ValidatorAddress: sdk.ValAddress(accAddr).String(),
					Amount:           sdk.NewCoin("unregistered", sdk.NewInt(100)),
				},
			},
			setupMocks: func() {
				// Setup expectations
				mockERC20.On("IsDenomRegistered", mock.Anything, "unregistered").Return(false)
				mockNext.On("AnteHandle", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
			},
			expectedErr: false,
		},
		{
			name: "MsgDelegate with conversion error",
			msgs: []sdk.Msg{
				&stakingtypes.MsgDelegate{
					DelegatorAddress: accAddr.String(),
					ValidatorAddress: sdk.ValAddress(accAddr).String(),
					Amount:           sdk.NewCoin("error", sdk.NewInt(100)),
				},
			},
			setupMocks: func() {
				// Setup expectations
				mockERC20.On("IsDenomRegistered", mock.Anything, "error").Return(true)
				// Mock that the balance is not sufficient, so conversion is needed
				mockBank.On("GetBalance", mock.Anything, mock.Anything, "error").Return(sdk.NewCoin("error", sdk.NewInt(0)))
				mockERC20.On("ConvertCoin", mock.Anything, mock.Anything).Return(&erc20types.MsgConvertCoinResponse{}, errors.New("error"))
			},
			expectedErr: true,
		},
		{
			name: "Other message type",
			msgs: []sdk.Msg{
				&stakingtypes.MsgBeginRedelegate{
					DelegatorAddress:    sdk.ValAddress("delegator").String(),
					ValidatorSrcAddress: sdk.ValAddress("validator1").String(),
					ValidatorDstAddress: sdk.ValAddress("validator2").String(),
					Amount:              sdk.NewCoin("registered", sdk.NewInt(100)),
				},
			},
			setupMocks: func() {
				mockNext.On("AnteHandle", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mocks
			mockERC20.ExpectedCalls = nil
			mockBank.ExpectedCalls = nil
			mockNext.ExpectedCalls = nil

			// Setup mocks
			tc.setupMocks()

			// Create a transaction
			tx := createTestTx(t, tc.msgs)

			// Call the decorator
			_, err := decorator.AnteHandle(ctx, tx, false, mockNext.AnteHandle)

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// Verify mock expectations
			mockERC20.AssertExpectations(t)
			mockBank.AssertExpectations(t)
			mockNext.AssertExpectations(t)
		})
	}
}

// Helper function to create a test transaction
func createTestTx(t *testing.T, msgs []sdk.Msg) sdk.Tx {
	// Create a dummy private key
	_, pk := createDummyPubKey(t)

	// Create a transaction builder
	txBuilder := simapp.MakeTestEncodingConfig().TxConfig.NewTxBuilder()

	// Set messages
	err := txBuilder.SetMsgs(msgs...)
	require.NoError(t, err)

	// Create a dummy signature
	sig := signing.SignatureV2{
		PubKey: pk,
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
			Signature: []byte("dummy_signature"),
		},
		Sequence: 0,
	}

	// Set signatures
	err = txBuilder.SetSignatures(sig)
	require.NoError(t, err)

	// Build the transaction
	tx := txBuilder.GetTx()

	return tx
}

// Helper function to create a dummy public key
func createDummyPubKey(t *testing.T) (string, cryptotypes.PubKey) {
	pk := simapp.CreateTestPubKeys(1)[0]
	return sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), pk.Address().Bytes()), pk
}
