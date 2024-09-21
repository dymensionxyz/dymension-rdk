package timeupgrade_test

import (
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBeginBlocker_Errors(t *testing.T) {
	app := utils.Setup(t, false)
	tUpgradeKeeper, ctx := testkeepers.NewTestTimeupgradeKeeperFromApp(app)
	upgradeKeeper, _ := testkeepers.NewTestUpgradeKeeperFromApp(app)

	t.Run("UpgradeTime not found", func(t *testing.T) {
		tUpgradeKeeper.UpgradeTime.Remove(ctx)

		timeupgrade.BeginBlocker(ctx, tUpgradeKeeper, upgradeKeeper)

		_, err := tUpgradeKeeper.UpgradeTime.Get(ctx)
		require.Error(t, err, "UpgradeTime should not exist")

		// Verify that no upgrade was added to the upgradeKeeper
		plan, exists := upgradeKeeper.GetUpgradePlan(ctx)
		require.False(t, exists, "No upgrade plan should exist in upgradeKeeper")
		require.Empty(t, plan)
	})

	t.Run("UpgradeTime parsing error", func(t *testing.T) {
		invalidTime := types.Timestamp{Seconds: -1, Nanos: -1}
		err := tUpgradeKeeper.UpgradeTime.Set(ctx, invalidTime)
		require.NoError(t, err)

		timeupgrade.BeginBlocker(ctx, tUpgradeKeeper, upgradeKeeper)

		_, err = tUpgradeKeeper.UpgradeTime.Get(ctx)
		require.Error(t, err, "UpgradeTime should have been removed")

		// Verify that no upgrade was added to the upgradeKeeper
		plan, exists := upgradeKeeper.GetUpgradePlan(ctx)
		require.False(t, exists, "No upgrade plan should exist in upgradeKeeper")
		require.Empty(t, plan)
	})

	t.Run("UpgradePlan not found", func(t *testing.T) {
		pastTime := types.Timestamp{Seconds: time.Now().Add(-1 * time.Hour).Unix()}
		err := tUpgradeKeeper.UpgradeTime.Set(ctx, pastTime)
		require.NoError(t, err)
		tUpgradeKeeper.UpgradePlan.Remove(ctx)

		timeupgrade.BeginBlocker(ctx, tUpgradeKeeper, upgradeKeeper)

		_, err = tUpgradeKeeper.UpgradeTime.Get(ctx)
		require.Error(t, err, "UpgradeTime should have been removed")
		_, err = tUpgradeKeeper.UpgradePlan.Get(ctx)
		require.Error(t, err, "UpgradePlan should not exist")

		// Verify that no upgrade was added to the upgradeKeeper
		plan, exists := upgradeKeeper.GetUpgradePlan(ctx)
		require.False(t, exists, "No upgrade plan should exist in upgradeKeeper")
		require.Empty(t, plan)
	})

	t.Run("ScheduleUpgrade error", func(t *testing.T) {
		pastTime := types.Timestamp{Seconds: time.Now().Add(-1 * time.Hour).Unix()}
		err := tUpgradeKeeper.UpgradeTime.Set(ctx, pastTime)
		require.NoError(t, err)

		invalidPlan := upgradetypes.Plan{Name: "", Height: -1} // Invalid plan
		err = tUpgradeKeeper.UpgradePlan.Set(ctx, invalidPlan)
		require.NoError(t, err)

		timeupgrade.BeginBlocker(ctx, tUpgradeKeeper, upgradeKeeper)

		_, err = tUpgradeKeeper.UpgradeTime.Get(ctx)
		require.Error(t, err, "UpgradeTime should have been removed")
		_, err = tUpgradeKeeper.UpgradePlan.Get(ctx)
		require.Error(t, err, "UpgradePlan should have been removed")

		// Verify that no upgrade was added to the upgradeKeeper
		plan, exists := upgradeKeeper.GetUpgradePlan(ctx)
		require.False(t, exists, "No upgrade plan should exist in upgradeKeeper")
		require.Empty(t, plan)
	})

	t.Run("Upgrade not scheduled (future time)", func(t *testing.T) {
		futureTime := types.Timestamp{Seconds: time.Now().Add(1 * time.Hour).Unix()}
		err := tUpgradeKeeper.UpgradeTime.Set(ctx, futureTime)
		require.NoError(t, err)

		plan := upgradetypes.Plan{Name: "test", Height: 100}
		err = tUpgradeKeeper.UpgradePlan.Set(ctx, plan)
		require.NoError(t, err)

		timeupgrade.BeginBlocker(ctx, tUpgradeKeeper, upgradeKeeper)

		_, err = tUpgradeKeeper.UpgradeTime.Get(ctx)
		require.NoError(t, err, "UpgradeTime should not have been removed")
		_, err = tUpgradeKeeper.UpgradePlan.Get(ctx)
		require.NoError(t, err, "UpgradePlan should not have been removed")

		// Verify that no upgrade was added to the upgradeKeeper
		upgradePlan, exists := upgradeKeeper.GetUpgradePlan(ctx)
		require.False(t, exists, "No upgrade plan should exist in upgradeKeeper")
		require.Empty(t, upgradePlan)
	})
}

func TestBeginBlocker_HappyPath(t *testing.T) {
	app := utils.Setup(t, false)
	tUpgradeKeeper, ctx := testkeepers.NewTestTimeupgradeKeeperFromApp(app)
	upgradeKeeper, _ := testkeepers.NewTestUpgradeKeeperFromApp(app)

	// Set up a past upgrade time
	pastTime := types.Timestamp{Seconds: time.Now().Add(-1 * time.Hour).Unix()}
	err := tUpgradeKeeper.UpgradeTime.Set(ctx, pastTime)
	require.NoError(t, err)

	// Set up a valid upgrade plan
	initialPlan := upgradetypes.Plan{
		Name:   "test_upgrade",
		Height: 1000, // This height will be overwritten by BeginBlocker
		Info:   "Test upgrade",
	}
	err = tUpgradeKeeper.UpgradePlan.Set(ctx, initialPlan)
	require.NoError(t, err)

	// Set current block height
	currentHeight := int64(500)
	ctx = ctx.WithBlockHeight(currentHeight)

	// Run BeginBlocker
	timeupgrade.BeginBlocker(ctx, tUpgradeKeeper, upgradeKeeper)

	// Verify that UpgradeTime has been removed
	_, err = tUpgradeKeeper.UpgradeTime.Get(ctx)
	require.Error(t, err, "UpgradeTime should have been removed")

	// Verify that UpgradePlan has been removed
	_, err = tUpgradeKeeper.UpgradePlan.Get(ctx)
	require.Error(t, err, "UpgradePlan should have been removed")

	// Verify that the upgrade has been scheduled correctly in the upgradeKeeper
	scheduledPlan, exists := upgradeKeeper.GetUpgradePlan(ctx)
	require.True(t, exists, "An upgrade plan should exist in upgradeKeeper")
	require.Equal(t, initialPlan.Name, scheduledPlan.Name)
	require.Equal(t, currentHeight+1, scheduledPlan.Height, "The plan height should be set to the current height + 1")
	require.Equal(t, initialPlan.Info, scheduledPlan.Info)
}
