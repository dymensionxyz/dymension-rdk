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
	KeyMintDenom                            = []byte("MintDenom")
	KeyGenesisEpochProvisions               = []byte("GenesisEpochProvisions")
	KeyEpochIdentifier                      = []byte("EpochIdentifier")
	KeyReductionPeriodInEpochs              = []byte("ReductionPeriodInEpochs")
	KeyReductionFactor                      = []byte("ReductionFactor")
	KeyMintingRewardsDistributionStartEpoch = []byte("MintingRewardsDistributionStartEpoch")
)

// ParamTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintDenom string, genesisEpochProvisions sdk.Dec, epochIdentifier string,
	ReductionFactor sdk.Dec, reductionPeriodInEpochs int64, mintingRewardsDistributionStartEpoch int64,
) Params {
	return Params{
		MintDenom:                            mintDenom,
		GenesisEpochProvisions:               genesisEpochProvisions,
		EpochIdentifier:                      epochIdentifier,
		ReductionPeriodInEpochs:              reductionPeriodInEpochs,
		ReductionFactor:                      ReductionFactor,
		MintingRewardsDistributionStartEpoch: mintingRewardsDistributionStartEpoch,
	}
}

// minting params
func DefaultParams() Params {
	return Params{
		MintDenom:                            sdk.DefaultBondDenom,
		GenesisEpochProvisions:               sdk.NewDec(2_500_000).Mul(sdk.NewDec(1_000_000)).Quo(sdk.NewDec(24 * 365)), // 2.5MST first year, broken into hours ~= 285ST / hour
		EpochIdentifier:                      "day",                                                                      // 1 hour
		ReductionPeriodInEpochs:              24 * 365,                                                                   // 24hrs*365d = 8760
		ReductionFactor:                      sdk.NewDec(1).QuoInt64(2),
		MintingRewardsDistributionStartEpoch: 0,
	}
}

// validate params.
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateGenesisEpochProvisions(p.GenesisEpochProvisions); err != nil {
		return err
	}
	if err := epochtypes.ValidateEpochIdentifierInterface(p.EpochIdentifier); err != nil {
		return err
	}
	if err := validateReductionPeriodInEpochs(p.ReductionPeriodInEpochs); err != nil {
		return err
	}
	if err := validateReductionFactor(p.ReductionFactor); err != nil {
		return err
	}
	if err := validateMintingRewardsDistributionStartEpoch(p.MintingRewardsDistributionStartEpoch); err != nil {
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
		paramtypes.NewParamSetPair(KeyGenesisEpochProvisions, &p.GenesisEpochProvisions, validateGenesisEpochProvisions),
		paramtypes.NewParamSetPair(KeyEpochIdentifier, &p.EpochIdentifier, epochtypes.ValidateEpochIdentifierInterface),
		paramtypes.NewParamSetPair(KeyReductionPeriodInEpochs, &p.ReductionPeriodInEpochs, validateReductionPeriodInEpochs),
		paramtypes.NewParamSetPair(KeyReductionFactor, &p.ReductionFactor, validateReductionFactor),
		paramtypes.NewParamSetPair(KeyMintingRewardsDistributionStartEpoch, &p.MintingRewardsDistributionStartEpoch, validateMintingRewardsDistributionStartEpoch),
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

func validateGenesisEpochProvisions(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("genesis epoch provision must be non-negative")
	}

	return nil
}

func validateReductionPeriodInEpochs(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}

	return nil
}

func validateReductionFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDec(1)) {
		return fmt.Errorf("reduction factor cannot be greater than 1")
	}

	if v.IsNegative() {
		return fmt.Errorf("reduction factor cannot be negative")
	}

	return nil
}

func validateMintingRewardsDistributionStartEpoch(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("start epoch must be non-negative")
	}

	return nil
}
