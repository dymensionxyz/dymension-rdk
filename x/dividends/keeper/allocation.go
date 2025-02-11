package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

// Allocate rewards from active gauges. This function is called every block and
// every epoch. `t` indicates whether the allocation called for blocks or epochs.
func (k Keeper) Allocate(ctx sdk.Context, t types.VestingFrequency) error {
	params := k.MustGetParams(ctx)

	var (
		totalStakingPower    = k.stakingKeeper.GetLastTotalPower(ctx)
		totalStakingPowerDec = sdk.NewDecFromInt(totalStakingPower)
		gaugesToDeactivate   []uint64
	)

	err := k.IterateActiveGauges(ctx, func(gauge types.Gauge) (stop bool, err error) {
		// Check if it's time to allocate rewards for this gauge
		if gauge.VestingFrequency != t {
			return false, nil
		}

		var (
			gaugeAddress           = sdk.MustAccAddressFromBech32(gauge.Address)
			gaugeUnfilteredBalance = k.bankKeeper.GetAllBalances(ctx, gaugeAddress)
			gaugeBalance           = types.FilterDenoms(gaugeUnfilteredBalance, params.ApprovedDenoms)
			gaugeRewards           sdk.Coins
		)

		switch c := gauge.VestingCondition.Condition.(type) {
		case *types.VestingCondition_Limited:
			// Estimate how to evenly distribute rewards through epochs/blocks
			if c.Limited.NumUnits <= c.Limited.FilledUnits {
				gaugesToDeactivate = append(gaugesToDeactivate, gauge.Id)
				return false, nil
			}

			remainingUnits := c.Limited.NumUnits - c.Limited.FilledUnits
			gaugeRewards = gaugeBalance.QuoInt(math.NewInt(remainingUnits))
			c.Limited.FilledUnits += 1

		case *types.VestingCondition_Perpetual:
			gaugeRewards = gaugeBalance
		}

		switch gauge.QueryCondition.Condition.(type) {
		case *types.QueryCondition_Stakers:
			// Add rewards to validators
			gaugeRewardsDec := sdk.NewDecCoinsFromCoins(gaugeRewards...)
			k.AllocateStakers(ctx, gaugeRewardsDec, totalStakingPowerDec)

			// Adjust the balance of the gauge
			err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, gaugeAddress, distrtypes.ModuleName, gaugeRewards)
			if err != nil {
				return true, fmt.Errorf("send coins from gauge to x/distribution: %w", err)
			}
		}

		// Save the updated gauge back
		err = k.SetGauge(ctx, gauge)
		if err != nil {
			return true, fmt.Errorf("set gauge: %w", err)
		}

		// Fund the community pool with unapproved coins
		unapprovedFunds := gaugeUnfilteredBalance.Sub(gaugeBalance...)
		if !unapprovedFunds.IsZero() {
			err = k.distrKeeper.FundCommunityPool(ctx, unapprovedFunds, gaugeAddress)
			if err != nil {
				return true, fmt.Errorf("fund community pool: %w", err)
			}
		}

		return false, nil
	})
	if err != nil {
		return fmt.Errorf("iterate gauges: %w", err)
	}

	// Deactivate gauges that have been filled
	for _, id := range gaugesToDeactivate {
		err = k.DeactivateGauge(ctx, id)
		if err != nil {
			return fmt.Errorf("deactivate gauge: %w", err)
		}
	}

	return nil
}

func (k Keeper) AllocateStakers(ctx sdk.Context, gaugeRewards sdk.DecCoins, totalStakingPower sdk.Dec) {
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		var (
			valPower      = validator.GetConsensusPower(sdk.DefaultPowerReduction)
			powerFraction = sdk.NewDec(valPower).QuoTruncate(totalStakingPower)
			reward        = gaugeRewards.MulDecTruncate(powerFraction)
		)

		// TODO: send rewards to x/distribution
		k.distrKeeper.AllocateTokensToValidator(ctx, validator, reward)
		return false
	})
}
