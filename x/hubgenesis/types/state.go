package types

// validate state
func (s State) Validate() error {
	if !s.GenesisTokens.IsValid() {
		return ErrInvalidGenesisTokens
	}
	return nil
}
