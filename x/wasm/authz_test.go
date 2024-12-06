package wasm_test

import (
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/stretchr/testify/require"

	"github.com/dymensionxyz/dymension-rdk/x/wasm"
)

func TestContractExecutionAuthorization_Accept(t *testing.T) {
	tests := []struct {
		name          string
		auth          *wasm.ContractExecutionAuthorization
		msg           sdk.Msg
		expectedResp  authztypes.AcceptResponse
		expectedError string
	}{
		{
			name:          "invalid message type",
			auth:          wasm.NewContractExecutionAuthorization(nil, sdk.Coins{}),
			msg:           &types.MsgStoreCode{},
			expectedResp:  authztypes.AcceptResponse{},
			expectedError: "invalid message type",
		},
		{
			name:          "contract not authorized",
			auth:          wasm.NewContractExecutionAuthorization([]string{"contract1"}, sdk.Coins{}),
			msg:           &types.MsgExecuteContract{Contract: "contract2"},
			expectedResp:  authztypes.AcceptResponse{},
			expectedError: "contract not authorized",
		},
		{
			name:          "exceeds spend limit",
			auth:          wasm.NewContractExecutionAuthorization(nil, sdk.NewCoins(sdk.NewInt64Coin("token", 100))),
			msg:           &types.MsgExecuteContract{Funds: sdk.NewCoins(sdk.NewInt64Coin("token", 200))},
			expectedResp:  authztypes.AcceptResponse{},
			expectedError: "exceeds spend limit",
		},
		{
			name:          "spend limit updated",
			auth:          wasm.NewContractExecutionAuthorization(nil, sdk.NewCoins(sdk.NewInt64Coin("token", 100))),
			msg:           &types.MsgExecuteContract{Funds: sdk.NewCoins(sdk.NewInt64Coin("token", 50))},
			expectedResp:  authztypes.AcceptResponse{Accept: true, Updated: wasm.NewContractExecutionAuthorization(nil, sdk.NewCoins(sdk.NewInt64Coin("token", 50)))},
			expectedError: "",
		},
		{
			name:          "spend limit exhausted",
			auth:          wasm.NewContractExecutionAuthorization(nil, sdk.NewCoins(sdk.NewInt64Coin("token", 100))),
			msg:           &types.MsgExecuteContract{Funds: sdk.NewCoins(sdk.NewInt64Coin("token", 100))},
			expectedResp:  authztypes.AcceptResponse{Accept: true, Delete: true},
			expectedError: "",
		},
		{
			name:          "contracts field is empty",
			auth:          wasm.NewContractExecutionAuthorization(nil, sdk.Coins{}),
			msg:           &types.MsgExecuteContract{Contract: "contract1"},
			expectedResp:  authztypes.AcceptResponse{Accept: true, Updated: wasm.NewContractExecutionAuthorization(nil, sdk.Coins{})},
			expectedError: "",
		},
		{
			name:          "spend limit field is empty",
			auth:          wasm.NewContractExecutionAuthorization([]string{"contract1"}, sdk.Coins{}),
			msg:           &types.MsgExecuteContract{Contract: "contract1", Funds: sdk.NewCoins(sdk.NewInt64Coin("token", 100))},
			expectedResp:  authztypes.AcceptResponse{Accept: true, Updated: wasm.NewContractExecutionAuthorization([]string{"contract1"}, sdk.Coins{})},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.auth.Accept(sdk.Context{}, tt.msg)
			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedResp, resp)
		})
	}
}
