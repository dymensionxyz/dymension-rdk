package keeper_test

import (
	"testing"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
)

func TestEpochsInitAndExportGenesis(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	params := k.GetParams(ctx)
	seqs := k.GetAllSequencers(ctx)
	require.Equal(t, 1, len(seqs))
	expectedOperator := seqs[0].GetOperator().String()
	require.NotEmpty(t, expectedOperator)

	genState := k.ExportGenesis(ctx)
	assert.Equal(t, expectedOperator, genState.GenesisOperatorAddress)
	assert.Equal(t, params, genState.Params)

	// Test InitGenesis
	genState.Params.HistoricalEntries = 100

	// init dymint sequencer as expected for initGenesis
	pubkey := utils.CreateTestPubKeys(1)[0]
	tmPubkey, err := cryptocodec.ToTmProtoPublicKey(pubkey)
	require.NoError(t, err)
	dymintSeq := abci.ValidatorUpdate{
		PubKey: tmPubkey,
		Power:  1,
	}

	k.SetDymintSequencers(ctx, []abci.ValidatorUpdate{dymintSeq})

	_ = k.InitGenesis(ctx, *genState)
	assert.Equal(t, genState.Params, k.GetParams(ctx))
}
