package keeper

import (
	"fmt"
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// StakingKeeper returns the historical headers kept in store.
type StakingKeeper interface {
	GetHistoricalInfo(ctx sdk.Context, height int64) (stakingtypes.HistoricalInfo, bool)
	GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, found bool)
	UnbondingTime(ctx sdk.Context) time.Duration
}

type AuthAccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

var _ StakingKeeper = (*Keeper)(nil)

type Keeper struct {
	cdc               codec.BinaryCodec
	storeKey          storetypes.StoreKey
	paramstore        paramtypes.Subspace
	authAccountKeeper AuthAccountKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	// TODO: auth account keeper
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
