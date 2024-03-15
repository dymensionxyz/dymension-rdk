package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	// GenesisTriggerrerWhitelist is store's key for GenesisTriggerrerWhitelist Params
	KeyGenesisTriggerrerWhitelist = []byte("GenesisTriggerrerWhitelist")
)

// ParamTable for hub_genesis module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(genesisTriggererWhitelist []GenesisTriggerrerParams) Params {
	return Params{
		GenesisTriggerrerWhitelist: genesisTriggererWhitelist,
	}
}

func DefaultParams() Params {
	return NewParams([]GenesisTriggerrerParams{})
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyGenesisTriggerrerWhitelist, &p.GenesisTriggerrerWhitelist, validateGenesisTriggerrerWhitelist),
	}
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	return validateGenesisTriggerrerWhitelist(p.GenesisTriggerrerWhitelist)
}

// validateGenesisTriggerrerWhitelist validates the GenesisTriggerrerWhitelist param
func validateGenesisTriggerrerWhitelist(v interface{}) error {
	genesisTriggererWhitelist, ok := v.([]GenesisTriggerrerParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// Check for duplicated index in genesis triggerrer address
	rollappGenesisTriggerrerIndexMap := make(map[string]struct{})

	for i, item := range genesisTriggererWhitelist {
		// check Bech32 format
		if _, err := sdk.AccAddressFromBech32(item.Address); err != nil {
			return fmt.Errorf("genesisTriggererWhitelist[%d] format error: %s", i, err.Error())
		}

		// check duplicate
		if _, ok := rollappGenesisTriggerrerIndexMap[item.Address]; ok {
			return fmt.Errorf("duplicated genesis trigerrer address in genesisTriggererWhitelist")
		}
		rollappGenesisTriggerrerIndexMap[item.Address] = struct{}{}
	}

	return nil
}
