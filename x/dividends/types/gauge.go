package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewGauge(
	id uint64,
	address string,
	active bool,
	approvedDenoms []string,
	queryCondition QueryCondition,
	vestingDuration VestingDuration,
	vestingFrequency VestingFrequency,
) Gauge {
	return Gauge{
		Id:               id,
		Address:          address,
		Active:           active,
		ApprovedDenoms:   approvedDenoms,
		QueryCondition:   queryCondition,
		VestingDuration:  vestingDuration,
		VestingFrequency: vestingFrequency,
	}
}

func (g Gauge) GetAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(g.Address)
}

func GaugeAccountName(id uint64) string {
	return fmt.Sprintf("%s-gauge-%d", ModuleName, id)
}

// ValidateBasic performs basic validation of the Gauge fields.
func (g Gauge) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(g.Address); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	for _, denom := range g.ApprovedDenoms {
		if err := sdk.ValidateDenom(denom); err != nil {
			return fmt.Errorf("validate approved denom: %w", err)
		}
	}
	if err := g.QueryCondition.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid query condition: %w", err)
	}
	if err := g.VestingDuration.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid vesting condition: %w", err)
	}
	if g.VestingFrequency == VestingFrequency_VESTING_FREQUENCY_UNSPECIFIED {
		return fmt.Errorf("vesting frequency cannot be unspecified")
	}
	return nil
}

// ValidateBasic performs basic validation of the QueryCondition fields.
func (qc QueryCondition) ValidateBasic() error {
	switch c := qc.Condition.(type) {
	case *QueryCondition_Stakers:
		if c.Stakers == nil {
			return fmt.Errorf("stakers field should be non-nil (it may be empty)")
		}
		return nil
	default:
		return fmt.Errorf("invalid query condition type")
	}
}

// ValidateBasic performs basic validation of the VestingCondition fields.
func (vc VestingDuration) ValidateBasic() error {
	switch c := vc.Duration.(type) {
	case *VestingDuration_Perpetual:
		if c.Perpetual == nil {
			return fmt.Errorf("perpetual field should be non-nil (it may be empty)")
		}
		return nil
	case *VestingDuration_FixedTerm:
		if c.FixedTerm.NumTotal <= 0 {
			return fmt.Errorf("num_total must be positive")
		}
		if c.FixedTerm.NumDone < 0 {
			return fmt.Errorf("num_done cannot be negative")
		}
		if c.FixedTerm.NumTotal <= c.FixedTerm.NumDone {
			return fmt.Errorf("num_done must be less than num_total")
		}
		return nil
	default:
		return fmt.Errorf("invalid vesting condition type")
	}
}
