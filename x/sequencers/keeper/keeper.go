package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
)

// Implements ValidatorSet interface
// stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
// var _ stakingtypes.ValidatorSet = Keeper{}

// Implements DelegationSet interface
// var _ stakingtypes.DelegationSet = Keeper{}

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramstore paramtypes.Subspace
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
