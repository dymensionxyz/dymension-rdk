package types

import (
	"encoding/json"
	"fmt"

	tmos "github.com/tendermint/tendermint/libs/os"
)

// GenesisStateFromGenFile creates the core parameters for genesis initialization
// for the application.
//
// NOTE: The pubkey input is this machines pubkey.
func GenesisStateFromGenFile(genFile string) (genesisState map[string]json.RawMessage, genDoc map[string]interface{}, err error) {
	if !tmos.FileExists(genFile) {
		return genesisState, genDoc,
			fmt.Errorf("%s does not exist, run `init` first", genFile)
	}

	genDoc, err = GenesisDocFromFile(genFile)
	if err != nil {
		return genesisState, genDoc, err
	}

	bz, err := json.Marshal(genDoc["app_state"])
	if err != nil {
		return genesisState, genDoc, err
	}

	err = json.Unmarshal(bz, &genesisState)
	return genesisState, genDoc, err
}
