package keeper

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper provides a way to manage incentives module storage.
type Keeper struct {
	storeKey storetypes.StoreKey
}

// NewKeeper returns a new instance of the incentive module keeper struct.
func NewKeeper(
	storeKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		storeKey: storeKey,
	}
}

// Logger returns a logger instance for the incentives module.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
