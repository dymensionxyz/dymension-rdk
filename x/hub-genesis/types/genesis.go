package types

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(params Params, hub Hub) *GenesisState {
	return &GenesisState{
		Params: params,
		Hub:    hub,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), Hub{})
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}
