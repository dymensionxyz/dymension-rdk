package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgHubGenesisEvent = "hub_genesis_event"

var _ sdk.Msg = &MsgHubGenesisEvent{}

func NewMsgHubGenesisEvent(address, channelId, hubId string) *MsgHubGenesisEvent {
	return &MsgHubGenesisEvent{
		Address:   address,
		ChannelId: channelId,
		HubId:     hubId,
	}
}

func (msg *MsgHubGenesisEvent) Route() string {
	return RouterKey
}

func (msg *MsgHubGenesisEvent) Type() string {
	return TypeMsgHubGenesisEvent
}

func (msg *MsgHubGenesisEvent) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{address}
}

func (msg *MsgHubGenesisEvent) GetSignBytes() []byte {
	bz := moduleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgHubGenesisEvent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}

	if msg.ChannelId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "channel id cannot be empty")
	}

	if msg.HubId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "hub id cannot be empty")
	}

	return nil
}
