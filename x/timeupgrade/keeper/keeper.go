package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/collections"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	prototypes "github.com/gogo/protobuf/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dymensionxyz/dymension-rdk/utils/collcompat"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

type Keeper struct {
	cdc       codec.BinaryCodec
	storeKey  storetypes.StoreKey
	authority string

	UpgradePlan collections.Item[upgradetypes.Plan]
	UpgradeTime collections.Item[prototypes.Timestamp]
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authority string,
) Keeper {
	service := collcompat.NewKVStoreService(storeKey)
	sb := collections.NewSchemaBuilder(service)

	return Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		authority: authority,

		UpgradePlan: collections.NewItem[upgradetypes.Plan](sb, collections.NewPrefix(0), "plan", collcompat.ProtoValue[upgradetypes.Plan](cdc)),
		UpgradeTime: collections.NewItem[prototypes.Timestamp](sb, collections.NewPrefix(1), "time", collcompat.ProtoValue[prototypes.Timestamp](cdc)),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetUpgradeTime gets the upgrade time from the store
func (k Keeper) GetUpgradeTime(ctx sdk.Context) (time.Time, error) {
	upgradeTime, err := k.UpgradeTime.Get(ctx)
	if err != nil {
		return time.Time{}, err
	}

	upgradeTimeTimestamp, err := prototypes.TimestampFromProto(&upgradeTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse upgrade time: %w", err)
	}

	return upgradeTimeTimestamp, nil
}

// CleanTimeUpgrade removes the upgrade time and plan from the store
func (k Keeper) CleanTimeUpgrade(ctx sdk.Context) error {
	err := k.UpgradeTime.Remove(ctx)
	if err != nil {
		return err
	}

	err = k.UpgradePlan.Remove(ctx)
	if err != nil {
		return err
	}
	return nil
}
