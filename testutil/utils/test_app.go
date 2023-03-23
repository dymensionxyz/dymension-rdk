package utils

import (
	"encoding/json"
	"time"

	dbm "github.com/tendermint/tm-db"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"

	"github.com/dymensionxyz/rollapp/app"
	"github.com/tendermint/tendermint/libs/log"

	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
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

func setup(withGenesis bool, invCheckPeriod uint) (*app.App, app.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := app.MakeEncodingConfig()
	testApp := app.NewRollapp(
		log.NewNopLogger(), db, nil, true, map[int64]bool{}, app.DefaultNodeHome, invCheckPeriod, encCdc, EmptyAppOptions{},
	)
	if withGenesis {
		return testApp, app.NewDefaultGenesisState(encCdc.Marshaler)
	}
	return testApp, app.GenesisState{}
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(isCheckTx bool) *app.App {
	testApp, genesisState := setup(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		pks := CreateTestPubKeys(1)

		pk, err := cryptocodec.ToTmProtoPublicKey(pks[0])
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		(*testApp).InitChain(
			abci.RequestInitChain{
				Validators: []abci.ValidatorUpdate{abci.ValidatorUpdate{
					PubKey: pk,
					Power:  1,
				}},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return testApp
}
