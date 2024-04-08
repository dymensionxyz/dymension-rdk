package types

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(params Params, locked Locked) *GenesisState {
	return &GenesisState{
		Params: params,
		Locked: locked,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), Locked{})
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	if err := data.Locked.Validate(); err != nil {
		return err
	}

	return nil
}
