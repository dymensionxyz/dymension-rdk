package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/dymint/da/registry"
)

var (
	DefaultDA     = "celestia"
	DefaultCommit = "abcde"
	KeyDa         = []byte("da")
	KeyVersion    = []byte("commit")
	CommitLength  = 40
)

// ParamKeyTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	da string,
	commit string,
) Params {
	return Params{
		Da:     da,
		Commit: commit,
	}
}

// DefaultParams returns default x/denommetadata module parameters.
func DefaultParams() Params {
	return Params{
		Da:     DefaultDA,
		Commit: DefaultCommit,
	}
}

func (p Params) Validate() error {
	err := assertValidDa(p.Da)
	if err != nil {
		return err
	}
	err = assertValidCommit(p.Commit)
	if err != nil {
		return err
	}
	return nil
}

func assertValidDa(i any) error {
	if registry.GetClient(i.(string)) == nil {
		return ErrDANotSupported
	}
	return nil

}

func assertValidCommit(i any) error {

	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid commit")
	}
	if len(i.(string)) != CommitLength {
		return fmt.Errorf("invalid commit")
	}
	return nil
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDa, &p.Da, assertValidDa),
		paramtypes.NewParamSetPair(KeyVersion, &p.Commit, assertValidCommit),
	}
}
