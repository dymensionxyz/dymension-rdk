package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// SetState sets the state state.
func (k Keeper) SetState(ctx sdk.Context, state types.State) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.StateKey, k.cdc.MustMarshal(&state))
}

// GetState returns the state state.
func (k Keeper) GetState(ctx sdk.Context) types.State {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.StateKey)
	if bz == nil {
		return types.State{}
	}
	var state types.State
	k.cdc.MustUnmarshal(bz, &state)
	return state
}
