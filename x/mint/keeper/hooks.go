package keeper

import (
	"fmt"

	epochstypes "github.com/dymensionxyz/dymension-rdk/x/epochs/types"
	"github.com/dymensionxyz/dymension-rdk/x/mint/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochInfo epochstypes.EpochInfo) {
}

// AfterEpochEnd is a hook which is executed after the end of an epoch.
// This hook should attempt to mint and distribute coins according to
// the configuration set via parameters. In addition, it handles the logic
// for reducing minted coins according to the parameters.
// For an attempt to mint to occur:
// - given epochIdentifier must be equal to the mint epoch identifier set via parameters.
// - given epochNumber must be greater than or equal to the mint start epoch set via parameters.
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochInfo epochstypes.EpochInfo) {
	epochIdentifier := epochInfo.Identifier
	epochNumber := epochInfo.CurrentEpoch
	params := k.GetParams(ctx)

	// not distribute rewards if it's not time yet for rewards distribution
	if epochIdentifier != params.EpochIdentifier || epochNumber < params.MintingRewardsDistributionStartEpoch {
		return
	}

	if epochNumber == params.MintingRewardsDistributionStartEpoch {
		k.SetLastReductionEpochNum(ctx, epochNumber)
	}

	// fetch stored minter & params
	minter := k.GetMinter(ctx)

	// Check if we have hit an epoch where we update the inflation parameter.
	// Since epochs only update based on BFT time data, it is safe to store the "reductioning period time"
	// in terms of the number of epochs that have transpired.
	if epochNumber >= params.ReductionPeriodInEpochs+k.GetLastReductionEpochNum(ctx) {
		// Reduce the reward per reduction period
		minter.EpochProvisions = minter.NextEpochProvisions(params)
		k.SetMinter(ctx, minter)
		k.SetLastReductionEpochNum(ctx, epochNumber)
	}

	// mint coins, update supply
	mintedCoin := minter.EpochProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	ctx.Logger().Info("AfterEpochEnd, minted coins", types.ModuleName, "mintedCoins", mintedCoins, "height", ctx.BlockHeight())

	// send the minted coins to the fee collector account
	err = k.DistributeMintedCoin(ctx, mintedCoin)
	if err != nil {
		panic(err)
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeEpochNumber, fmt.Sprintf("%d", epochNumber)),
			sdk.NewAttribute(types.AttributeKeyEpochProvisions, minter.EpochProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for incentives keeper.
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Return the wrapper struct.
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks.
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochInfo epochstypes.EpochInfo) {
	h.k.BeforeEpochStart(ctx, epochInfo)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochInfo epochstypes.EpochInfo) {
	h.k.AfterEpochEnd(ctx, epochInfo)
}
