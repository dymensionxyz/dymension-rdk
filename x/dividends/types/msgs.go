package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateGauge)(nil)
	_ sdk.Msg = (*MsgUpdateParams)(nil)
)

func (m MsgCreateGauge) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	if err := m.QueryCondition.ValidateBasic(); err != nil {
		return errors.Join(sdkerrors.ErrInvalidRequest, err)
	}
	if err := m.VestingCondition.ValidateBasic(); err != nil {
		return errors.Join(sdkerrors.ErrInvalidRequest, err)
	}
	if m.VestingFrequency == VestingFrequency_VESTING_FREQUENCY_UNSPECIFIED {
		return sdkerrors.ErrInvalidRequest.Wrap("vesting frequency cannot be zero")
	}
	return nil
}

func (m MsgCreateGauge) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

func (m MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	if err := m.NewParams.Validate(); err != nil {
		return err
	}
	return nil
}

func (m MsgUpdateParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}
