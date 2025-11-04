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

	transfer  porttypes.IBCModule // used to skip the transfer stack, and mint tokens directly
	convertor keeper.Keeper
}

// NewIBCModule creates a new IBCModule for the hub module with decimal conversion middleware
// transfer: the base transfer keeper (used to skip middleware and mint tokens directly)
// next: the next middleware in the stack (or the complete stack so far)
func NewDecimalConversionMiddleware(
	transfer porttypes.IBCModule,
	next porttypes.IBCModule,
	hubKeeper keeper.Keeper,
) DecimalConversionMiddleware {
	return DecimalConversionMiddleware{
		IBCModule: next,
		transfer:  transfer,
		convertor: hubKeeper,
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

	// Convert the amount from bridge token (custom decimals) to rollapp token (18 decimals)
	convertedAmt, err := m.convertor.ConvertFromBridgeAmt(ctx, amount)
	if err != nil {
		return uevent.NewErrorAcknowledgement(ctx, errorsmod.Wrapf(err, "convert amount from bridge token"))
	}

	delta := sdk.NewCoin(packetData.Denom, convertedAmt.Sub(amount))

	// Mint the missing amount of tokens to the receiver
	if err := m.convertor.MintCoins(ctx, receiver, delta); err != nil {
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

	// FIXME: Handle the acknowledgement packet

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
	// FIXME: Handle the timeout packet

	// For timeouts, we need to handle the refund with the original (pre-conversion) amount
	// The underlying transfer module will handle the refund correctly
	return im.IBCModule.OnTimeoutPacket(ctx, packet, relayer)
}
