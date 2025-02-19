package types

import "fmt"

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:      DefaultParams(),
		Gauges:      nil,
		LastGaugeId: 0,
	}
}

func (g GenesisState) Validate() error {
	err := g.Params.Validate()
	if err != nil {
		return fmt.Errorf("validate params: %w", err)
	}

	for _, gauge := range g.Gauges {
		err = gauge.ValidateBasic()
		if err != nil {
			return fmt.Errorf("validate gauge: %w", err)
		}
	}

	return nil
}
