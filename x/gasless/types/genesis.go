package types

import (
	"fmt"
)

// DefaultGenesis returns the default Capability genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (genState GenesisState) Validate() error {
	if err := genState.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	for _, txToTankIDs := range genState.TxToGasTankIds {
		if err := txToTankIDs.Validate(); err != nil {
			return fmt.Errorf("invalid txToTankIDs: %w", err)
		}
	}

	for _, tank := range genState.GasTanks {
		if err := tank.Validate(); err != nil {
			return fmt.Errorf("invalid tank: %w", err)
		}
	}

	for _, consumer := range genState.GasConsumers {
		if err := consumer.Validate(); err != nil {
			return fmt.Errorf("invalid consumer: %w", err)
		}
	}

	return nil
}
