package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(params Params, state State) *GenesisState {
	return &GenesisState{
		Params: params,
		State:  state,
	}
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), State{})
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	if err := data.State.Validate(); err != nil {
		return err
	}
	nSeqsExpected := data.State.NumUnackedTransfers
	nSeqsHave := uint64(len(data.UnackedTransferSeqNums))
	if nSeqsExpected != nSeqsHave {
		return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "got different number of unacked transfers than expected: expect: %d, actual: %d", nSeqsExpected, nSeqsHave)
	}

	return nil
}
