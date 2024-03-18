package types

import (
	errorsmod "cosmossdk.io/errors"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateDenomMetadata{}
	_ sdk.Msg = &MsgUpdateDenomMetadata{}
)

const (
	TypeMsgCreateDenomMetadata = "create_denom_metadata"
	TypeMsgUpdateDenomMetadata = "update_denom_metadata"
)

// NewMsgCreateDenomMetadata creates new instance of MsgCreateDenomMetadata
func NewMsgCreateDenomMetadata(
	sender sdk.Address,
	tokenMetadata banktypes.Metadata,
) *MsgCreateDenomMetadata {
	return &MsgCreateDenomMetadata{
		SenderAddress: sender.String(),
		TokenMetadata: tokenMetadata,
	}
}

// Route returns the name of the module
func (msg MsgCreateDenomMetadata) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgCreateDenomMetadata) Type() string { return TypeMsgCreateDenomMetadata }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateDenomMetadata) ValidateBasic() error {

	// this also checks for empty addresses
	if _, err := sdk.AccAddressFromBech32(msg.SenderAddress); err != nil {
		return errorsmod.Wrapf(err, "invalid sender address: %s", err.Error())
	}

	err := msg.TokenMetadata.Validate()
	if err != nil {
		return err
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCreateDenomMetadata) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateDenomMetadata) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}

// NewMsgCreateDenomMetadata creates new instance of MsgCreateDenomMetadata
func NewMsgUpdateDenomMetadata(
	sender sdk.Address,
	tokenMetadata banktypes.Metadata,
) *MsgUpdateDenomMetadata {
	return &MsgUpdateDenomMetadata{
		SenderAddress: sender.String(),
		TokenMetadata: tokenMetadata,
	}
}

// Route returns the name of the module
func (msg MsgUpdateDenomMetadata) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUpdateDenomMetadata) Type() string { return TypeMsgUpdateDenomMetadata }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateDenomMetadata) ValidateBasic() error {

	// this also checks for empty addresses
	if _, err := sdk.AccAddressFromBech32(msg.SenderAddress); err != nil {
		return errorsmod.Wrapf(err, "invalid sender address: %s", err.Error())
	}

	err := msg.TokenMetadata.Validate()
	if err != nil {
		return err
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgUpdateDenomMetadata) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUpdateDenomMetadata) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.SenderAddress)
	return []sdk.AccAddress{from}
}
