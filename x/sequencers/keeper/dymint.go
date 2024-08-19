package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// MustSetDymintValidatorUpdates  - ABCI expects the result of init genesis to return the same value as passed in InitChainer,
// so we save it to return later.
func (k Keeper) MustSetDymintValidatorUpdates(ctx sdk.Context, updates []abci.ValidatorUpdate) {
	if len(updates) != 1 {
		panic(errorsmod.Wrapf(gerrc.ErrOutOfRange, "expect 1 abci validator update: got: %d", len(updates)))
	}
	u := updates[0]
	k.cdc.MustMarshal(&u)
	ctx.KVStore(k.storeKey).Set(types.ValidatorUpdateKey, k.cdc.MustMarshal(&u))
}
