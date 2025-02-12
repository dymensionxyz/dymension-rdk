package types

import (
	"fmt"

	epochtypes "github.com/dymensionxyz/dymension-rdk/x/epochs/types"
)

func DefaultParams() Params {
	return Params{
		DistrEpochIdentifier: "day",
	}
}

func (p Params) Validate() error {
	err := epochtypes.ValidateEpochIdentifierString(p.DistrEpochIdentifier)
	if err != nil {
		return fmt.Errorf("validate distribution epoch identifier: %w", err)
	}

	return nil
}
