package types

import (
	fmt "fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	InflationAnnualDuration = time.Duration(365 * 24 * 60 * time.Minute)
)

// NewMinter returns a new Minter object with the given epoch
// provisions values.
func NewMinter(denom string, inflationRate sdk.Dec) Minter {
	return Minter{
		MintDenom:            denom,
		CurrentInflationRate: inflationRate,
	}
}

// InitialMinter returns an initial Minter object.
func InitialMinter() Minter {
	return NewMinter(sdk.DefaultBondDenom, sdk.NewDecWithPrec(8, 2)) // 8%
}

// DefaultInitialMinter returns a default initial Minter object for a new chain.
func DefaultInitialMinter() Minter {
	return InitialMinter()
}

// validate minter.
func ValidateMinter(minter Minter) error {
	if minter.MintDenom != "" {
		err := sdk.ValidateDenom(minter.MintDenom)
		if err != nil {
			return err
		}

		// validate it's not ibc or tokenfactory
		if strings.HasPrefix(strings.ToLower(minter.MintDenom), "ibc") || strings.HasPrefix(strings.ToLower(minter.MintDenom), "factory") {
			return fmt.Errorf("denom is not allowed in minter (%s)", minter.MintDenom)
		}
	}

	return validateInflationRate(minter.CurrentInflationRate)
}
