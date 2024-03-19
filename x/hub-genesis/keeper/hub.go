package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// SetHub set a specific hub in the store from its index
func (k Keeper) SetHub(ctx sdk.Context, hub types.Hub) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HubKeyPrefix))
	b := k.cdc.MustMarshal(&hub)
	store.Set(types.HubKey(
		hub.HubId,
	), b)
}

// GetHub returns a hub from its index
func (k Keeper) GetHub(
	ctx sdk.Context,
	hubId string,
) (val types.Hub, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HubKeyPrefix))

	b := store.Get(types.HubKey(
		hubId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
