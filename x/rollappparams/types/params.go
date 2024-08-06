package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/dymint/da/registry"
)

var (
	DefaultDA               = "celestia"
	DefaultCommit           = "74fad6a00713cba62352c2451c6b7ab73571c515"
	DefaultMaxBlockGas      = 400000000
	DefaultMaxBlockSize     = 500000
	MinAcceptedMaxBlockSize = 100000
	KeyDa                   = []byte("da")
	KeyCommit               = []byte("commit")
	KeyBlockMaxGas          = []byte("blockmaxgas")
	KeyBlockMaxSize         = []byte("blockmaxsize")
	CommitLength            = 40
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
		Blockmaxgas:  uint32(DefaultMaxBlockGas),
		Blockmaxsize: uint32(DefaultMaxBlockSize),
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
	err = assertValidBlockMaxGas(p.Blockmaxgas)
	if err != nil {
		return err
	}
	err = assertValidBlockMaxSize(p.Blockmaxsize)
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

func assertValidBlockMaxGas(i any) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid block max gas")
	}

	return nil
}

func assertValidBlockMaxSize(i any) error {
	size, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid block max size")
	}
	if size < uint32(MinAcceptedMaxBlockSize) {
		return fmt.Errorf("block max size cannot be smaller than %d", MinAcceptedMaxBlockSize)
	}
	if size > uint32(DefaultMaxBlockSize) {
		return fmt.Errorf("block max size cannot be greater than %d", DefaultMaxBlockSize)
	}
	return nil
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	fmt.Println(p)
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDa, &p.Da, assertValidDa),
		paramtypes.NewParamSetPair(KeyCommit, &p.Commit, assertValidCommit),
		paramtypes.NewParamSetPair(KeyBlockMaxGas, &p.Blockmaxgas, assertValidBlockMaxGas),
		paramtypes.NewParamSetPair(KeyBlockMaxSize, &p.Blockmaxsize, assertValidBlockMaxSize),
	}
}
