package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

func (k Keeper) Allocate(ctx sdk.Context) error {
	var (
		totalStakingPower    = k.stakingKeeper.GetLastTotalPower(ctx)
		totalStakingPowerDec = sdk.NewDecFromInt(totalStakingPower)
	)

	err := k.IterateGauges(ctx, func(_ uint64, gauge types.Gauge) (stop bool, err error) {
		var (
			address      = sdk.MustAccAddressFromBech32(gauge.Address)
			gaugeRewards = sdk.NewDecCoinsFromCoins(k.bankKeeper.GetAllBalances(ctx, address)...)
		)

		switch gauge.VestingCondition.Condition.(type) {
		case *types.VestingCondition_Block:
		case *types.VestingCondition_Epoch:
		}

		switch gauge.QueryCondition.Condition.(type) {
		case *types.QueryCondition_Stakers:
			k.AllocateStakers(ctx, totalStakingPowerDec, gaugeRewards)
		case *types.QueryCondition_Funds:
		}

		return false, nil
	})
	if err != nil {
		return fmt.Errorf("iterate gauges: %w", err)
	}

	return nil
}

func (k Keeper) AllocateStakers(ctx sdk.Context, stakingPower sdk.Dec, gaugeRewards sdk.DecCoins) {
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		var (
			valPower      = validator.GetConsensusPower(sdk.DefaultPowerReduction)
			powerFraction = sdk.NewDec(valPower).QuoTruncate(stakingPower)
			reward        = gaugeRewards.MulDecTruncate(powerFraction)
		)

		k.distrKeeper.AllocateTokensToValidator(ctx, validator, reward)
		return false
	})
}
