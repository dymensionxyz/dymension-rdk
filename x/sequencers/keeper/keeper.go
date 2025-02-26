package keeper

import (
	"fmt"
	"time"

	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dymensionxyz/dymension-rdk/utils/collcompat"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// StakingKeeper returns the historical headers kept in store.
type StakingKeeper interface {
	GetHistoricalInfo(ctx sdk.Context, height int64) (stakingtypes.HistoricalInfo, bool)
	GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, found bool)
	UnbondingTime(ctx sdk.Context) time.Duration
}

var _ StakingKeeper = (*Keeper)(nil)

// AccountBumpFilterFunc is a function signature that filters accounts whose sequence should be bumped.
// IT is passed the account proto name to avoid re-computing it, it is also passed the account in case
// casting is needed.
type AccountBumpFilterFunc = func(accountProtoName string, account authtypes.AccountI) (bool, error)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramstore paramtypes.Subspace
	authority  string // address of the authorized actor that can execute consensus msgs

	accountKeeper      types.AccountKeeper
	rollapParamsKeeper types.RollappParamsKeeper
	accountBumpFilters []AccountBumpFilterFunc
	upgradeKeeper      upgradekeeper.Keeper

	whitelistedRelayers collections.Map[sdk.ValAddress, types.WhitelistedRelayers]
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	authority string,
	accountKeeper types.AccountKeeper,
	rollapParamsKeeper types.RollappParamsKeeper,
	upgradeKeeper upgradekeeper.Keeper,
	accountBumpFilters []AccountBumpFilterFunc,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	sb := collections.NewSchemaBuilder(collcompat.NewKVStoreService(storeKey))

	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		paramstore:         ps,
		authority:          authority,
		accountKeeper:      accountKeeper,
		rollapParamsKeeper: rollapParamsKeeper,
		accountBumpFilters: accountBumpFilters,
		upgradeKeeper:      upgradeKeeper,
		whitelistedRelayers: collections.NewMap(
			sb,
			types.WhitelistedRelayersPrefix(),
			"whitelisted_relayers",
			collcompat.ValAddressKey,
			collcompat.ProtoValue[types.WhitelistedRelayers](cdc),
		),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
