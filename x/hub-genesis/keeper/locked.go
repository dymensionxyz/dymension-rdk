package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// SetLocked sets the locked state.
func (k Keeper) SetLocked(ctx sdk.Context, locked types.Locked) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LockedKey, k.cdc.MustMarshal(&locked))
}

// GetLocked returns the locked state.
func (k Keeper) GetLocked(ctx sdk.Context) types.Locked {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LockedKey)
	if bz == nil {
		return types.Locked{}
	}
	var locked types.Locked
	k.cdc.MustUnmarshal(bz, &locked)
	return locked
}
