package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/dist/types"
)

// GetParams returns the total set of distribution parameters.
func (k Keeper) GetParams(clientCtx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(clientCtx, &params)
	return params
}

// SetParams sets the distribution parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetCommunityTax returns the current distribution community tax.
func (k Keeper) GetCommunityTax(ctx sdk.Context) (percent sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyCommunityTax, &percent)
	return percent
}

// GetBaseProposerReward returns the current distribution base proposer rate.
func (k Keeper) GetProposerReward(ctx sdk.Context) (percent sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyProposerReward, &percent)
	return percent
}

// GetBonusProposerReward returns the current distribution bonus proposer reward
// rate.
func (k Keeper) GetMembersReward(ctx sdk.Context) (percent sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyMembersReward, &percent)
	return percent
}

// GetWithdrawAddrEnabled returns the current distribution withdraw address
// enabled parameter.
func (k Keeper) GetWithdrawAddrEnabled(ctx sdk.Context) (enabled bool) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyWithdrawAddrEnabled, &enabled)
	return enabled
}
