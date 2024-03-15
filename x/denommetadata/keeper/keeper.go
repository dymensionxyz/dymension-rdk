package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	denommetadatatypes "github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

// Keeper of this module maintains distributing tokens to all stakers.
type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	bankKeeper denommetadatatypes.BankKeeper
}

// NewKeeper creates new instances of the Keeper
func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	bk denommetadatatypes.BankKeeper,
	feeCollector string,
	authority string,
) Keeper {
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		bankKeeper: bk,
	}
}
