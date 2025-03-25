package types

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	epochtypes "github.com/dymensionxyz/dymension-rdk/x/epochs/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	KeyMintEpochIdentifier            = []byte("MintEpochIdentifier")
	KeyMintStartEpoch                 = []byte("MintStartEpoch")
	KeyInflationChangeEpochIdentifier = []byte("InflationChangeEpochIdentifier")
	KeyInflationRateChange            = []byte("InflationRateChange")
	KeyTargetInflationRate            = []byte("TargetInflationRate")
)

// ParamKeyTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintEpochIdentifier string,
	mintStartEpoch int64, inflationEpochIdentifier string,
	inflationRateChange sdk.Dec, targetInflationRate sdk.Dec,
) Params {
	return Params{
		MintEpochIdentifier:            mintEpochIdentifier,
		MintStartEpoch:                 mintStartEpoch,
		InflationChangeEpochIdentifier: inflationEpochIdentifier,
		InflationRateChange:            inflationRateChange,
		TargetInflationRate:            targetInflationRate,
	}
}

// minting params
func DefaultParams() Params {
	return Params{
		MintEpochIdentifier:            "hour",
		MintStartEpoch:                 1,
		InflationChangeEpochIdentifier: "year",
		InflationRateChange:            sdk.NewDecWithPrec(1, 2), // 1% annual inflation change
		TargetInflationRate:            sdk.NewDecWithPrec(2, 2), // 2%
	}
}

// validate params.
func (p Params) Validate() error {
	if err := epochtypes.ValidateEpochIdentifierInterface(p.MintEpochIdentifier); err != nil {
		return err
	}
	if err := epochtypes.ValidateEpochIdentifierInterface(p.InflationChangeEpochIdentifier); err != nil {
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
		paramtypes.NewParamSetPair(KeyMintEpochIdentifier, &p.MintEpochIdentifier, epochtypes.ValidateEpochIdentifierInterface),
		paramtypes.NewParamSetPair(KeyInflationChangeEpochIdentifier, &p.InflationChangeEpochIdentifier, epochtypes.ValidateEpochIdentifierInterface),
		paramtypes.NewParamSetPair(KeyMintStartEpoch, &p.MintStartEpoch, validateInt),
		paramtypes.NewParamSetPair(KeyInflationRateChange, &p.InflationRateChange, validateInflationRate),
		paramtypes.NewParamSetPair(KeyTargetInflationRate, &p.TargetInflationRate, validateInflationRate),
	}
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
