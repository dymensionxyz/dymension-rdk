package keeper_test

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func TestHappyPath(t *testing.T) {
	// prepare test
	var (
		app       = utils.Setup(t, false)
		_, ctx    = testkeepers.NewTestSequencerKeeperFromApp(app)
		authority = authtypes.NewModuleAddress(types.ModuleName).String()
	)

	t.Log(authority)

	// prepare test data
	var (
		operator    = utils.Proposer.GetOperator()
		rewardAddr1 = utils.AccAddress()
		rewardAddr2 = utils.AccAddress()
		relayers1   = []string{
			utils.AccAddress().String(),
			utils.AccAddress().String(),
			utils.AccAddress().String(),
		}
		relayers2 = []string{
			utils.AccAddress().String(),
			utils.AccAddress().String(),
			utils.AccAddress().String(),
		}
	)
	anyPubKey, err := codectypes.NewAnyWithValue(utils.ConsPrivKey.PubKey())
	require.NoError(t, err)

	// helper method to validate test results
	validateResults := func(rewardAddr string, relayers []string) {
		// validate sequencer
		actualSequencer, ok := app.SequencersKeeper.GetSequencer(ctx, operator)
		require.True(t, ok)
		consAddr, err := actualSequencer.GetConsAddr()
		require.NoError(t, err)
		require.Equal(t, consAddr, sdk.ConsAddress(utils.ConsPrivKey.PubKey().Address()))

		// validate reward address
		actualRewardAddr, ok := app.SequencersKeeper.GetRewardAddr(ctx, operator)
		require.True(t, ok)
		require.Equal(t, rewardAddr, actualRewardAddr.String())

		// validate relayers
		actualRelayers, err := app.SequencersKeeper.GetWhitelistedRelayers(ctx, operator)
		require.NoError(t, err)
		require.ElementsMatch(t, relayers, actualRelayers.Relayers)
	}

	t.Run("ConsensusMsgUpsertSequencer", func(t *testing.T) {
		msg := &types.ConsensusMsgUpsertSequencer{
			Signer:     authority,
			Operator:   operator.String(),
			ConsPubKey: anyPubKey,
			RewardAddr: rewardAddr1.String(),
			Relayers:   relayers1,
		}

		err = msg.ValidateBasic()
		require.NoError(t, err)

		// call msg server
		_, err = app.MsgServiceRouter().Handler(new(types.ConsensusMsgUpsertSequencer))(ctx, msg)
		require.NoError(t, err)

		// validate results
		validateResults(rewardAddr1.String(), relayers1)
	})

	t.Run("ConsensusMsgUpsertSequencer: Unauthorized", func(t *testing.T) {
		msg := &types.ConsensusMsgUpsertSequencer{
			Signer:     utils.AccAddress().String(),
			Operator:   operator.String(),
			ConsPubKey: anyPubKey,
			RewardAddr: rewardAddr1.String(),
			Relayers:   relayers1,
		}

		err = msg.ValidateBasic()
		require.NoError(t, err)

		// call msg server
		_, err = app.MsgServiceRouter().Handler(new(types.ConsensusMsgUpsertSequencer))(ctx, msg)
		require.Error(t, err)

		// previous values are unchanged
		validateResults(rewardAddr1.String(), relayers1)
	})

	t.Run("MsgUpdateRewardAddress", func(t *testing.T) {
		msg := &types.MsgUpdateRewardAddress{
			Operator:   operator.String(),
			RewardAddr: rewardAddr2.String(),
		}

		err = msg.ValidateBasic()
		require.NoError(t, err)

		// call msg server
		_, err = app.MsgServiceRouter().Handler(new(types.MsgUpdateRewardAddress))(ctx, msg)
		require.NoError(t, err)

		// validate results
		validateResults(rewardAddr2.String(), relayers1)
	})

	t.Run("MsgUpdateWhitelistedRelayers", func(t *testing.T) {
		msg := &types.MsgUpdateWhitelistedRelayers{
			Operator: operator.String(),
			Relayers: relayers2,
		}

		err = msg.ValidateBasic()
		require.NoError(t, err)

		// call msg server
		_, err = app.MsgServiceRouter().Handler(new(types.MsgUpdateWhitelistedRelayers))(ctx, msg)
		require.NoError(t, err)

		// validate results
		validateResults(rewardAddr2.String(), relayers2)
	})
}

func TestUpgradeDRS(t *testing.T) {
	// prepare test
	var (
		app    = utils.Setup(t, false)
		_, ctx = testkeepers.NewTestSequencerKeeperFromApp(app)
	)

	tests := []struct {
		name        string
		drsVersion  uint64
		expectError bool
	}{
		{
			name:        "Success: Update DRS version",
			drsVersion:  2,
			expectError: false,
		},
		{
			name:        "Success: Update to higher version",
			drsVersion:  10,
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Get initial params
			initialParams := app.RollappParamsKeeper.GetParams(ctx)

			// Create message
			msg := &types.MsgUpgradeDRS{
				Authority:  authtypes.NewModuleAddress("gov").String(),
				DrsVersion: tc.drsVersion,
			}

			// Validate basic
			err := msg.ValidateBasic()
			require.NoError(t, err)

			// Execute message
			handler := app.MsgServiceRouter().Handler(new(types.MsgUpgradeDRS))
			_, err = handler(ctx, msg)

			if tc.expectError {
				require.Error(t, err)
				// Verify params haven't changed
				currentParams := app.RollappParamsKeeper.GetParams(ctx)
				require.Equal(t, initialParams.DrsVersion, currentParams.DrsVersion)
			} else {
				require.NoError(t, err)
				// Verify params have been updated
				currentParams := app.RollappParamsKeeper.GetParams(ctx)
				require.Equal(t, uint32(tc.drsVersion), currentParams.DrsVersion)
				require.NotEqual(t, initialParams.DrsVersion, currentParams.DrsVersion)
			}
		})
	}

	t.Run("Multiple updates", func(t *testing.T) {
		versions := []uint64{3, 5, 8}

		for _, version := range versions {
			msg := &types.MsgUpgradeDRS{
				Authority:  authtypes.NewModuleAddress("gov").String(),
				DrsVersion: version,
			}

			err := msg.ValidateBasic()
			require.NoError(t, err)

			handler := app.MsgServiceRouter().Handler(new(types.MsgUpgradeDRS))
			_, err = handler(ctx, msg)
			require.NoError(t, err)

			// Verify version was updated
			currentParams := app.RollappParamsKeeper.GetParams(ctx)
			require.Equal(t, uint32(version), currentParams.DrsVersion)
		}
	})
}
