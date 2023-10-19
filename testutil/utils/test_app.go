package utils

import (
	"encoding/json"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	etherencoding "github.com/evmos/evmos/v12/encoding"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"github.com/dymensionxyz/dymension-rdk/app"
	"github.com/dymensionxyz/dymension-rdk/app/params"
	"github.com/tendermint/tendermint/libs/log"

	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"

	seqcli "github.com/dymensionxyz/dymension-rdk/x/sequencers/client/cli"
	seqtypes "github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   -1,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

func setup(withGenesis bool, invCheckPeriod uint, isEVM bool) (*app.App, app.GenesisState) {
	db := dbm.NewMemDB()

	encCdc := app.MakeEncodingConfig()
	if isEVM {
		ethEncodingConfig := etherencoding.MakeConfig(app.ModuleBasics)
		encCdc = params.EncodingConfig{
			InterfaceRegistry: ethEncodingConfig.InterfaceRegistry,
			Codec:             ethEncodingConfig.Codec,
			TxConfig:          ethEncodingConfig.TxConfig,
			Amino:             ethEncodingConfig.Amino,
		}
	}
	testApp := app.NewRollapp(
		log.NewNopLogger(), db, nil, true, map[int64]bool{}, app.DefaultNodeHome, invCheckPeriod, encCdc, EmptyAppOptions{},
	)
	if withGenesis {
		return testApp, app.NewDefaultGenesisState(encCdc.Codec)
	}
	return testApp, app.GenesisState{}
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(t *testing.T, isCheckTx bool) *app.App {
	t.Helper()
	pks := CreateTestPubKeys(1)
	pk, err := cryptocodec.ToTmProtoPublicKey(pks[0])
	require.NoError(t, err)
	seq, err := seqtypes.NewSequencer(sdk.ValAddress("dsadas"), pks[0], 1)
	require.NoError(t, err)

	app, genesisState := setup(true, 5, true)

	//setting genesis sequencer as returned later from InitChain
	genesisState, err = seqcli.AddSequencerToGenesis(app.AppCodec(), genesisState, seq)
	require.NoError(t, err)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)
	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			Time:            time.Time{},
			ChainId:         "test_100-1",
			ConsensusParams: DefaultConsensusParams,
			Validators:      []abci.ValidatorUpdate{{PubKey: pk, Power: 1}},
			AppStateBytes:   stateBytes,
			InitialHeight:   0,
		},
	)

	return app
}
