package types_test

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/stretchr/testify/require"

	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func TestMsgUpdateRewardAddress(t *testing.T) {
	valAddr := sdk.ValAddress(utils.AccAddress())
	rewardAddr := utils.AccAddress()

	tests := []struct {
		name          string
		input         types.MsgUpdateRewardAddress
		errorIs       error
		errorContains string
	}{
		{
			name: "valid",
			input: types.MsgUpdateRewardAddress{
				Operator:   valAddr.String(),
				RewardAddr: rewardAddr.String(),
			},
			errorIs:       nil,
			errorContains: "",
		},
		{
			name: "invalid operator",
			input: types.MsgUpdateRewardAddress{
				Operator:   "invalid_operator",
				RewardAddr: rewardAddr.String(),
			},
			errorIs:       gerrc.ErrInvalidArgument,
			errorContains: "get operator addr from bech32",
		},
		{
			name: "invalid reward addr",
			input: types.MsgUpdateRewardAddress{
				Operator:   valAddr.String(),
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

func TestMsgUpdateWhitelistedRelayers(t *testing.T) {
	valAddr := sdk.ValAddress(utils.AccAddress())
	addr := utils.AccAddress()
	relayers := []string{
		utils.AccAddress().String(),
		utils.AccAddress().String(),
	}

	tests := []struct {
		name          string
		input         types.MsgUpdateWhitelistedRelayers
		errorIs       error
		errorContains string
	}{
		{
			name: "valid",
			input: types.MsgUpdateWhitelistedRelayers{
				Operator: valAddr.String(),
				Relayers: relayers,
			},
			errorIs:       nil,
			errorContains: "",
		},
		{
			name: "invalid relayer addr",
			input: types.MsgUpdateWhitelistedRelayers{
				Operator: valAddr.String(),
				Relayers: []string{"invalid"},
			},
			errorIs:       gerrc.ErrInvalidArgument,
			errorContains: "validate whitelisted relayers",
		},
		{
			name: "duplicated relayers",
			input: types.MsgUpdateWhitelistedRelayers{
				Operator: valAddr.String(),
				Relayers: []string{addr.String(), addr.String()},
			},
			errorIs:       gerrc.ErrInvalidArgument,
			errorContains: "validate whitelisted relayers",
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

func TestConsensusMsgUpsertSequencer(t *testing.T) {
	valAddr := utils.AccAddress()
	rewardAddr := utils.AccAddress()
	anyPubKey, err := codectypes.NewAnyWithValue(ed25519.GenPrivKey().PubKey())
	require.NoError(t, err)
	relayers := []string{
		utils.AccAddress().String(),
		utils.AccAddress().String(),
	}

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
				Relayers:   relayers,
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
				Relayers:   relayers,
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
				Relayers:   relayers,
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
				Relayers:   relayers,
			},
			errorIs:       gerrc.ErrInvalidArgument,
			errorContains: "get reward addr from bech32",
		},
		{
			name: "invalid relayer addr",
			input: types.ConsensusMsgUpsertSequencer{
				Operator:   valAddr.String(),
				ConsPubKey: anyPubKey,
				RewardAddr: rewardAddr.String(),
				Relayers:   []string{"invalid"},
			},
			errorIs:       gerrc.ErrInvalidArgument,
			errorContains: "validate whitelisted relayers",
		},
		{
			name: "duplicated relayers",
			input: types.ConsensusMsgUpsertSequencer{
				Operator:   valAddr.String(),
				ConsPubKey: anyPubKey,
				RewardAddr: rewardAddr.String(),
				Relayers:   []string{rewardAddr.String(), rewardAddr.String()},
			},
			errorIs:       gerrc.ErrInvalidArgument,
			errorContains: "validate whitelisted relayers",
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
