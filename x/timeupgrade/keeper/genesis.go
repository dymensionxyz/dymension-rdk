package keeper

import (
	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, s *types.GenesisState) {
	if s.EmptyTimestamp() && s.EmptyPlan() {
		return
	}
	if err := k.UpgradePlan.Set(ctx, s.Plan); err != nil {
		panic(errorsmod.Wrap(err, "set upgrade plan"))
	}
	if err := k.UpgradeTime.Set(ctx, *s.Timestamp); err != nil {
		panic(errorsmod.Wrap(err, "set upgrade time"))
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	ret := types.GenesisState{}
	p, err := k.UpgradePlan.Get(ctx)
	if errorsmod.IsOf(err, collections.ErrNotFound) {
		return &ret
	}
	if err != nil {
		panic(errorsmod.Wrap(err, "upgrade plan"))
	}
	t, err := k.UpgradeTime.Get(ctx)
	if err != nil && !errorsmod.IsOf(err, collections.ErrNotFound) {
		panic(errorsmod.Wrap(err, "upgrade time"))
	}
	return &types.GenesisState{
		Plan:      p,
		Timestamp: &t,
	}
}
