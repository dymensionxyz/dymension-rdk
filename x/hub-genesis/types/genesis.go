package types

import (
	"errors"
	"fmt"

	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

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

	accountSet := make(map[string]struct{})
	for _, a := range g.GenesisAccounts {
		if err := a.ValidateBasic(); err != nil {
			return errors.Join(gerrc.ErrInvalidArgument, err)
		}
		if _, exists := accountSet[a.Address]; exists {
			return fmt.Errorf("duplicate genesis account: %s", a.Address)
		}
		accountSet[a.Address] = struct{}{}
	}

	return nil
}
