package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.Msg = (*MsgCreateSequencer)(nil)
	_ sdk.Msg = (*MsgUpdateSequencer)(nil)
)

func (m *MsgCreateSequencer) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *MsgCreateSequencer) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgUpdateSequencer) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *MsgUpdateSequencer) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
