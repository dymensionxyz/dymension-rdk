package types_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type params struct {
	valAddr     sdk.ValAddress
	pubKey      cryptotypes.PubKey
	description stakingtypes.Description
}

var validDescription = stakingtypes.Description{
	Moniker:         "test-moniker",
	Identity:        "test-identity",
	Website:         "test-website",
	SecurityContact: "test-security",
	Details:         "test-details",
}

func TestMsgCreateSequencer_ValidateBasic(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.ValAddress(pubkey.Address())

	tests := []struct {
		name string
		msg  params
		err  error
	}{
		{
			name: "valid message",
			msg: params{
				valAddr:     addr,
				pubKey:      pubkey,
				description: validDescription,
			},
			err: nil,
		}, {
			name: "empty address",
			msg: params{
				valAddr:     sdk.ValAddress([]byte("")),
				pubKey:      pubkey,
				description: validDescription,
			},
			err: types.ErrEmptyDelegatorAddr,
		}, {
			name: "missing pubkey",
			msg: params{
				valAddr:     addr,
				pubKey:      nil,
				description: validDescription,
			},
			err: types.ErrEmptyValidatorPubKey,
		}, {
			name: "missing moniker",
			msg: params{
				valAddr: addr,
				pubKey:  pubkey,
				description: stakingtypes.Description{
					Moniker:         "",
					Identity:        "test-identity",
					Website:         "test-website",
					SecurityContact: "test-security",
					Details:         "test-details",
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}

	for _, tt := range tests {
		msg, err := types.NewMsgCreateSequencer(tt.msg.valAddr, tt.msg.pubKey, tt.msg.description)
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, err)
			err := msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
