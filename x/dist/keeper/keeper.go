package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/rollapp/x/dist/types"

	distkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	distkeeper.Keeper

	authKeeper    disttypes.AccountKeeper
	bankKeeper    disttypes.BankKeeper
	stakingKeeper disttypes.StakingKeeper
	seqKeeper     types.SequencerKeeper

	feeCollectorName string
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	ak disttypes.AccountKeeper, bk disttypes.BankKeeper, sk disttypes.StakingKeeper, seqk types.SequencerKeeper,
	feeCollectorName string, blockedAddrs map[string]bool,
) Keeper {
	k := distkeeper.NewKeeper(cdc, key, paramSpace, ak, bk, sk, feeCollectorName, blockedAddrs)
	return Keeper{
		Keeper:           k,
		authKeeper:       ak,
		bankKeeper:       bk,
		stakingKeeper:    sk,
		seqKeeper:        seqk,
		feeCollectorName: feeCollectorName,
	}
}
