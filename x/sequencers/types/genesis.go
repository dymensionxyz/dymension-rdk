package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

var DefaultAddressPermissions = []AddressPermissions{} // no one allowed

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:             DefaultParams(),
		AddressPermissions: DefaultAddressPermissions,
	}
}

func (gs GenesisState) ValidateGenesis() error {
	err := gs.Params.Validate()
	if err != nil {
		return err
	}

	for _, accPerms := range gs.AddressPermissions {
		if err := accPerms.Validate(); err != nil {
			return err
		}
	}

	_, err = sdk.ValAddressFromBech32(gs.GenesisOperatorAddress)
	if err != nil {
		return fmt.Errorf("genesis operator address is invalid: %w", err)
	}

	return nil
}
