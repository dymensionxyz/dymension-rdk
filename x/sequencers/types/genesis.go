package types

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *stakingtypes.GenesisState {
	// return &stakingtypes.GenesisState{
	// 	// this line is used by starport scaffolding # genesis/types/default
	// 	Params: DefaultParams(),
	// }

	return stakingtypes.DefaultGenesisState()
}
