package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/dymensionxyz/dymension-rdk/x/dist/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	distkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	distkeeper.Keeper

	authKeeper    disttypes.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	seqKeeper     types.SequencerKeeper
	erc20k        types.ERC20Keeper

	feeCollectorName string
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	ak disttypes.AccountKeeper,
	bk types.BankKeeper,
	sk types.StakingKeeper,
	seqk types.SequencerKeeper,
	erc20k types.ERC20Keeper,
	feeCollectorName string,
) Keeper {
	k := distkeeper.NewKeeper(cdc, key, paramSpace, ak, bk, sk, feeCollectorName)
	return Keeper{
		Keeper:           k,
		authKeeper:       ak,
		bankKeeper:       bk,
		stakingKeeper:    sk,
		seqKeeper:        seqk,
		erc20k:           erc20k,
		feeCollectorName: feeCollectorName,
	}
}
