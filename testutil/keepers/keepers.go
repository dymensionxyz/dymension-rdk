package keepers

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	epochkeeper "github.com/dymensionxyz/dymension-rdk/x/epochs/keeper"
	mintkeeper "github.com/dymensionxyz/dymension-rdk/x/mint/keeper"
	seqkeeper "github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"

	app "github.com/dymensionxyz/dymension-rdk/testutil/app"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func NewTestEpochKeeperFromApp(t *testing.T, app *app.App) (*epochkeeper.Keeper, sdk.Context) {
	k := &app.EpochsKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})

	return k, ctx
}

func NewTestSequencerKeeperFromApp(t *testing.T, app *app.App) (*seqkeeper.Keeper, sdk.Context) {
	k := &app.SequencersKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

func NewTestMintKeeperFromApp(t *testing.T, app *app.App) (*mintkeeper.Keeper, sdk.Context) {
	k := &app.MintKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

// TODO: when dymension release denommetadata module, replace this with denommetadata keeper and remove this
func NewTestBankKeeperFromApp(t *testing.T, app *app.App) (*bankkeeper.Keeper, sdk.Context) {
	k := &app.BankKeeper
	bankStoreKey := storetypes.NewKVStoreKey(banktypes.StoreKey)
	bank_t_storeKey := storetypes.NewTransientStoreKey("t_" + banktypes.StoreKey)
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(bankStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(bank_t_storeKey, storetypes.StoreTypeTransient, nil)
	if err := stateStore.LoadLatestVersion(); err != nil {
		panic(fmt.Errorf("failed to load latest version: %w", err))
	}
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()}).WithMultiStore(stateStore)
	return k, ctx
}
