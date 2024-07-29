package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/dymint/da"
	"github.com/dymensionxyz/dymint/da/registry"
)

var (
	DefaultDA      = "celestia"
	DefaultVersion = ""
	KeyDa          = []byte("da")
	KeyVersion     = []byte("version")
)

// ParamKeyTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	da string,
	version string,
) Params {
	return Params{
		Da:      da,
		Version: version,
	}
}

// DefaultParams returns default x/denommetadata module parameters.
func DefaultParams() Params {
	return Params{
		Da:      DefaultDA,
		Version: DefaultVersion,
	}
}

func (p Params) Validate() error {
	err := assertValidDa(p.Da)
	if err != nil {
		return err
	}
	err = assertValidVersion(p.Version)
	if err != nil {
		return err
	}
	return nil
}

func assertValidDa(i any) error {
	if registry.GetClient(i.(string)) == nil {
		return da.ErrNonexistentDA
	}
	return nil

}

func assertValidVersion(i any) error {
	return nil
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDa, &p.Da, assertValidDa),
		paramtypes.NewParamSetPair(KeyVersion, &p.Version, assertValidVersion),
	}
}
