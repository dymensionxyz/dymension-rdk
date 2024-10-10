package keeper_test

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func TestHappyPath(t *testing.T) {
	// prepare test
	var (
		app       = utils.Setup(t, false)
		k, ctx    = testkeepers.NewTestSequencerKeeperFromApp(app)
		msgServer = keeper.NewMsgServerImpl(*k)
	)

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
			Operator:   operator.String(),
			ConsPubKey: anyPubKey,
			RewardAddr: rewardAddr1.String(),
			Relayers:   relayers1,
		}

		err = msg.ValidateBasic()
		require.NoError(t, err)

		// call msg server
		_, err = msgServer.UpsertSequencer(ctx, msg)
		require.NoError(t, err)

		// validate results
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
		_, err = msgServer.UpdateRewardAddress(ctx, msg)
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
		_, err = msgServer.UpdateWhitelistedRelayers(ctx, msg)
		require.NoError(t, err)

		// validate results
		validateResults(rewardAddr2.String(), relayers2)
	})
}
