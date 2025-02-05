package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

func (k Keeper) Allocate(ctx sdk.Context, t types.VestingFrequency) error {
	var (
		totalStakingPower    = k.stakingKeeper.GetLastTotalPower(ctx)
		totalStakingPowerDec = sdk.NewDecFromInt(totalStakingPower)
	)

	err := k.IterateGauges(ctx, func(_ uint64, gauge types.Gauge) (stop bool, err error) {
		// Check if it's time to allocate rewards for this gauge
		if gauge.VestingFrequency != t {
			return false, nil
		}

		var (
			address      = sdk.MustAccAddressFromBech32(gauge.Address)
			gaugeBalance = sdk.NewDecCoinsFromCoins(k.bankKeeper.GetAllBalances(ctx, address)...)
			gaugeRewards sdk.DecCoins
		)

		switch c := gauge.VestingCondition.Condition.(type) {
		case *types.VestingCondition_Limited:
			// Estimate how to evenly distribute rewards through epochs/blocks
			if c.Limited.NumUnits >= c.Limited.FilledUnits {
				// TODO: remove this gauge, there's nothing to fill anymore
				return false, nil
			}

			remainingUnits := c.Limited.NumUnits - c.Limited.FilledUnits
			gaugeRewards = gaugeBalance.QuoDec(sdk.NewDec(remainingUnits))

		case *types.VestingCondition_Perpetual:
			gaugeRewards = gaugeBalance
		}

		switch gauge.QueryCondition.Condition.(type) {
		case *types.QueryCondition_Stakers:
			k.AllocateStakers(ctx, gaugeRewards, totalStakingPowerDec)
		}

		return false, nil
	})
	if err != nil {
		return fmt.Errorf("iterate gauges: %w", err)
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
