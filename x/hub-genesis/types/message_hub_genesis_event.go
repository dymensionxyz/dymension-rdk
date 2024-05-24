package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgHubGenesisEvent = "hub_genesis_event"

var _ sdk.Msg = &MsgHubGenesisEvent{}

func NewMsgHubGenesisEvent(address, channelId, hubId string) *MsgHubGenesisEvent {
	return &MsgHubGenesisEvent{}
}

func (msg *MsgHubGenesisEvent) Route() string {
	return RouterKey
}

func (msg *MsgHubGenesisEvent) Type() string {
	return TypeMsgHubGenesisEvent
}

func (msg *MsgHubGenesisEvent) GetSigners() []sdk.AccAddress {
	return nil
}

func (msg *MsgHubGenesisEvent) GetSignBytes() []byte {
	return nil
}

func (msg *MsgHubGenesisEvent) ValidateBasic() error {
	return nil
}
