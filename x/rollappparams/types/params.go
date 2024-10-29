package types

import (
	"fmt"
	"regexp"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/dymint/da/registry"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

const (
	// length of the version commit string.
	CommitLength = 40
	// Data availability used by the RollApp. Default value used is mock da.
	DefaultDA = "mock"
)

// Parameter store keys.
var (
	KeyDa      = []byte("da")
	KeyCommit  = []byte("commit")
	KeyVersion = []byte("version")
	// git commit for the version used for the rollapp binary. it must be overwritten in the build process
	Version = uint64(0)
	Commit  = "<commit>"
	// default max block size accepted (equivalent to block max size it can fit into a celestia blob).
	// regexp used to validate version commit
	CommitRegExp = regexp.MustCompile(`^[a-z0-9]*$`)
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	da string,
	version uint64,
	commit string,
) Params {
	return Params{
		Da:      da,
		Version: version,
		Commit:  commit,
	}
}

// DefaultParams returns default x/rollappparams module parameters.
func DefaultParams() Params {
	return Params{
		Da:      DefaultDA,
		Version: Version,
		Commit:  Commit,
	}
}

func (p Params) Validate() error {
	err := ValidateDa(p.Da)
	if err != nil {
		return err
	}
	err = ValidateVersion(p.Version)
	if err != nil {
		return err
	}
	err = ValidateCommit(p.Commit)
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
	version, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid version type param type: %w", gerrc.ErrInvalidArgument)
	}
	if version == 0 {
		return fmt.Errorf("invalid DRS version: Version %d: %w", version, gerrc.ErrInvalidArgument)
	}

	return nil

}

func ValidateCommit(i any) error {

	commit, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid commit type param type: %w", gerrc.ErrInvalidArgument)
	}
	if len(commit) != CommitLength {
		return fmt.Errorf("invalid commit length: param length: %d accepted: %d: %w", len(commit), CommitLength, gerrc.ErrInvalidArgument)
	}
	if !CommitRegExp.MatchString(commit) {
		return fmt.Errorf("invalid commit: it must be alphanumeric %w", gerrc.ErrInvalidArgument)
	}

	return nil
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDa, &p.Da, ValidateDa),
		paramtypes.NewParamSetPair(KeyVersion, &p.Version, ValidateVersion),
		paramtypes.NewParamSetPair(KeyCommit, &p.Commit, ValidateCommit),
	}
}
