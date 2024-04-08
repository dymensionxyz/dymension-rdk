package types

// validate locked
func (l Locked) Validate() error {
	if !l.Tokens.IsValid() {
		return ErrInvalidGenesisTokens
	}
	return nil
}
