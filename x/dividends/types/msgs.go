package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (msg MsgCreateGauge) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	if err := msg.QueryCondition.ValidateBasic(); err != nil {
		return err
	}
	if err := msg.VestingCondition.ValidateBasic(); err != nil {
		return err
	}
	if msg.VestingFrequency == VestingFrequency_VESTING_FREQUENCY_UNSPECIFIED {
		return sdkerrors.ErrInvalidRequest.Wrap("vesting frequency cannot be zero")
	}
	return nil
}

func (msg MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	if err := msg.NewParams.Validate(); err != nil {
		return err
	}
	return nil
}
