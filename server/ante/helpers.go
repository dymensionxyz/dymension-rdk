package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
)

// IsIBCRelayerMsg checks if all the messages in the transaction are IBC relayer messages
func IsIBCRelayerMsg(msgs []sdk.Msg) bool {
	for _, msg := range msgs {
		switch msg.(type) {
		// IBC Client Messages
		case *clienttypes.MsgCreateClient, *clienttypes.MsgUpdateClient,
			*clienttypes.MsgUpgradeClient, *clienttypes.MsgSubmitMisbehaviour:

		// IBC Connection Messages
		case *conntypes.MsgConnectionOpenInit, *conntypes.MsgConnectionOpenTry,
			*conntypes.MsgConnectionOpenAck, *conntypes.MsgConnectionOpenConfirm:

		// IBC Channel Messages
		case *channeltypes.MsgChannelOpenInit, *channeltypes.MsgChannelOpenTry,
			*channeltypes.MsgChannelOpenAck, *channeltypes.MsgChannelOpenConfirm,
			*channeltypes.MsgChannelCloseInit, *channeltypes.MsgChannelCloseConfirm:

		// IBC Packet Messages
		case *channeltypes.MsgRecvPacket, *channeltypes.MsgAcknowledgement,
			*channeltypes.MsgTimeout, *channeltypes.MsgTimeoutOnClose:

		default:
			return false
		}
	}
	return true
}
