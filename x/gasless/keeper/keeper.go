package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/dymensionxyz/dymension-rdk/utils/collcompat"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

// Keeper of the gasless store.
type Keeper struct {
	cdc               codec.BinaryCodec
	storeKey          storetypes.StoreKey
	paramSpace        paramstypes.Subspace
	interfaceRegistry codectypes.InterfaceRegistry

	// accountKeeper types.AccountKeeper
	bankKeeper types.BankKeeper
	wasmKeeper *wasmkeeper.Keeper

	usageIdentifierToGasTankIDSet collections.KeySet[collections.Pair[string, uint64]]
	lastUsedGasTankIDMap          collections.Map[string, uint64]
}

// NewKeeper creates a new gasless Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramSpace paramstypes.Subspace,
	interfaceRegistry codectypes.InterfaceRegistry,
	bankKeeper types.BankKeeper,
	wasmKeeper *wasmkeeper.Keeper,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		paramSpace:        paramSpace,
		interfaceRegistry: interfaceRegistry,
		bankKeeper:        bankKeeper,
		wasmKeeper:        wasmKeeper,
		usageIdentifierToGasTankIDSet: collections.NewKeySet[collections.Pair[string, uint64]](
			collections.NewSchemaBuilder(collcompat.NewKVStoreService(storeKey)),
			types.UsageIdentifierToGasTankIdsKeyPrefix,
			"usageIdentifierToGasTankID",
			collections.PairKeyCodec(collections.StringKey, collections.Uint64Key)),
		lastUsedGasTankIDMap: collections.NewMap(
			collections.NewSchemaBuilder(collcompat.NewKVStoreService(storeKey)),
			types.LastUsedGasTankKey,
			"lastUsedGasTankID",
			collections.StringKey,
			collections.Uint64Value,
		),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
