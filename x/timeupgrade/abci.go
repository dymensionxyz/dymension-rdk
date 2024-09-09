package timeupgrade

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper, upgradeKeeper upgradekeeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	upgradeTimeTimestamp, err := k.GetUpgradeTime(ctx)
	if err != nil {
		err = k.CleanTimeUpgrade(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to clean time upgrade: %w", err))
		}
		return
	}

	if ctx.BlockTime().After(upgradeTimeTimestamp) {
		err = setPlanToNextBlock(ctx, k, upgradeKeeper)
		if err != nil {
			err = k.CleanTimeUpgrade(ctx)
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

	err = k.CleanTimeUpgrade(ctx)
	if err != nil {
		return err
	}

	return nil
}
