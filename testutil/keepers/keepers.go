package keepers

import (
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	epochkeeper "github.com/dymensionxyz/dymension-rdk/x/epochs/keeper"
	hubgenkeeper "github.com/dymensionxyz/dymension-rdk/x/hub-genesis/keeper"
	hubkeeper "github.com/dymensionxyz/dymension-rdk/x/hub/keeper"
	mintkeeper "github.com/dymensionxyz/dymension-rdk/x/mint/keeper"
	rollappparamskeeper "github.com/dymensionxyz/dymension-rdk/x/rollappparams/keeper"
	seqkeeper "github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	timeupgradekeeper "github.com/dymensionxyz/dymension-rdk/x/timeupgrade/keeper"
)

func NewTestEpochKeeperFromApp(app *app.App) (*epochkeeper.Keeper, sdk.Context) {
	k := &app.EpochsKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

func NewTestSequencerKeeperFromApp(app *app.App) (*seqkeeper.Keeper, sdk.Context) {
	k := &app.SequencersKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

func NewTestTimeupgradeKeeperFromApp(app *app.App) (timeupgradekeeper.Keeper, sdk.Context) {
	k := app.TimeUpgradeKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

func NewTestUpgradeKeeperFromApp(app *app.App) (upgradekeeper.Keeper, sdk.Context) {
	k := app.UpgradeKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

func NewTestMintKeeperFromApp(app *app.App) (*mintkeeper.Keeper, sdk.Context) {
	k := &app.MintKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

func NewTestHubGenesisKeeperFromApp(app *app.App) (*hubgenkeeper.Keeper, sdk.Context) {
	k := &app.HubGenesisKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

func NewTestRollappParamsKeeperFromApp(app *app.App) (*rollappparamskeeper.Keeper, sdk.Context) {
	k := &app.RollappParamsKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

func NewTestHubKeeperFromApp(app *app.App) (*hubkeeper.Keeper, sdk.Context) {
	k := &app.HubKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-]1", Time: time.Now().UTC()})
	return k, ctx
}
