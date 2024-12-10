package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	hubgenesistypes "github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// IsIBCRelayerMsg checks if all the messages in the transaction are IBC relayer messages
func IsIBCRelayerMsg(msgs []sdk.Msg) bool {
	return len(msgs) > 0 && countIBCMsgs(msgs) == len(msgs)
}

func countIBCMsgs(msgs []sdk.Msg) int {
	count := 0
	for _, msg := range msgs {
		switch msg.(type) {
		// IBC Client Messages
		case *clienttypes.MsgCreateClient, *clienttypes.MsgUpdateClient,
			*clienttypes.MsgUpgradeClient, *clienttypes.MsgSubmitMisbehaviour:
			count++

		// IBC Connection Messages
		case *conntypes.MsgConnectionOpenInit, *conntypes.MsgConnectionOpenTry,
			*conntypes.MsgConnectionOpenAck, *conntypes.MsgConnectionOpenConfirm:
			count++

		// IBC Channel Messages
		case *channeltypes.MsgChannelOpenInit, *channeltypes.MsgChannelOpenTry,
			*channeltypes.MsgChannelOpenAck, *channeltypes.MsgChannelOpenConfirm,
			*channeltypes.MsgChannelCloseInit, *channeltypes.MsgChannelCloseConfirm:
			count++

		// IBC Packet Messages
		case *channeltypes.MsgRecvPacket, *channeltypes.MsgAcknowledgement,
			*channeltypes.MsgTimeout, *channeltypes.MsgTimeoutOnClose:
			count++

		// Not strictly an IBC message, but rather a custom message for dymension
		case *hubgenesistypes.MsgSendTransfer:
			count++
		}
	}
	return count
}
