package types

import (
	"encoding/json"
	"fmt"
	"os"
)

// GenesisDocFromFile reads JSON data from a file and unmarshalls it into a GenesisDoc.
func GenesisDocFromFile(genDocFile string) (map[string]interface{}, error) {
	jsonBlob, err := os.ReadFile(genDocFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't read GenesisDoc file: %w", err)
	}

	genDoc, err := GenesisDocFromJSON(jsonBlob)
	if err != nil {
		return nil, fmt.Errorf("error reading GenesisDoc at %s: %w", genDocFile, err)
	}
	return genDoc, nil
}

// GenesisDocFromJSON unmarshalls JSON data into a GenesisDoc.
func GenesisDocFromJSON(jsonBlob []byte) (map[string]interface{}, error) {
	genDoc := make(map[string]interface{})
	err := json.Unmarshal(jsonBlob, &genDoc)
	if err != nil {
		return nil, err
	}

	return genDoc, err
}
