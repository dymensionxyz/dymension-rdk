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

// GenesisTriggererWhitelist returns the GenesisTriggererWhitelist param
func (k Keeper) GenesisTriggererWhitelist(ctx sdk.Context) (res []types.GenesisTriggererParams) {
	k.paramstore.Get(ctx, types.KeyGenesisTriggererWhitelist, &res)
	return
}

func (k Keeper) IsAddressInGenesisTriggererWhiteList(ctx sdk.Context, address string) bool {
	whitelist := k.GenesisTriggererWhitelist(ctx)
	for _, item := range whitelist {
		if item.Address == address {
			return true
		}
	}
	return false
}
