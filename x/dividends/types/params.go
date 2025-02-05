package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	epochtypes "github.com/dymensionxyz/dymension-rdk/x/epochs/types"
)

func DefaultParams() Params {
	return Params{
		DistrEpochIdentifier: "day",
		ApprovedDenoms:       []string{sdk.DefaultBondDenom},
	}
}

func (p Params) Validate() error {
	err := epochtypes.ValidateEpochIdentifierString(p.DistrEpochIdentifier)
	if err != nil {
		return fmt.Errorf("validate distribution epoch identifier: %w", err)
	}

	for _, denom := range p.ApprovedDenoms {
		if err := sdk.ValidateDenom(denom); err != nil {
			return fmt.Errorf("validate approved denom: %w", err)
		}
	}

	return nil
}
