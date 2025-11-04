package keeper

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/utils/collcompat"
	"github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

type Keeper struct {
	cdc                   codec.BinaryCodec
	storeKey              storetypes.StoreKey
	registeredHubDenoms   collections.KeySet[string]
	decimalConversionPair collections.Item[types.DecimalConversionPair]
	bankKeeper            types.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(collcompat.NewKVStoreService(storeKey))
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		bankKeeper: bankKeeper,
		registeredHubDenoms: collections.NewKeySet(
			sb,
			collections.NewPrefix(types.RegisteredHubDenomsKeyPrefix),
			"registered_hub_denoms",
			collections.StringKey,
		),
		decimalConversionPair: collections.NewItem(
			sb,
			collections.NewPrefix(types.DecimalConversionPairKeyPrefix),
			"decimal_conversion_pair",
			collcompat.ProtoValue[types.DecimalConversionPair](cdc),
		),
	}
}

func (k Keeper) SetHubDenom(ctx sdk.Context, denom string) error {
	return k.registeredHubDenoms.Set(ctx, denom)
}

func (k Keeper) HasHubDenom(ctx sdk.Context, denom string) (bool, error) {
	return k.registeredHubDenoms.Has(ctx, denom)
}

func (k Keeper) GetAllHubDenoms(ctx sdk.Context) (denoms []string, _ error) {
	return denoms, k.registeredHubDenoms.Walk(ctx, new(collections.Range[string]), func(d string) (_ bool, _ error) {
		denoms = append(denoms, d)
		return
	})
}

// SetDecimalConversionPair sets the decimal conversion pair
func (k Keeper) SetDecimalConversionPair(ctx sdk.Context, pair types.DecimalConversionPair) error {
	return k.decimalConversionPair.Set(ctx, pair)
}

// GetDecimalConversionPair retrieves the decimal conversion pair
func (k Keeper) GetDecimalConversionPair(ctx sdk.Context) (types.DecimalConversionPair, error) {
	return k.decimalConversionPair.Get(ctx)
}
