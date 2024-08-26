package types

import (
	"fmt"
	"regexp"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/dymint/da/registry"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

const (
	// current supported DA
	DefaultDA = "mock"
	// version used for the rollapp binary. it must be overwritten in Makefile.
	DefaultVersion = "3a19edd887a9b576a866750bc9d480ada53d2c0d"
	// default max gas accepted per block. limited to 400M.
	DefaultBlockMaxGas = 400000000
	// default max block size accepted (equivalent to block max size it can fit into a celestia blob).
	DefaultBlockMaxSize = 500000
	// default minimum block size. not specific reason to set it to 100K, but we need to avoid no transactions can be included in a block.
	MinBlockMaxSize = 100000
	// default minimum value for max gas used in a block. set to 10M to avoid using too small values that limit performance and avoid no transactions can be included in a block.
	MinBlockMaxGas = 10000000
	// length of the version commit string.
	VersionLength = 40
)

// Parameter store keys.
var (
	KeyDa           = []byte("da")
	KeyVersion      = []byte("version")
	KeyBlockMaxGas  = []byte("blockmaxgas")
	KeyBlockMaxSize = []byte("blockmaxsize")
	VersionRegExp   = regexp.MustCompile(`^[a-z0-9]*$`)
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	da string,
	version string,
	blockMaxGas uint32,
	blockMaxSize uint32,
) Params {
	return Params{
		Da:           da,
		Version:      version,
		Blockmaxgas:  blockMaxGas,
		Blockmaxsize: blockMaxSize,
	}
}

// DefaultParams returns default x/rollappparams module parameters.
func DefaultParams() Params {
	return Params{
		Da:           DefaultDA,
		Version:      DefaultVersion,
		Blockmaxgas:  uint32(DefaultBlockMaxGas),
		Blockmaxsize: uint32(DefaultBlockMaxSize),
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
	err = ValidateBlockMaxGas(p.Blockmaxgas)
	if err != nil {
		return err
	}
	err = ValidateBlockMaxSize(p.Blockmaxsize)
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

	version, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid version type param type: %w", gerrc.ErrInvalidArgument)
	}
	if len(version) != VersionLength {
		return fmt.Errorf("invalid version length: param length: %d accepted: %d: %w", len(version), VersionLength, gerrc.ErrInvalidArgument)
	}
	if !VersionRegExp.MatchString(version) {
		return fmt.Errorf("invalid version: it must be alphanumeric %w", gerrc.ErrInvalidArgument)
	}

	return nil
}

func ValidateBlockMaxGas(i any) error {
	gas, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid block max gas param type: %w", gerrc.ErrInvalidArgument)
	}
	if gas < uint32(MinBlockMaxGas) {
		return fmt.Errorf("invalid block max gas value: used: %d minimum accepted: %d: %w", gas, MinBlockMaxGas, gerrc.ErrInvalidArgument)
	}
	return nil
}

func ValidateBlockMaxSize(i any) error {
	size, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid block max size param type : %w", gerrc.ErrInvalidArgument)
	}
	if size < uint32(MinBlockMaxSize) {
		return fmt.Errorf("invalid block max size value: used %d: minimum accepted %d : %w", size, MinBlockMaxSize, gerrc.ErrInvalidArgument)
	}
	if size > uint32(DefaultBlockMaxSize) {
		return fmt.Errorf("invalid block max size value: used %d: max accepted %d : %w", size, DefaultBlockMaxSize, gerrc.ErrInvalidArgument)
	}
	return nil
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDa, &p.Da, ValidateDa),
		paramtypes.NewParamSetPair(KeyVersion, &p.Version, ValidateVersion),
		paramtypes.NewParamSetPair(KeyBlockMaxGas, &p.Blockmaxgas, ValidateBlockMaxGas),
		paramtypes.NewParamSetPair(KeyBlockMaxSize, &p.Blockmaxsize, ValidateBlockMaxSize),
	}
}
