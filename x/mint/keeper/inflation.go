package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleInflationChange(ctx sdk.Context) (inflationRate sdk.Dec, err error) {
	params := k.GetParams(ctx)
	minter := k.GetMinter(ctx)

	if minter.CurrentInflationRate.LT(params.TargetInflationRate) {
		// Increase inflation
		newInflation := minter.CurrentInflationRate.Add(params.InflationRateChange)
		if newInflation.GT(params.TargetInflationRate) {
			newInflation = params.TargetInflationRate
		}
		minter.CurrentInflationRate = newInflation
	} else if minter.CurrentInflationRate.GT(params.TargetInflationRate) {
		// Decrease inflation
		newInflation := minter.CurrentInflationRate.Sub(params.InflationRateChange)
		if newInflation.LT(params.TargetInflationRate) {
			newInflation = params.TargetInflationRate
		}
		minter.CurrentInflationRate = newInflation
	}

	k.SetMinter(ctx, minter)

	return minter.CurrentInflationRate, nil
}
