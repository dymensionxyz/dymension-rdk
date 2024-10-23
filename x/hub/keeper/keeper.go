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
	if err := k.registeredHubDenoms.Set(ctx, denom); err != nil {
		return err
	}
	return nil
}

func (k Keeper) HasHubDenom(ctx sdk.Context, denom string) (bool, error) {
	ok, err := k.registeredHubDenoms.Has(ctx, denom)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (k Keeper) GetAllHubDenoms(ctx sdk.Context) ([]string, error) {
	var denoms []string
	if err := k.IterateHubDenoms(ctx, func(denom string) (bool, error) {
		denoms = append(denoms, denom)
		return false, nil
	}); err != nil {
		return nil, err
	}
	return denoms, nil
}

func (k Keeper) IterateHubDenoms(ctx sdk.Context, cb func(denom string) (bool, error)) error {
	iter, err := k.registeredHubDenoms.Iterate(ctx, new(collections.Range[string]))
	if err != nil {
		return err
	}
	defer iter.Close()

	for iter.Valid() {
		denom, err := iter.Key()
		if err != nil {
			return err
		}
		stop, err := cb(denom)
		if err != nil {
			return err
		}
		if stop {
			break
		}
		iter.Next()
	}
	return nil
}
