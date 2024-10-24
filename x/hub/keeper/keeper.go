package keeper

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/utils/collcompat"
	hubtypes "github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

type Keeper struct {
	cdc                 codec.BinaryCodec
	storeKey            storetypes.StoreKey
	registeredHubDenoms collections.KeySet[string]
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		registeredHubDenoms: collections.NewKeySet(
			collections.NewSchemaBuilder(collcompat.NewKVStoreService(storeKey)),
			collections.NewPrefix(hubtypes.RegisteredHubDenomsKeyPrefix),
			"registered_hub_denoms",
			collections.StringKey,
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
