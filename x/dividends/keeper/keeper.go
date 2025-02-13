package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/utils/collcompat"
	"github.com/dymensionxyz/dymension-rdk/x/dividends/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	authority string // authority is the x/gov module account

	getBalanceFn GetGaugeBalanceFunc

	schema      collections.Schema
	params      collections.Item[types.Params]
	lastGaugeID collections.Sequence
	// [status, GaugeID] -> Gauge; status is true for active, false for inactive
	gauges collections.Map[collections.Pair[bool, uint64], types.Gauge]

	stakingKeeper types.StakingKeeper
	accountKeeper types.AccountKeeper
	distrKeeper   types.DistributionKeeper
	bankKeeper    types.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	stakingKeeper types.StakingKeeper,
	accountKeeper types.AccountKeeper,
	distrKeeper types.DistributionKeeper,
	bankKeeper types.BankKeeper,
	authority string,
) Keeper {
	sb := collections.NewSchemaBuilder(collcompat.NewKVStoreService(storeKey))

	k := Keeper{
		authority:    authority,
		getBalanceFn: nil,                  // set in SetGetBalanceFunc
		schema:       collections.Schema{}, // set later
		params: collections.NewItem(
			sb,
			types.ParamsKey,
			"params",
			collcompat.ProtoValue[types.Params](cdc),
		),
		lastGaugeID: collections.NewSequence(
			sb,
			types.LastGaugeKey,
			"last_gauge_id",
		),
		gauges: collections.NewMap(
			sb,
			types.GaugesKey,
			"gauges",
			collections.PairKeyCodec(collections.BoolKey, collections.Uint64Key),
			collcompat.ProtoValue[types.Gauge](cdc),
		),
		stakingKeeper: stakingKeeper,
		accountKeeper: accountKeeper,
		distrKeeper:   distrKeeper,
		bankKeeper:    bankKeeper,
	}

	// SchemaBuilder CANNOT be used after Build is called,
	// so we build it after all collections are initialized
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.schema = schema

	return k
}

func (k Keeper) SetGetBalanceFunc(fn GetGaugeBalanceFunc) {
	k.getBalanceFn = fn
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Schema() collections.Schema {
	return k.schema
}
