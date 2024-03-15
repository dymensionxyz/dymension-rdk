package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GenesisTriggerrerWhitelist returns the GenesisTriggerrerWhitelist param
func (k Keeper) GenesisTriggerrerWhitelist(ctx sdk.Context) (res []types.GenesisTriggerrerParams) {
	k.paramstore.Get(ctx, types.KeyGenesisTriggerrerWhitelist, &res)
	return
}

func (k Keeper) IsAddressInGenesisTriggerrerWhiteList(ctx sdk.Context, address string) bool {
	whitelist := k.GenesisTriggerrerWhitelist(ctx)
	for _, item := range whitelist {
		if item.Address == address {
			return true
		}
	}
	return false
}
