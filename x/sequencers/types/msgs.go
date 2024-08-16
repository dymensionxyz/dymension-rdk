package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

var (
	_ sdk.Msg = (*MsgCreateSequencer)(nil)
	_ sdk.Msg = (*MsgUpdateSequencer)(nil)
)

func (m *MsgCreateSequencer) ValidateBasic() error {
	if _, err := m.GetSigner(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get signer")
	}
	// TODO implement me
	panic("implement me")
}

func (m *MsgCreateSequencer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

func (m *MsgCreateSequencer) MustGetSigner() sdk.AccAddress {
	addr, err := m.GetSigner()
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *MsgCreateSequencer) GetSigner() (sdk.AccAddress, error) {
	addr, err := sdk.AccAddressFromBech32(m.Creator)
	return addr, errorsmod.Wrap(err, "acc addr from bech32")
}

func (m *MsgUpdateSequencer) ValidateBasic() error {
	if _, err := m.GetSigner(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get signer")
	}
	// TODO implement me
	panic("implement me")
}

func (m *MsgUpdateSequencer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

func (m *MsgUpdateSequencer) MustGetSigner() sdk.AccAddress {
	addr, err := m.GetSigner()
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *MsgUpdateSequencer) GetSigner() (sdk.AccAddress, error) {
	addr, err := sdk.AccAddressFromBech32(m.Creator)
	return addr, errorsmod.Wrap(err, "acc addr from bech32")
}

func (m *MsgUpdateSequencer) MustRewardAccAddr() sdk.AccAddress {
	s := m.GetPayload().GetRewardAddr()
	return sdk.MustAccAddressFromBech32(s)
}

// Validator is a convenience method - it returns a validator object which already
// has implementations of various useful methods like obtaining various type conversions
// for the public key.
func (m *KeyAndSig) Validator() stakingtypes.Validator {
	return stakingtypes.Validator{ConsensusPubkey: m.PubKey}
}

// Build - a helper used to fill the data according to protocol
func (m *MsgUpdateSequencer) Build(
	creator auth.AccountI,
	chainID string,
) {
	s := m.GetPayload().GetRewardAddr()
	return sdk.MustAccAddressFromBech32(s)
}
