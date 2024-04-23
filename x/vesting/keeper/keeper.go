package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dymensionxyz/dymension-rdk/x/vesting/types"
)

// Keeper of this module maintains distributing tokens to all stakers.
type Keeper struct {
	cdc codec.BinaryCodec
	ps  paramtypes.Subspace
}

// NewKeeper creates new instances of the vesting Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc: cdc,
		ps:  ps,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams returns the total set of denommetadata parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.ps.GetParamSet(ctx, &params)
	return
}

// SetParams sets the total set of denommetadata parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.ps.SetParamSet(ctx, &params)
}

// IsAddressPermissioned checks if the given address is permissioned to create or update denom metadata
func (k Keeper) IsAddressPermissioned(ctx sdk.Context, address string) bool {
	params := k.GetParams(ctx)
	for _, PermissionedAddress := range params.AllowedAddresses {
		if PermissionedAddress == address {
			return true
		}
	}
	return false
}
