package types

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(state State) *GenesisState {
	return &GenesisState{
		State: state,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(State{})
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.State.Validate(); err != nil {
		return err
	}

	return nil
}
