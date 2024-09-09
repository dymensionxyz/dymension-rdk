package timeupgrade

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	prototypes "github.com/gogo/protobuf/types"

	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper, upgradeKeeper upgradekeeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	upgradeTimeTimestamp, err := getUpgradeTime(ctx, k)
	if err != nil {
		err = cleanTimeUpgrade(ctx, k)
		if err != nil {
			panic(fmt.Errorf("failed to clean time upgrade: %w", err))
		}
		return
	}

	if ctx.BlockTime().After(upgradeTimeTimestamp) {
		err = setPlanToNextBlock(ctx, k, upgradeKeeper)
		if err != nil {
			err = cleanTimeUpgrade(ctx, k)
			if err != nil {
				panic(fmt.Errorf("failed to clean time upgrade: %w", err))
			}
			return
		}
	}
}

// setPlanToNextBlock sets the upgrade plan to the next block and schedules the upgrade
func setPlanToNextBlock(ctx sdk.Context, k keeper.Keeper, upgradeKeeper upgradekeeper.Keeper) error {
	plan, err := k.UpgradePlan.Get(ctx)
	if err != nil {
		return err
	}

	plan.Height = ctx.BlockHeight() + 1
	err = upgradeKeeper.ScheduleUpgrade(ctx, plan)
	if err != nil {
		return err
	}

	err = cleanTimeUpgrade(ctx, k)
	if err != nil {
		return err
	}

	return nil
}

// cleanTimeUpgrade removes the upgrade time and plan from the store
func cleanTimeUpgrade(ctx sdk.Context, k keeper.Keeper) error {
	err := k.UpgradeTime.Remove(ctx)
	if err != nil {
		return err
	}

	err = k.UpgradePlan.Remove(ctx)
	if err != nil {
		return err
	}
	return nil
}

// getUpgradeTime gets the upgrade time from the store
func getUpgradeTime(ctx sdk.Context, k keeper.Keeper) (time.Time, error) {
	upgradeTime, err := k.UpgradeTime.Get(ctx)
	if err != nil {
		return time.Time{}, err
	}

	upgradeTimeTimestamp, err := prototypes.TimestampFromProto(&upgradeTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse upgrade time: %w", err)
	}

	return upgradeTimeTimestamp, nil
}
