package types_test

import (
	"strings"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dymensionxyz/rollapp/testutil/utils"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateSequencer_ValidateBasic(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address()).String()
	pkAny, err := codectypes.NewAnyWithValue(pubkey)
	require.NoError(t, err)

	tests := []struct {
		name string
		msg  types.MsgCreateSequencer
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgCreateSequencer{
				DelegatorAddress: "invalid_address",
				SequencerAddress: addr,
				Pubkey:           pkAny,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgCreateSequencer{
				DelegatorAddress: utils.AccAddress(),
				SequencerAddress: addr,
				Pubkey:           pkAny,
			},
		}, {
			name: "invalid sequencer address",
			msg: types.MsgCreateSequencer{
				DelegatorAddress: utils.AccAddress(),
				SequencerAddress: "invalid_address",
				Pubkey:           pkAny,
			},
			err: ErrInvalidSequencerAddress,
		}, {
			name: "invalid pubkey",
			msg: types.MsgCreateSequencer{
				DelegatorAddress: utils.AccAddress(),
				SequencerAddress: utils.AccAddress(),
				Pubkey:           pkAny,
			},
			err: sdkerrors.ErrInvalidPubKey,
		}, {
			name: "valid description",
			msg: types.MsgCreateSequencer{
				DelegatorAddress: utils.AccAddress(),
				SequencerAddress: addr,
				Pubkey:           pkAny,
				Description: Description{
					Moniker:         strings.Repeat("a", MaxMonikerLength),
					Identity:        strings.Repeat("a", MaxIdentityLength),
					Website:         strings.Repeat("a", MaxWebsiteLength),
					SecurityContact: strings.Repeat("a", MaxSecurityContactLength),
					Details:         strings.Repeat("a", MaxDetailsLength)},
			},
		}, {
			name: "invalid moniker length",
			msg: types.MsgCreateSequencer{
				DelegatorAddress: utils.AccAddress(),
				SequencerAddress: addr,
				Pubkey:           pkAny,
				Description: Description{
					Moniker: strings.Repeat("a", MaxMonikerLength+1)},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
