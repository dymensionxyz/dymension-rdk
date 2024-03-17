package keepers

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	epochkeeper "github.com/dymensionxyz/dymension-rdk/x/epochs/keeper"
	mintkeeper "github.com/dymensionxyz/dymension-rdk/x/mint/keeper"
	seqkeeper "github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	"github.com/dymensionxyz/rollapp/app"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
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
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}
