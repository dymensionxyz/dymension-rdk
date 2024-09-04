package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _, _ sdk.Msg = &MsgSoftwareUpgrade{}, &MsgCancelUpgrade{}

func (m *MsgSoftwareUpgrade) ValidateBasic() error {
	return m.OriginalUpgrade.ValidateBasic()
}

func (m *MsgSoftwareUpgrade) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.GetOriginalUpgrade().Authority)
	return []sdk.AccAddress{addr}
}

func (m *MsgCancelUpgrade) ValidateBasic() error {
	//TODO implement me
	panic("implement me")
}

func (m *MsgCancelUpgrade) GetSigners() []sdk.AccAddress {
	//TODO implement me
	panic("implement me")
}
