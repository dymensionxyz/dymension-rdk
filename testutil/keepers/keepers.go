package keepers

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	epochkeeper "github.com/dymensionxyz/dymension-rdk/x/epochs/keeper"
	mintkeeper "github.com/dymensionxyz/dymension-rdk/x/mint/keeper"
	seqkeeper "github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	testutils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
)

func NewTestEpochKeeperFromApp(t *testing.T, app *testutils.App) (*epochkeeper.Keeper, sdk.Context) {
	k := &app.EpochsKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})

	return k, ctx
}

func NewTestSequencerKeeperFromApp(t *testing.T, app *testutils.App) (*seqkeeper.Keeper, sdk.Context) {
	k := &app.SequencersKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}

func NewTestMintKeeperFromApp(t *testing.T, app *testutils.App) (*mintkeeper.Keeper, sdk.Context) {
	k := &app.MintKeeper
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "rollapp-1", Time: time.Now().UTC()})
	return k, ctx
}
