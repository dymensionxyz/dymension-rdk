package keepers

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	denommetadatakeeper "github.com/dymensionxyz/dymension-rdk/x/denommetadata/keeper"
	epochkeeper "github.com/dymensionxyz/dymension-rdk/x/epochs/keeper"
	hubgenkeeper "github.com/dymensionxyz/dymension-rdk/x/hub-genesis/keeper"
	mintkeeper "github.com/dymensionxyz/dymension-rdk/x/mint/keeper"
	seqkeeper "github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	app "github.com/dymensionxyz/dymension-rdk/testutil/app"
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

func NewTestMintKeeperFromApp(app *app.App) (*mintkeeper.Keeper, sdk.Context) {
	k := &app.MintKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

func NewTestDenommetadataKeeperFromApp(app *app.App) (*denommetadatakeeper.Keeper, sdk.Context) {
	k := &app.DenommetadataKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

func NewTestHubGenesisKeeperFromApp(app *app.App) (*hubgenkeeper.Keeper, sdk.Context) {
	k := &app.HubGenesisKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}
