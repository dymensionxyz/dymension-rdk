package keepers

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/rollapp/app"
	epochkeeper "github.com/dymensionxyz/rollapp/x/epochs/keeper"
	seqkeeper "github.com/dymensionxyz/rollapp/x/sequencers/keeper"

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
