package keeper

import (
	"fmt"

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
