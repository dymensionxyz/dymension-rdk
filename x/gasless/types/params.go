package types

import (
	fmt "fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	GasTankAddressPrefix         = "GasTankAddress"
	ModuleAddressNameSplitter    = "|"
	MaximumAuthorizedActorsLimit = 5
)

// gasless module's params default values
var (
	DefaultMinimumGasDeposit = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(10_000_000)))
)

var (
	KeyMinimumGasDeposit = []byte("MinimumGasDeposit")
)

var _ paramstypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(minGasDeposit sdk.Coins) Params {
	return Params{
		MinimumGasDeposit: minGasDeposit,
	}
}

// DefaultParams returns a default params for the liquidity module.
func DefaultParams() Params {
	return NewParams(DefaultMinimumGasDeposit)
}

// ParamSetPairs implements ParamSet.
func (params *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyMinimumGasDeposit, &params.MinimumGasDeposit, validateMinimumGasDeposit),
	}
}

// Validate validates Params.
func (p Params) Validate() error {
	if err := validateMinimumGasDeposit(p.MinimumGasDeposit); err != nil {
		return fmt.Errorf("invalid minimum gas deposit: %w", err)
	}
	return nil
}

func validateMinimumGasDeposit(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := v.Validate(); err != nil {
		return fmt.Errorf("invalid minimum gas deposit fee: %w", err)
	}

	return nil
}
