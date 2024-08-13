package types

import (
	"fmt"
	"regexp"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/dymint/da/registry"
)

const (
	DefaultDA               = "celestia"
	DefaultCommit           = "74fad6a00713cba62352c2451c6b7ab73571c515"
	DefaultBlockMaxGas      = 400000000
	DefaultBlockMaxSize     = 500000
	MinAcceptedBlockMaxSize = 100000
	MinAcceptedBlockMaxGas  = 1000000
	CommitLength            = 40
)

// Parameter store keys.
var (
	KeyDa           = []byte("da")
	KeyCommit       = []byte("commit")
	KeyBlockMaxGas  = []byte("blockmaxgas")
	KeyBlockMaxSize = []byte("blockmaxsize")
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	da string,
	commit string,
	blockMaxGas uint32,
	blockMaxSize uint32,
) Params {
	return Params{
		Da:           da,
		Commit:       commit,
		Blockmaxgas:  blockMaxGas,
		Blockmaxsize: blockMaxSize,
	}
}

// DefaultParams returns default x/rollappparams module parameters.
func DefaultParams() Params {
	return Params{
		Da:           DefaultDA,
		Commit:       DefaultCommit,
		Blockmaxgas:  uint32(DefaultBlockMaxGas),
		Blockmaxsize: uint32(DefaultBlockMaxSize),
	}
}

func (p Params) Validate() error {
	err := ValidateDa(p.Da)
	if err != nil {
		return err
	}
	err = ValidateCommit(p.Commit)
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
		return ErrDANotSupported
	}

	return nil

}

func ValidateCommit(i any) error {

	commit, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid commit")
	}
	if len(commit) != CommitLength {
		return fmt.Errorf("invalid commit length")
	}
	if !regexp.MustCompile(`^[a-z0-9]*$`).MatchString(commit) {
		return fmt.Errorf("commit must be alphanumeric")

	}

	return nil
}

func ValidateBlockMaxGas(i any) error {
	gas, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid block max gas")
	}
	if gas < uint32(MinAcceptedBlockMaxGas) {
		return fmt.Errorf("block max gas cannot be smaller than %d", MinAcceptedBlockMaxGas)
	}
	return nil
}

func ValidateBlockMaxSize(i any) error {
	size, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid block max size")
	}
	if size < uint32(MinAcceptedBlockMaxSize) {
		return fmt.Errorf("block max size cannot be smaller than %d", MinAcceptedBlockMaxSize)
	}
	if size > uint32(DefaultBlockMaxSize) {
		return fmt.Errorf("block max size cannot be greater than %d", DefaultBlockMaxSize)
	}
	return nil
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	fmt.Println(p)
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDa, &p.Da, ValidateDa),
		paramtypes.NewParamSetPair(KeyCommit, &p.Commit, ValidateCommit),
		paramtypes.NewParamSetPair(KeyBlockMaxGas, &p.Blockmaxgas, ValidateBlockMaxGas),
		paramtypes.NewParamSetPair(KeyBlockMaxSize, &p.Blockmaxsize, ValidateBlockMaxSize),
	}
}
