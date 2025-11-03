package hub

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"

	"github.com/dymensionxyz/sdk-utils/utils/uevent"

	"github.com/dymensionxyz/dymension-rdk/x/hub/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

var (
	_ porttypes.ICS4Wrapper = &DecimalConversionMiddleware{}
	_ porttypes.IBCModule   = &DecimalConversionMiddleware{}
)

// DecimalConversionMiddleware implements the ICS26 callbacks for decimal conversion middleware
type DecimalConversionMiddleware struct {
	porttypes.ICS4Wrapper
	porttypes.IBCModule

	transfer  porttypes.IBCModule
	hubKeeper keeper.Keeper
}

// NewDecimalConversionMiddleware creates a new DecimalConversionMiddleware for decimal conversion
func NewDecimalConversionMiddleware(
	transfer porttypes.IBCModule,
	completeStack porttypes.IBCModule,
	ics4Wrapper porttypes.ICS4Wrapper,
	hubKeeper keeper.Keeper,
) DecimalConversionMiddleware {
	return DecimalConversionMiddleware{
		IBCModule:   completeStack,
		transfer:    transfer,
		ICS4Wrapper: ics4Wrapper,
		hubKeeper:   hubKeeper,
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
	pair, found, err := m.hubKeeper.GetConversionForDenom(ctx, packetData.Denom)
	if err != nil {
		return uevent.NewErrorAcknowledgement(ctx, errorsmod.Wrapf(err, "get decimal conversion pair"))
	}

	// No conversion needed, continue with the complete stack
	if !found {
		return m.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	// First, let the underlying transfer module handle the packet
	// This will mint the original tokens to the receiver
	ack := m.transfer.OnRecvPacket(ctx, packet, relayer)

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

	// Convert the amount
	coin := sdk.NewCoin(packetData.Denom, amount)
	convertedCoin, err := types.ConvertCoin(coin, pair)
	if err != nil {
		return uevent.NewErrorAcknowledgement(ctx, errorsmod.Wrapf(err, "convert coin"))
	}

	// Burn the original tokens from the receiver
	if err := m.hubKeeper.BurnCoins(ctx, receiver, coin); err != nil {
		return uevent.NewErrorAcknowledgement(ctx, errorsmod.Wrapf(err, "burn original coins from receiver"))
	}

	// Mint the converted tokens to the receiver
	if err := m.hubKeeper.MintCoins(ctx, receiver, convertedCoin); err != nil {
		return uevent.NewErrorAcknowledgement(ctx, errorsmod.Wrapf(err, "mint converted coins to receiver"))
	}

	return ack
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im DecimalConversionMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	// For acknowledgements, we don't need to convert anything
	// The conversion was already done in SendPacket
	return im.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

// OnTimeoutPacket implements the IBCModule interface
func (im DecimalConversionMiddleware) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	// For timeouts, we need to handle the refund with the original (pre-conversion) amount
	// The underlying transfer module will handle the refund correctly
	return im.IBCModule.OnTimeoutPacket(ctx, packet, relayer)
}

// FIXME: we have issue here, as doing IBC transfer of the converted token assume it's source denom, and thus it won't be burned
// we probably need to override the transfer keeper, and convert the coin before sending the packet

// SendPacket wraps IBC ChannelKeeper's SendPacket function to convert token amounts
// Note: The transfer module has already moved tokens from sender to escrow before calling this
func (m *DecimalConversionMiddleware) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	destinationPort string, destinationChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	packet := new(transfertypes.FungibleTokenPacketData)
	if err = transfertypes.ModuleCdc.UnmarshalJSON(data, packet); err != nil {
		return 0, errorsmod.Wrapf(errortypes.ErrJSONUnmarshal, "unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	// Check if there's a decimal conversion pair for this denom
	pair, found, err := m.hubKeeper.GetConversionForDenom(ctx, packet.Denom)
	if err != nil {
		return 0, errorsmod.Wrapf(err, "get decimal conversion pair")
	}

	if !found {
		// No conversion needed, pass through
		return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
	}

	// Parse the amount
	amount, ok := sdk.NewIntFromString(packet.Amount)
	if !ok {
		return 0, errorsmod.Wrapf(errortypes.ErrInvalidRequest, "invalid amount: %s", packet.Amount)
	}

	// Convert the amount
	coin := sdk.NewCoin(packet.Denom, amount)
	convertedCoin, err := types.ConvertCoin(coin, pair)
	if err != nil {
		return 0, errorsmod.Wrapf(err, "convert coin")
	}

	// FIXME: what happens to the truncated amount?

	// Update packet data with converted values
	packet.Denom = convertedCoin.Denom
	packet.Amount = convertedCoin.Amount.String()

	// Marshal the updated packet
	data, err = transfertypes.ModuleCdc.MarshalJSON(packet)
	if err != nil {
		return 0, errorsmod.Wrapf(errors.Join(errortypes.ErrJSONMarshal, err), "marshal ICS-20 transfer packet data")
	}

	return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
}
