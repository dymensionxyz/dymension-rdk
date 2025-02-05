package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewGauge(
	id uint64,
	address string,
	queryCondition QueryCondition,
	vestingCondition VestingCondition,
	vestingFrequency VestingFrequency,
) Gauge {
	return Gauge{
		Id:               id,
		Address:          address,
		QueryCondition:   queryCondition,
		VestingCondition: vestingCondition,
		VestingFrequency: vestingFrequency,
	}
}

// ValidateBasic performs basic validation of the Gauge fields.
func (g Gauge) ValidateBasic() error {
	if g.Id == 0 {
		return fmt.Errorf("gauge id cannot be zero")
	}
	if _, err := sdk.AccAddressFromBech32(g.Address); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	if err := g.QueryCondition.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid query condition: %w", err)
	}
	if err := g.VestingCondition.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid vesting condition: %w", err)
	}
	if g.VestingFrequency == VestingFrequency_VESTING_FREQUENCY_UNSPECIFIED {
		return fmt.Errorf("vesting frequency cannot be unspecified")
	}
	return nil
}

// ValidateBasic performs basic validation of the QueryCondition fields.
func (qc QueryCondition) ValidateBasic() error {
	switch qc.Condition.(type) {
	case *QueryCondition_Stakers:
		return nil
	default:
		return fmt.Errorf("invalid query condition type")
	}
}

// ValidateBasic performs basic validation of the VestingCondition fields.
func (vc VestingCondition) ValidateBasic() error {
	switch c := vc.Condition.(type) {
	case *VestingCondition_Perpetual:
		return nil
	case *VestingCondition_Limited:
		if c.Limited.NumUnits < 0 {
			return fmt.Errorf("num_units cannot be negative")
		}
		if c.Limited.FilledUnits <= 0 {
			return fmt.Errorf("filled_units must be greater than zero")
		}
		return nil
	default:
		return fmt.Errorf("invalid vesting condition type")
	}
}
