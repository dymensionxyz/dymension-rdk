package types

import (
	"encoding/json"
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func (gs GenesisState) ValidateGenesis() error {
	err := gs.Params.Validate()
	if err != nil {
		return err
	}
	for _, s := range gs.GetSequencers() {
		if s.Validator == nil {
			return errorsmod.Wrap(gerrc.ErrInvalidArgument, "validator is nil")
		}
		if _, err := s.Validator.ConsPubKey(); err != nil {
			return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "cons pub key")
		}
		if _, err := sdk.ValAddressFromBech32(s.Validator.OperatorAddress); err != nil {
			return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "operator addr")
		}
		if s.RewardAddr != "" {
			if _, err := s.RewardAcc(); err != nil {
				return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "reward acc")
			}
		}
	}
	return nil
}

// MustClone returns a deep copy - intended for tests
func (gs GenesisState) MustClone() GenesisState {
	bz, err := json.Marshal(gs)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bz, &gs)
	if err != nil {
		panic(err)
	}
	return gs
}

// RewardAcc will try to parse an acc address from the sequencer reward addr assuming it is not empty string
func (s Sequencer) RewardAcc() (sdk.AccAddress, error) {
	return sdk.AccAddressFromBech32(s.GetRewardAddr())
}

func (s Sequencer) MustRewardAcc() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(s.GetRewardAddr())
}
