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
	return types.NewParams(
		k.DA(ctx),
		k.Commit(ctx),
		k.MaxBlockGas(ctx),
		k.MaxBlockSize(ctx),
	)
}

// SetParams sets the total set of rollapp parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) MaxBlockGas(ctx sdk.Context) (res uint32) {
	k.paramSpace.Get(ctx, types.KeyBlockMaxGas, &res)
	return
}

func (k Keeper) MaxBlockSize(ctx sdk.Context) (res uint32) {
	k.paramSpace.Get(ctx, types.KeyBlockMaxSize, &res)
	return
}

func (k Keeper) DA(ctx sdk.Context) (res string) {
	k.paramSpace.Get(ctx, types.KeyDa, &res)
	return
}

func (k Keeper) Commit(ctx sdk.Context) (res string) {
	k.paramSpace.Get(ctx, types.KeyCommit, &res)
	return
}
