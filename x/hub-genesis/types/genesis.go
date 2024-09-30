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

// ValidateBasic performs basic validation of the genesis state.
func (g GenesisState) ValidateBasic() error {
	if err := g.Params.Validate(); err != nil {
		return err
	}
	if err := g.State.Validate(); err != nil {
		return err
	}

	if g.State.OutboundTransfersEnabled {
		return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "outbound transfers should be disabled in genesis")
	}

	if g.State.HubPortAndChannel != nil {
		return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "hub port and channel should not be set in genesis")
	}

	return nil
}
