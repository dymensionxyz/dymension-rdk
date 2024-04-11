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

// GenesisTriggererAllowlist returns the GenesisTriggererAllowlist param
func (k Keeper) GenesisTriggererAllowlist(ctx sdk.Context) (res []types.GenesisTriggererParams) {
	k.paramstore.Get(ctx, types.KeyGenesisTriggererAllowlist, &res)
	return
}

func (k Keeper) IsAddressInGenesisTriggererAllowList(ctx sdk.Context, address string) bool {
	Allowlist := k.GenesisTriggererAllowlist(ctx)
	for _, item := range Allowlist {
		if item.Address == address {
			return true
		}
	}
	return false
}
