package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
)

func (k Keeper) BeginBlock(ctx sdk.Context) error {
	return k.Allocate(ctx, types.VestingFrequency_VESTING_FREQUENCY_BLOCK)
}
