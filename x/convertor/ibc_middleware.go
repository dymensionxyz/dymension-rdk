package convertor

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"

	"github.com/dymensionxyz/dymension-rdk/x/convertor/keeper"
	"github.com/dymensionxyz/sdk-utils/utils/uevent"
)

var (
	_ porttypes.IBCModule = &DecimalConversionMiddleware{}
)

// DecimalConversionMiddleware implements the ICS26 callbacks for decimal conversion middleware
type DecimalConversionMiddleware struct {
	porttypes.IBCModule
	convertor keeper.Keeper
}

// NewIBCModule creates a new IBCModule for the hub module with decimal conversion middleware
// next: the next middleware in the stack (or the complete stack so far)
func NewDecimalConversionMiddleware(
	next porttypes.IBCModule,
	convertor keeper.Keeper,
) DecimalConversionMiddleware {
	return DecimalConversionMiddleware{
		IBCModule: next,
		convertor: convertor,
	}
}

// OnRecvPacket handles incoming packets. It first lets the underlying transfer module
// process the packet (which mints tokens to receiver), then performs decimal conversion
// by burning the original tokens and minting the converted tokens to the receiver.
func (m DecimalConversionMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {

	// Parse packet data to check if conversion is needed
	packetData := new(transfertypes.FungibleTokenPacketData)
	transfertypes.ModuleCdc.MustUnmarshalJSON(packet.GetData(), packetData)

	// Check if there's a decimal conversion pair for this denom
	required, err := m.convertor.ConversionRequired(ctx, packetData.Denom)
	if err != nil {
		return uevent.NewErrorAcknowledgement(ctx, errorsmod.Wrapf(err, "get decimal conversion pair"))
	}

	// No conversion needed, continue with the complete stack
	if !required {
		return m.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	// First, let the underlying transfer module handle the packet
	// This will mint the original tokens to the receiver
	ack := m.IBCModule.OnRecvPacket(ctx, packet, relayer)

	if !ack.Success() {
		// If the underlying transfer failed, don't attempt conversion
		return ack
	}

	// Parse the receiver address
	receiver, err := sdk.AccAddressFromBech32(packetData.Receiver)
	if err != nil {
		return uevent.NewErrorAcknowledgement(ctx, errorsmod.Wrapf(err, "invalid receiver address"))
	}

	// Parse the amount
	amount, ok := sdk.NewIntFromString(packetData.Amount)
	if !ok {
		return uevent.NewErrorAcknowledgement(ctx, errorsmod.Wrapf(errortypes.ErrInvalidRequest, "invalid amount: %s", packetData.Amount))
	}

	// Convert the amount from bridge token (custom decimals) to rollapp token (18 decimals)
	convertedAmt, err := m.convertor.ConvertFromBridgeAmt(ctx, amount)
	if err != nil {
		return uevent.NewErrorAcknowledgement(ctx, errorsmod.Wrapf(err, "convert amount from bridge token"))
	}

	// Calculate the delta to mint (difference between converted and original amount)
	// The delta represents additional tokens needed to reach the full 18-decimal precision
	delta := sdk.NewCoin(packetData.Denom, convertedAmt.Sub(amount))

	// Mint the missing amount of tokens to the receiver
	if err := m.convertor.MintCoins(ctx, receiver, delta); err != nil {
		return uevent.NewErrorAcknowledgement(ctx, errorsmod.Wrapf(err, "mint converted coins to receiver"))
	}

	return ack
}

// OnAcknowledgementPacket implements the IBCModule interface
func (m DecimalConversionMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	// The conversion was already done in SendPacket
	err := m.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	if err != nil {
		return err
	}

	// parse the acknowledgement, if it's error, it means there's a refund, and we need to convert token's decimals
	var ack channeltypes.Acknowledgement
	if err := transfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(errortypes.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}

	// no error, nothing to do
	if ack.GetError() == "" {
		return nil
	}

	// parse the packet data
	var packetData transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &packetData); err != nil {
		return errorsmod.Wrapf(errortypes.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %v", err)
	}

	// check if there's a decimal conversion pair for this denom
	required, err := m.convertor.ConversionRequired(ctx, packetData.Denom)
	if err != nil {
		return errorsmod.Wrapf(err, "get decimal conversion pair")
	}

	// no conversion needed, nothing to do
	if !required {
		return nil
	}

	sender, err := sdk.AccAddressFromBech32(packetData.Sender)
	if err != nil {
		return errorsmod.Wrapf(err, "invalid sender address")
	}

	// Parse the amount
	amount, ok := sdk.NewIntFromString(packetData.Amount)
	if !ok {
		return errorsmod.Wrapf(errortypes.ErrInvalidRequest, "invalid amount: %s", packetData.Amount)
	}
	// convert the amount from bridge token (custom decimals) to rollapp token (18 decimals)
	convertedAmt, err := m.convertor.ConvertFromBridgeAmt(ctx, amount)
	if err != nil {
		return err
	}

	// On refund, user received back 'amount' but originally sent 'convertedAmt'
	// So we need to mint the difference back to them
	delta := sdk.NewCoin(packetData.Denom, convertedAmt.Sub(amount))

	err = m.convertor.MintCoins(ctx, sender, delta)
	if err != nil {
		return errorsmod.Wrapf(err, "mint converted coins to sender")
	}

	return nil
}

// OnTimeoutPacket implements the IBCModule interface
func (m DecimalConversionMiddleware) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	// The underlying transfer module will handle the refund with the converted (rollapp) amount
	err := m.IBCModule.OnTimeoutPacket(ctx, packet, relayer)
	if err != nil {
		return err
	}

	// parse the packet data
	var packetData transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &packetData); err != nil {
		return errorsmod.Wrapf(errortypes.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %v", err)
	}

	// check if there's a decimal conversion pair for this denom
	required, err := m.convertor.ConversionRequired(ctx, packetData.Denom)
	if err != nil {
		return errorsmod.Wrapf(err, "get decimal conversion pair")
	}

	// no conversion needed, nothing to do
	if !required {
		return nil
	}

	sender, err := sdk.AccAddressFromBech32(packetData.Sender)
	if err != nil {
		return errorsmod.Wrapf(err, "invalid sender address")
	}

	// Parse the amount
	amount, ok := sdk.NewIntFromString(packetData.Amount)
	if !ok {
		return errorsmod.Wrapf(errortypes.ErrInvalidRequest, "invalid amount: %s", packetData.Amount)
	}

	// convert the amount from bridge token (custom decimals) to rollapp token (18 decimals)
	convertedAmt, err := m.convertor.ConvertFromBridgeAmt(ctx, amount)
	if err != nil {
		return err
	}

	// On timeout, user received back 'amount' but originally sent 'convertedAmt'
	// So we need to mint the difference back to them
	delta := sdk.NewCoin(packetData.Denom, convertedAmt.Sub(amount))

	err = m.convertor.MintCoins(ctx, sender, delta)
	if err != nil {
		return errorsmod.Wrapf(err, "mint converted coins to sender")
	}

	return nil
}
