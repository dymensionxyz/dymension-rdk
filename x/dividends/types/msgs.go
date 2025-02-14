package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateGauge)(nil)
	_ sdk.Msg = (*MsgUpdateGauge)(nil)
	_ sdk.Msg = (*MsgDeactivateGauge)(nil)
	_ sdk.Msg = (*MsgUpdateParams)(nil)
)

func (m MsgCreateGauge) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	for _, denom := range m.ApprovedDenoms {
		if err := sdk.ValidateDenom(denom); err != nil {
			return fmt.Errorf("validate approved denom: %w", err)
		}
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

func (m MsgUpdateGauge) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	for _, denom := range m.ApprovedDenoms {
		if err := sdk.ValidateDenom(denom); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrapf("validate approved denom: %s", err)
		}
	}
	return nil
}

func (m MsgUpdateGauge) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

func (m MsgDeactivateGauge) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	return nil
}

func (m MsgDeactivateGauge) GetSigners() []sdk.AccAddress {
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
