package types

import (
	fmt "fmt"

	errorsmod "cosmossdk.io/errors"

	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

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
	metadatas []DenomMetadata,
) *MsgCreateDenomMetadata {
	return &MsgCreateDenomMetadata{
		SenderAddress: sender.String(),
		Metadatas:     metadatas,
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

	for _, metadata := range msg.Metadatas {
		err := metadata.TokenMetadata.Validate()
		if err != nil {
			return err
		}

		denomTrace := transfertypes.ParseDenomTrace(metadata.DenomTrace)
		// If path is empty, then the denom is not ibc denom
		if denomTrace.Path != "" {
			denom := denomTrace.IBCDenom()
			if denom != metadata.TokenMetadata.Base {
				return fmt.Errorf("denom parse from denom trace does not match metadata base denom. base denom: %s, expected: %s", metadata.TokenMetadata.Base, denom)
			}
		}
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
	metadatas []DenomMetadata,
) *MsgUpdateDenomMetadata {
	return &MsgUpdateDenomMetadata{
		SenderAddress: sender.String(),
		Metadatas:     metadatas,
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

	for _, metadata := range msg.Metadatas {
		err := metadata.TokenMetadata.Validate()
		if err != nil {
			return err
		}

		denomTrace := transfertypes.ParseDenomTrace(metadata.DenomTrace)
		// If path is empty, then the denom is not ibc denom
		if denomTrace.Path != "" {
			denom := denomTrace.IBCDenom()
			if denom != metadata.TokenMetadata.Base {
				return fmt.Errorf("denom parse from denom trace does not match metadata base denom. base denom: %s, expected: %s", metadata.TokenMetadata.Base, denom)
			}
		}
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
