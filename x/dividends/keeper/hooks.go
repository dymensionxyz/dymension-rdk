package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/dymensionxyz/dymension-rdk/x/epochs/types"
)

func (k Keeper) BeforeEpochStart(sdk.Context, epochstypes.EpochInfo) error {
	return nil
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochId epochstypes.EpochInfo) error {
	params := k.MustGetParams(ctx)
	if epochId.Identifier == params.DistrEpochIdentifier {
		err := k.Allocate(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

type Hooks struct {
	keeper Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Hooks returns the hook wrapper struct.
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochId epochstypes.EpochInfo) {
	err := h.keeper.BeforeEpochStart(ctx, epochId)
	if err != nil {
		h.keeper.Logger(ctx).Error("Error in BeforeEpochStart", "error", err)
	}
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochId epochstypes.EpochInfo) {
	err := h.keeper.AfterEpochEnd(ctx, epochId)
	if err != nil {
		h.keeper.Logger(ctx).Error("Error in AfterEpochEnd", "error", err)
	}
}
