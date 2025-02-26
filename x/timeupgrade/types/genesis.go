package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	prototypes "github.com/gogo/protobuf/types"
)

func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

func (gs GenesisState) ValidateGenesis() error {
	if gs.EmptyTimestamp() && gs.EmptyPlan() {
		return nil
	}
	if gs.EmptyTimestamp() {
		return gerrc.ErrInvalidArgument.Wrap("timestamp empty")
	}
	if gs.EmptyPlan() {
		return gerrc.ErrInvalidArgument.Wrap("plan empty")
	}
	if err := gs.Plan.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "plan")
	}
	_, err := prototypes.TimestampFromProto(gs.GetTimestamp())
	if err != nil {
		return errorsmod.Wrap(err, "timestamp")
	}
	return nil
}

func (gs GenesisState) EmptyPlan() bool {
	return gs.Plan == nil
}

func (gs GenesisState) EmptyTimestamp() bool {
	return gs.Timestamp == nil
}
