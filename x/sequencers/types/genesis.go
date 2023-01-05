package types

import (
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
