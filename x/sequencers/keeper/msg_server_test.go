package keeper_test

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func TestUpsertHappyPath(t *testing.T) {
	// prepare test
	var (
		app       = utils.Setup(t, false)
		k, ctx    = testkeepers.NewTestSequencerKeeperFromApp(app)
		msgServer = keeper.NewMsgServerImpl(*k)
	)

	// prepare test data
	var (
		operator   = utils.Proposer.GetOperator()
		rewardAddr = utils.AccAddress()
		relayers   = []string{
			utils.AccAddress().String(),
			utils.AccAddress().String(),
			utils.AccAddress().String(),
		}
	)
	anyPubKey, err := codectypes.NewAnyWithValue(utils.ConsPrivKey.PubKey())
	require.NoError(t, err)

	msg := &types.ConsensusMsgUpsertSequencer{
		Operator:   operator.String(),
		ConsPubKey: anyPubKey,
		RewardAddr: rewardAddr.String(),
		Relayers:   relayers,
	}

	err = msg.ValidateBasic()
	require.NoError(t, err)

	// call msg server
	_, err = msgServer.UpsertSequencer(ctx, msg)
	require.NoError(t, err)

	// validate results
	actualRewardAddr, ok := app.SequencersKeeper.GetRewardAddr(ctx, operator)
	require.True(t, ok)
	require.Equal(t, msg.RewardAddr, actualRewardAddr.String())

	actualRelayers, err := app.SequencersKeeper.GetWhitelistedRelayers(ctx, operator)
	require.NoError(t, err)
	require.ElementsMatch(t, msg.Relayers, actualRelayers.Relayers)
}
