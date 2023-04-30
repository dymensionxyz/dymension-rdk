package erc20middleware

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"

	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// var _ porttypes.Middleware = &AppModule{}

// // SendPacket implements types.Middleware
// func (*IBCMiddleware) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI) error {
// 	panic("unimplemented")
// }

// // WriteAcknowledgement implements types.Middleware
// func (*IBCMiddleware) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
// 	panic("unimplemented")
// }

// OnRecvPacket implements the IBCModule interface.
// It receives the tokens through the default ICS20 OnRecvPacket callback logic
// and then automatically converts the Cosmos Coin to their ERC20 token
// representation.
// If the acknowledgement fails, this callback will default to the ibc-core
// packet callback.
func (im AppModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	ack := im.IBCModule.OnRecvPacket(ctx, packet, relayer)

	// return if the acknowledgement is an error ACK
	if !ack.Success() {
		return ack
	}

	return im.keeper.OnRecvPacket(ctx, packet, ack)
}

// OnAcknowledgementPacket implements the IBCModule interface.
// It refunds the token transferred and then automatically converts the
// Cosmos Coin to their ERC20 token representation.
func (im AppModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := transfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(errortypes.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}

	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return errorsmod.Wrapf(errortypes.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	if err := im.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer); err != nil {
		return err
	}

	return im.keeper.OnAcknowledgementPacket(ctx, packet, data, ack)
}

// OnTimeoutPacket implements the IBCModule interface.
// It refunds the token transferred and then automatically converts the
// Cosmos Coin to their ERC20 token representation.
func (im AppModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return errorsmod.Wrapf(errortypes.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	if err := im.IBCModule.OnTimeoutPacket(ctx, packet, relayer); err != nil {
		return err
	}

	return im.keeper.OnTimeoutPacket(ctx, packet, data)
}
