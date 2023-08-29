package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Implements DelegationSet interface
var _ types.DelegationSet = Keeper{}

// keeper of the staking store
type Keeper struct {
	stakingkeeper.Keeper
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, ak types.AccountKeeper, bk types.BankKeeper,
	ps paramtypes.Subspace,
) Keeper {
	k := stakingkeeper.NewKeeper(cdc, key, ak, bk, ps)
	return Keeper{
		Keeper: k,
	}
}

// Override this function, which is called by genutil when the genesis state is created
// We don't want to return the validator set
func (k Keeper) ApplyAndReturnValidatorSetUpdates(ctx sdk.Context) (updates []abci.ValidatorUpdate, err error) {
	_, err = k.Keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	return updates, err
}

// Set the validator hooks
func (k *Keeper) SetHooks(sh types.StakingHooks) *Keeper {
	k.Keeper = *k.Keeper.SetHooks(sh)
	return k
}
