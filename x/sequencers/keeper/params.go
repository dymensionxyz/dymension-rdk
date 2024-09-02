package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// GetParams returns the total set of sequencers parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the sequencers parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	// TODO: this is hack needed to make the light client work
	//  we need to synchronise this value with the value on the Hub
	params.UnbondingTime = time.Hour * 24 * 14

	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyUnbondingTime, &res)
	return
}

// HistoricalEntries = number of historical info entries
// to persist in store
func (k Keeper) HistoricalEntries(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyHistoricalEntries, &res)
	return
}
