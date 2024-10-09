package types

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:          DefaultParams(),
		GenesisAccounts: []GenesisAccount{},
	}
}

// ValidateBasic performs basic validation of the genesis state.
func (g GenesisState) ValidateBasic() error {
	if err := g.Params.Validate(); err != nil {
		return err
	}

	for _, acc := range g.GenesisAccounts {
		if err := acc.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}
