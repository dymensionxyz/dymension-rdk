package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/dymensionxyz/dymension-rdk/x/dist/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	distkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	distkeeper.Keeper

	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	paramSpace paramtypes.Subspace

	authKeeper    types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	seqKeeper     types.SequencerKeeper

	blockedAddrs     map[string]bool
	feeCollectorName string
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper, seqk types.SequencerKeeper,
	feeCollectorName string, blockedAddrs map[string]bool,
) Keeper {

	k := distkeeper.NewKeeper(cdc, key, paramSpace, ak, bk, sk, feeCollectorName)

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		Keeper: k,

		storeKey:         key,
		cdc:              cdc,
		paramSpace:       paramSpace,
		authKeeper:       ak,
		bankKeeper:       bk,
		stakingKeeper:    sk,
		seqKeeper:        seqk,
		blockedAddrs:     blockedAddrs,
		feeCollectorName: feeCollectorName,
	}
}
