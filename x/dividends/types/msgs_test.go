package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	testutils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateGauge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name          string
		msg           types.MsgUpdateGauge
		expectedError error
	}{
		{
			name: "valid message",
			msg: types.MsgUpdateGauge{
				Authority: testutils.AccAddress().String(),
				ApprovedDenoms: []string{
					"validdenom1",
					"validdenom2",
					"ibc/7B2A4F6E798182988D77B6B884919AF617A73503FDAC27C916CD7A69A69013CF",
					"erc20/0xB69c34f580d74396Daeb327D35B4fb4677353Fa9",
				},
				//ApprovedDenoms: []string{"validdenom1", "validdenom2", ibcBase},
			},
			expectedError: nil,
		},
		{
			name: "invalid authority address",
			msg: types.MsgUpdateGauge{
				Authority:      "invalid_address",
				ApprovedDenoms: []string{"validdenom1", "validdenom2"},
			},
			expectedError: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid approved denom",
			msg: types.MsgUpdateGauge{
				Authority:      testutils.AccAddress().String(),
				ApprovedDenoms: []string{"invalid denom"},
			},
			expectedError: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty approved denoms",
			msg: types.MsgUpdateGauge{
				Authority:      testutils.AccAddress().String(),
				ApprovedDenoms: []string{},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
