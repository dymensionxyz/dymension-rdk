package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	// GenesisTriggererWhitelist is store's key for GenesisTriggererWhitelist Params
	KeyGenesisTriggererWhitelist = []byte("GenesisTriggererWhitelist")
)

// ParamTable for hub_genesis module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(genesisTriggererWhitelist []GenesisTriggererParams) Params {
	return Params{
		GenesisTriggererWhitelist: genesisTriggererWhitelist,
	}
}

func DefaultParams() Params {
	return NewParams([]GenesisTriggererParams{})
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyGenesisTriggererWhitelist, &p.GenesisTriggererWhitelist, validateGenesisTriggererWhitelist),
	}
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	return validateGenesisTriggererWhitelist(p.GenesisTriggererWhitelist)
}

// validateGenesisTriggererWhitelist validates the GenesisTriggererWhitelist param
func validateGenesisTriggererWhitelist(v interface{}) error {
	genesisTriggererWhitelist, ok := v.([]GenesisTriggererParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// Check for duplicated index in genesis triggerer address
	rollappGenesisTriggererIndexMap := make(map[string]struct{})

	for i, item := range genesisTriggererWhitelist {
		// check Bech32 format
		if _, err := sdk.AccAddressFromBech32(item.Address); err != nil {
			return fmt.Errorf("genesisTriggererWhitelist[%d] format error: %s", i, err.Error())
		}

		// check duplicate
		if _, ok := rollappGenesisTriggererIndexMap[item.Address]; ok {
			return fmt.Errorf("duplicated genesis trigerrer address in genesisTriggererWhitelist")
		}
		rollappGenesisTriggererIndexMap[item.Address] = struct{}{}
	}

	return nil
}
