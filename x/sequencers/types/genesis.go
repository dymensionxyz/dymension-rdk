package types

import (
	fmt "fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		Sequencers: []stakingtypes.Validator{},
		Exported:   false,
	}
}

func (gs GenesisState) ValidateGenesis() error {
	//TODO: Add validation when gentx for sequencers works
	// if len(data.Validators) == 0 {
	// 	return types.ErrNoSequencerOnGenesis
	// }

	// Check for duplicated index in sequencer
	sequencerIndexMap := make(map[string]bool)

	for _, elem := range gs.Sequencers {
		index := elem.OperatorAddress
		if _, ok := sequencerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for sequencer")
		}
		sequencerIndexMap[index] = true
	}

	return gs.Params.Validate()
}
