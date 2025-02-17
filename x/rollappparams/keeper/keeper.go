package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"
)

// Keeper of this module maintains distributing tokens to all stakers.
type Keeper struct {
	paramSpace paramtypes.Subspace
}

// NewKeeper creates new instances of the Keeper
func NewKeeper(
	paramSpace paramtypes.Subspace,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		paramSpace: paramSpace,
	}
}

// GetParams returns the total set of rollapp parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetVersion(ctx sdk.Context, version uint32) error {
	if err := types.ValidateVersion(version); err != nil {
		return err
	}
	k.paramSpace.Set(ctx, types.KeyVersion, version)
	return nil
}

func (k Keeper) SetDA(ctx sdk.Context, da string) error {
	if err := types.ValidateDa(da); err != nil {
		return err
	}
	k.paramSpace.Set(ctx, types.KeyDa, da)
	return nil
}

func (k Keeper) SetMinGasPrices(ctx sdk.Context, minGasPrices sdk.DecCoins) error {
	if err := types.ValidateMinGasPrices(minGasPrices); err != nil {
		return err
	}
	k.paramSpace.Set(ctx, types.KeyMinGasPrices, minGasPrices)
	return nil
}

func (k Keeper) DA(ctx sdk.Context) (res string) {
	k.paramSpace.Get(ctx, types.KeyDa, &res)
	return
}

func (k Keeper) Version(ctx sdk.Context) (res uint32) {
	k.paramSpace.Get(ctx, types.KeyVersion, &res)
	return
}

func (k Keeper) MinGasPrices(ctx sdk.Context) (res sdk.DecCoins) {
	k.paramSpace.Get(ctx, types.KeyMinGasPrices, &res)
	return
}

func (k Keeper) FreeIBC(ctx sdk.Context) bool {
	var res bool
	k.paramSpace.Get(ctx, types.KeyFreeIBC, &res)
	return res
}
