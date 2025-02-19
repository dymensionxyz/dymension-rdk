package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

type GetGaugeBalanceFunc func(ctx sdk.Context, address sdk.AccAddress, denoms []string) sdk.Coins

// Allocate rewards from active gauges. This function is called every block and
// every epoch. `t` indicates whether the allocation called for blocks or epochs.
func (k Keeper) Allocate(ctx sdk.Context, frequency types.VestingFrequency) error {
	var (
		totalStakingPower    = k.stakingKeeper.GetLastTotalPower(ctx)
		totalStakingPowerDec = sdk.NewDecFromInt(totalStakingPower)
		gaugesToDeactivate   []uint64
	)

	err := k.IterateActiveGauges(ctx, func(gauge types.Gauge) (stop bool, err error) {
		// Check if it's time to allocate rewards for this gauge
		if gauge.VestingFrequency != frequency {
			return false, nil
		}

		var (
			gaugeAddress = gauge.GetAccAddress()
			gaugeBalance = k.getBalanceFn(ctx, gaugeAddress, gauge.ApprovedDenoms)
			gaugeRewards sdk.Coins
		)

		switch c := gauge.VestingDuration.Duration.(type) {
		case *types.VestingDuration_FixedTerm:
			// Estimate how to evenly distribute rewards through epochs/blocks
			if c.FixedTerm.NumTotal <= c.FixedTerm.NumDone {
				gaugesToDeactivate = append(gaugesToDeactivate, gauge.Id)
				return false, nil
			}

			remainingUnits := c.FixedTerm.NumTotal - c.FixedTerm.NumDone
			gaugeRewards = gaugeBalance.QuoInt(math.NewInt(remainingUnits))
			c.FixedTerm.NumDone += 1

		case *types.VestingDuration_Perpetual:
			gaugeRewards = gaugeBalance
		}

		// Gauge rewards might be zero if the gauge has no balance or the gauge
		// balance is so small that it's rounded down to zero after integer division
		if gaugeRewards.IsZero() {
			return false, nil
		}

		switch gauge.QueryCondition.Condition.(type) {
		case *types.QueryCondition_Stakers:
			// Fund the distribution module with the rewards from the gauge
			err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, gaugeAddress, distrtypes.ModuleName, gaugeRewards)
			if err != nil {
				return true, fmt.Errorf("send coins from gauge to x/distribution: %w", err)
			}

			// Add rewards to validators. AllocateStakers changes the validator's balance record,
			// but does not actually send coins to the validator's account. That's why we need to
			// send the coins to the distribution module first.
			gaugeRewardsDec := sdk.NewDecCoinsFromCoins(gaugeRewards...)
			k.AllocateStakers(ctx, gaugeRewardsDec, totalStakingPowerDec)
		}

		// Save the updated gauge back
		err = k.SetGauge(ctx, gauge)
		if err != nil {
			return true, fmt.Errorf("set gauge: %w", err)
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

		k.distrKeeper.AllocateTokensToValidator(ctx, validator, reward)
		return false
	})
}

func (k Keeper) GetBalanceFunc() GetGaugeBalanceFunc {
	return func(ctx sdk.Context, address sdk.AccAddress, denoms []string) sdk.Coins {
		var coins []sdk.Coin
		for _, denom := range denoms {
			balance := k.bankKeeper.GetBalance(ctx, address, denom)
			coins = append(coins, balance)
		}
		return sdk.NewCoins(coins...)
	}
}
