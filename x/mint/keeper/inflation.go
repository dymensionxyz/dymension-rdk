package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleInflationChange(ctx sdk.Context) (inflationRate sdk.Dec, err error) {
	params := k.GetParams(ctx)
	minter := k.GetMinter(ctx)

	newInfaltion := minter.CurrentInflationRate.Mul(params.InflationRateChange)
	if newInfaltion.GT(params.TargetInflationRate) {
		newInfaltion = params.TargetInflationRate
	}
	minter.CurrentInflationRate = newInfaltion

	k.SetMinter(ctx, minter)

	return newInfaltion, nil
}
