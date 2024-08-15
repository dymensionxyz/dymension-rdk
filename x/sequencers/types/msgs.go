package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
