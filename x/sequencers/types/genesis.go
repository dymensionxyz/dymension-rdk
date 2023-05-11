package types

import (
	fmt "fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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
	if len(gs.Sequencers) == 0 {
		return ErrNoSequencerOnGenesis
	}

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

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range g.Sequencers {
		if err := g.Sequencers[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}
