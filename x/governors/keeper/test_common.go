package keeper // noalias

import (
	"bytes"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

// does a certain by-power index record exist
func GovernorByPowerIndexExists(ctx sdk.Context, keeper Keeper, power []byte) bool {
	store := ctx.KVStore(keeper.storeKey)
	return store.Has(power)
}

// update governor for testing
func TestingUpdateGovernor(keeper Keeper, ctx sdk.Context, governor types.Governor, apply bool) types.Governor {
	keeper.SetGovernor(ctx, governor)

	// Remove any existing power key for governor.
	store := ctx.KVStore(keeper.storeKey)
	deleted := false

	iterator := sdk.KVStorePrefixIterator(store, types.GovernorsByPowerIndexKey)
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		valAddr := types.ParseGovernorPowerRankKey(iterator.Key())
		if bytes.Equal(valAddr, governor.GetOperator()) {
			if deleted {
				panic("found duplicate power index key")
			} else {
				deleted = true
			}

			store.Delete(iterator.Key())
		}
	}

	keeper.SetGovernorByPowerIndex(ctx, governor)

	if !apply {
		ctx, _ = ctx.CacheContext()
	}
	err := keeper.ApplyGovernorSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	governor, found := keeper.GetGovernor(ctx, governor.GetOperator())
	if !found {
		panic("governor expected but not found")
	}

	return governor
}

// RandomGovernor returns a random governor given access to the keeper and ctx
func RandomGovernor(r *rand.Rand, keeper Keeper, ctx sdk.Context) (val types.Governor, ok bool) {
	vals := keeper.GetAllGovernors(ctx)
	if len(vals) == 0 {
		return types.Governor{}, false
	}

	i := r.Intn(len(vals))

	return vals[i], true
}
