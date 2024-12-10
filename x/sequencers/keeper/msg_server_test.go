package keeper_test

import (
	"fmt"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
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
		msgServer = keeper.NewMsgServerImpl(app.SequencersKeeper)
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

	t.Run("MsgBumpAccountSequences", func(t *testing.T) {
		// all accounts are module accounts. this ensures if the assumption
		// change in the future we will realize it.
		accs := app.AccountKeeper.GetAllAccounts(ctx)
		for _, acc := range accs {
			require.IsType(t, acc, &authtypes.ModuleAccount{})
		}

		// add a new account
		newAcc := utils.AccAddress()
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, newAcc)
		app.AccountKeeper.SetAccount(ctx, acc)

		// now we invoke bump account sequences and we should see this new acc
		// sequence bumped.
		msg := &types.MsgBumpAccountSequences{
			Authority: authority,
		}
		resp, err := msgServer.BumpAccountSequences(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// ensure accounts are correctly bumped
		accs = app.AccountKeeper.GetAllAccounts(ctx)
		for _, acc := range accs {
			switch concreteAccount := acc.(type) {
			case *authtypes.ModuleAccount:
				// module accounts should not be bumped
				require.Zero(t, concreteAccount.GetSequence())
			case *authtypes.BaseAccount:
				// base accounts should be bumped
				require.Equal(t, uint64(keeper.BumpSequence), concreteAccount.GetSequence())
			}
		}
	})
}

func TestUpgradeDRS(t *testing.T) {
	// prepare test
	var (
		app    = utils.Setup(t, false)
		_, ctx = testkeepers.NewTestSequencerKeeperFromApp(app)
	)

	tests := []struct {
		name       string
		drsVersion uint64
	}{
		{
			name:       "Same DRS version, don't upgrade",
			drsVersion: 1,
		},
		{
			name:       "Different DRS version, upgrade required",
			drsVersion: 2,
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
			require.NoError(t, err)

			plan, ok := app.UpgradeKeeper.GetUpgradePlan(ctx)
			if initialParams.DrsVersion == uint32(tc.drsVersion) {
				// Verify there is no upgrade plan created
				require.False(t, ok)
			} else {
				// Verify upgrade drs plan exists
				require.True(t, ok)
				require.Equal(t, plan.Name, fmt.Sprint("upgrade-drs-", tc.drsVersion))
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

			plan, ok := app.UpgradeKeeper.GetUpgradePlan(ctx)
			// Verify upgrade drs plan exists
			require.True(t, ok)
			require.Equal(t, plan.Name, fmt.Sprint("upgrade-drs-", version))
		}
	})
}
