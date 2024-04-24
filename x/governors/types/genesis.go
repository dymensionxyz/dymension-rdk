package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// NewGenesisState creates a new GenesisState instanc e
func NewGenesisState(params Params, Governors []Governor, delegations []stakingtypes.Delegation) *GenesisState {
	return &GenesisState{
		Params:      params,
		Governors:   Governors,
		Delegations: delegations,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// GetGenesisStateFromAppState returns x/staking GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
