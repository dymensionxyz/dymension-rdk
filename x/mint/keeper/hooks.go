package keeper

import (
	"fmt"

	epochstypes "github.com/dymensionxyz/dymension-rdk/x/epochs/types"
	"github.com/dymensionxyz/dymension-rdk/x/mint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeforeEpochStart is a hook which is executed before the start of an epoch. It is a no-op for mint module.
func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochInfo epochstypes.EpochInfo) {
}

// AfterEpochEnd is a hook which is executed after the end of an epoch.
// This hook should attempt to mint and distribute coins according to
// the configuration set via parameters.
// In addition, it handles the logic for updating the inflation according to the parameters.
// For a mint to occur:
// - given epochIdentifier must be equal to the mint epoch identifier set via parameters.
// - given epochNumber must be greater than or equal to the mint start epoch set via parameters.
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochInfo epochstypes.EpochInfo) {
	epochIdentifier := epochInfo.Identifier
	epochNumber := epochInfo.CurrentEpoch
	params := k.GetParams(ctx)

	// Update inflation
	if epochIdentifier == params.InflationChangeEpochIdentifier {
		newInflation, err := k.HandleInflationChange(ctx)
		if err != nil {
			k.Logger(ctx).Error("error updating inflation", "error", err)
			return
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeInflation,
				sdk.NewAttribute(types.AttributeEpochNumber, fmt.Sprintf("%d", epochNumber)),
				sdk.NewAttribute(types.AttributeKeyInflationRate, newInflation.String()),
			),
		)
	}

	// Mint coins
	if epochIdentifier == params.MintEpochIdentifier && epochNumber >= params.MintStartEpoch {
		mintedCoins, err := k.HandleMintingEpoch(ctx)
		if err != nil {
			k.Logger(ctx).Error("error minting coins", "error", err)
			return
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeMint,
				sdk.NewAttribute(types.AttributeEpochNumber, fmt.Sprintf("%d", epochNumber)),
				sdk.NewAttribute(types.AttributeKeyMintedCoins, mintedCoins.String()),
			),
		)
		return
	}
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
