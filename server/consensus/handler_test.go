package consensus

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types3 "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func TestMapAdmissionHandler(t *testing.T) {
	allowedMessages := []string{
		proto.MessageName(&types.MsgUpdateRewardAddress{}),
		proto.MessageName(&types.MsgUpdateWhitelistedRelayers{}),
	}

	handler := AllowedMessagesHandler(allowedMessages)

	tests := []struct {
		name    string
		message sdk.Msg
		wantErr bool
	}{
		{
			name:    "Allowed message 1",
			message: &types.MsgUpdateRewardAddress{},
			wantErr: false,
		},
		{
			name:    "Allowed message 2",
			message: &types.MsgUpdateWhitelistedRelayers{},
			wantErr: false,
		},
		{
			name:    "Not allowed message",
			message: &types3.MsgSend{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler(sdk.Context{}, tt.message)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "is not allowed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
