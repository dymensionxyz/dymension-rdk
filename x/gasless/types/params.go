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
	MaximumTankCreationLimit     = uint64(10)
	MaximumAuthorizedActorsLimit = 5
)

// gasless module's params default values
var (
	DefaultTankCreationLimit = uint64(5)
	DefaultMinimumGasDeposit = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(10_000_000)))
)

var (
	KeyTankCreationLimit = []byte("TankCreationLimit")
	KeyMinimumGasDeposit = []byte("MinimumGasDeposit")
)

var _ paramstypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(tankCreationLimit uint64, minGasDeposit sdk.Coins) Params {
	return Params{
		TankCreationLimit: tankCreationLimit,
		MinimumGasDeposit: minGasDeposit,
	}
}

// DefaultParams returns a default params for the liquidity module.
func DefaultParams() Params {
	return NewParams(DefaultTankCreationLimit, DefaultMinimumGasDeposit)
}

// ParamSetPairs implements ParamSet.
func (params *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyTankCreationLimit, &params.TankCreationLimit, validateTankCreationLimit),
		paramstypes.NewParamSetPair(KeyMinimumGasDeposit, &params.MinimumGasDeposit, validateMinimumGasDeposit),
	}
}

// Validate validates Params.
func (p Params) Validate() error {
	if err := validateTankCreationLimit(p.TankCreationLimit); err != nil {
		return fmt.Errorf("invalid tank creation limit: %w", err)
	}
	if err := validateMinimumGasDeposit(p.MinimumGasDeposit); err != nil {
		return fmt.Errorf("invalid minimum gas deposit: %w", err)
	}
	return nil
}

func validateTankCreationLimit(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("tank creation limit must be positive: %d", v)
	}

	if v > MaximumTankCreationLimit {
		return fmt.Errorf("maximum tank creation allowed limit is : %d", MaximumTankCreationLimit)
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
