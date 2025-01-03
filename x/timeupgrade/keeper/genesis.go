package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, s *types.GenesisState) {
	k.SetUpgradePlan(ctx, s.Plan)
	k.SetUpgradeTime(ctx, s.Timestamp)

}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {

}
