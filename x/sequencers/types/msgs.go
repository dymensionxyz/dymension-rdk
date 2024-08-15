package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (m *KeyAndSig) GetSdkPubKey() (cryptotypes.PubKey, error) {
	c := m.PubKey.GetCachedValue()
	pubKey, ok := c.(cryptotypes.PubKey)
	if !ok {
		return nil, errorsmod.WithType(errorsmod.Wrap(gerrc.ErrInvalidArgument, "assert cryptotypes pub key"), c)
	}
	return pubKey, nil
}

func (m *KeyAndSig) MustGetConsAddr() sdk.ConsAddress {
	addr, err := m.Validator().GetConsAddr()
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *KeyAndSig) Validator() stakingtypes.Validator {
	return stakingtypes.Validator{ConsensusPubkey: m.PubKey}
}
