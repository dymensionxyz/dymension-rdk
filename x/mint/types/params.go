package types

import (
	"errors"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	epochtypes "github.com/dymensionxyz/dymension-rdk/x/epochs/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	KeyMintDenom                = []byte("MintDenom")
	KeyMintEpochIdentifier      = []byte("MintEpochIdentifier")
	KeyMintEpochSpreadFactor    = []byte("MintEpochSpreadFactor")
	KeyMintStartEpoch           = []byte("MintStartEpoch")
	KeyInflationEpochIdentifier = []byte("InflationEpochIdentifier")
	KeyInflationRateChange      = []byte("InflationRateChange")
	KeyTargetInflationRate      = []byte("TargetInflationRate")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintDenom string, mintEpochIdentifier string, mintEpochSpreadFactor int64,
	mintStartEpoch int64, inflationEpochIdentifier string,
	inflationRateChange sdk.Dec, targetInflationRate sdk.Dec,
) Params {
	return Params{
		MintDenom:                mintDenom,
		MintEpochIdentifier:      mintEpochIdentifier,
		MintEpochSpreadFactor:    mintEpochSpreadFactor,
		MintStartEpoch:           mintStartEpoch,
		InflationEpochIdentifier: inflationEpochIdentifier,
		InflationRateChange:      inflationRateChange,
		TargetInflationRate:      targetInflationRate,
	}
}

// minting params
func DefaultParams() Params {
	return Params{
		MintDenom:                sdk.DefaultBondDenom,
		MintEpochIdentifier:      "day",
		MintEpochSpreadFactor:    365,
		MintStartEpoch:           1,
		InflationEpochIdentifier: "year",
		InflationRateChange:      sdk.NewDecWithPrec(10, 2), // 10% inflation change
		TargetInflationRate:      sdk.NewDecWithPrec(2, 2),  // 2%
	}
}

// validate params.
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := epochtypes.ValidateEpochIdentifierInterface(p.MintEpochIdentifier); err != nil {
		return err
	}
	if err := validateInt(p.MintEpochSpreadFactor); err != nil {
		return err
	}
	if err := epochtypes.ValidateEpochIdentifierInterface(p.InflationEpochIdentifier); err != nil {
		return err
	}
	if err := validateInflationRate(p.InflationRateChange); err != nil {
		return err
	}
	if err := validateInflationRate(p.TargetInflationRate); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyMintEpochIdentifier, &p.MintEpochIdentifier, epochtypes.ValidateEpochIdentifierInterface),
		paramtypes.NewParamSetPair(KeyInflationEpochIdentifier, &p.InflationEpochIdentifier, epochtypes.ValidateEpochIdentifierInterface),
		paramtypes.NewParamSetPair(KeyMintEpochSpreadFactor, &p.MintEpochSpreadFactor, validateInt),
		paramtypes.NewParamSetPair(KeyMintStartEpoch, &p.MintStartEpoch, validateInt),
		paramtypes.NewParamSetPair(KeyInflationRateChange, &p.InflationRateChange, validateInflationRate),
		paramtypes.NewParamSetPair(KeyTargetInflationRate, &p.TargetInflationRate, validateInflationRate),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateInt(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("value must be positive: %d", v)
	}

	return nil
}

func validateInflationRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDec(1)) {
		return fmt.Errorf("inflation rate cannot be greater than 1")
	}

	if v.IsNegative() {
		return fmt.Errorf("inflation rate cannot be negative")
	}

	return nil
}
