package types_test

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/stretchr/testify/require"

	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func TestConsensusMsgUpsertSequencer(t *testing.T) {
	valAddr := utils.AccAddress()
	rewardAddr := utils.AccAddress()
	anyPubKey, err := codectypes.NewAnyWithValue(ed25519.GenPrivKey().PubKey())
	require.NoError(t, err)

	tests := []struct {
		name          string
		input         types.ConsensusMsgUpsertSequencer
		errorIs       error
		errorContains string
	}{
		{
			name: "valid",
			input: types.ConsensusMsgUpsertSequencer{
				Operator:   valAddr.String(),
				ConsPubKey: anyPubKey,
				RewardAddr: rewardAddr.String(),
			},
			errorIs:       nil,
			errorContains: "",
		},
		{
			name: "empty cons pub key",
			input: types.ConsensusMsgUpsertSequencer{
				Operator:   valAddr.String(),
				ConsPubKey: nil,
				RewardAddr: rewardAddr.String(),
			},
			errorIs:       gerrc.ErrInvalidArgument,
			errorContains: "pub key is nil",
		},
		{
			name: "invalid operator",
			input: types.ConsensusMsgUpsertSequencer{
				Operator:   "invalid_operator",
				ConsPubKey: anyPubKey,
				RewardAddr: rewardAddr.String(),
			},
			errorIs:       gerrc.ErrInvalidArgument,
			errorContains: "get operator addr from bech32",
		},
		{
			name: "invalid reward addr",
			input: types.ConsensusMsgUpsertSequencer{
				Operator:   valAddr.String(),
				ConsPubKey: anyPubKey,
				RewardAddr: "invalid_reward_addr",
			},
			errorIs:       gerrc.ErrInvalidArgument,
			errorContains: "get reward addr from bech32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.ValidateBasic()

			expectError := tt.errorIs != nil
			switch expectError {
			case true:
				require.Error(t, err)
				require.ErrorIs(t, err, tt.errorIs)
				require.Contains(t, err.Error(), tt.errorContains)
			case false:
				require.NoError(t, err)
			}
		})
	}
}
