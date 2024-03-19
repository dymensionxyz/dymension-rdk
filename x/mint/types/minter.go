package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	InflationAnnualDuration = time.Duration(365 * 24 * 60 * time.Minute)
)

// NewMinter returns a new Minter object with the given epoch
// provisions values.
func NewMinter(inflationRate sdk.Dec) Minter {
	return Minter{
		CurrentInflationRate: inflationRate,
	}
}

// InitialMinter returns an initial Minter object.
func InitialMinter() Minter {
	return NewMinter(sdk.NewDecWithPrec(8, 2)) // 8%
}

// DefaultInitialMinter returns a default initial Minter object for a new chain.
func DefaultInitialMinter() Minter {
	return InitialMinter()
}

// validate minter.
func ValidateMinter(minter Minter) error {
	return validateInflationRate(minter.CurrentInflationRate)
}
