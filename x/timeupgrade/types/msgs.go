package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
	prototypes "github.com/gogo/protobuf/types"
)

var _, _ sdk.Msg = &MsgSoftwareUpgrade{}, &MsgCancelUpgrade{}

func (m *MsgSoftwareUpgrade) ValidateBasic() error {
	if m.Drs == 0 {
		return sdkerrors.Wrapf(errors.ErrInvalidVersion, "invalid drs version")
	}
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrap(err, "authority")
	}

	_, err := prototypes.TimestampFromProto(m.UpgradeTime)
	if err != nil || m.UpgradeTime == nil || m.UpgradeTime.Seconds == 0 {
		return sdkerrors.Wrap(err, "upgrade time")
	}
	return nil
}

func (m *MsgSoftwareUpgrade) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func (m *MsgCancelUpgrade) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrap(err, "authority")
	}

	return nil
}

func (m *MsgCancelUpgrade) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}
