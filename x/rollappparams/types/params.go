package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/dymint/da/registry"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

const (

	// Data availability used by the RollApp. Default value used is mock da.
	DefaultDA = "mock"
)

// Parameter store keys.
var (
	KeyDa      = []byte("da")
	KeyVersion = []byte("version")
	// Default version set
	DrsVersion = uint32(1)
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	da string,
	drsVersion uint32,
) Params {
	return Params{
		Da:         da,
		DrsVersion: drsVersion,
	}
}

// DefaultParams returns default x/rollappparams module parameters.
func DefaultParams() Params {
	return Params{
		Da:         DefaultDA,
		DrsVersion: DrsVersion,
	}
}

func (p Params) Validate() error {
	err := ValidateDa(p.Da)
	if err != nil {
		return err
	}
	err = ValidateVersion(p.DrsVersion)
	if err != nil {
		return err
	}

	return nil
}

func ValidateDa(i any) error {
	if registry.GetClient(i.(string)) == nil {
		return fmt.Errorf("invalid DA type: DA %s: %w", i, gerrc.ErrInvalidArgument)
	}

	return nil

}

func ValidateVersion(i any) error {
	version, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid version type param type: %w", gerrc.ErrInvalidArgument)
	}
	if version <= 0 {
		return fmt.Errorf("invalid DRS version: Version must be positive")
	}

	return nil

}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDa, &p.Da, ValidateDa),
		paramtypes.NewParamSetPair(KeyVersion, &p.DrsVersion, ValidateVersion),
	}
}
