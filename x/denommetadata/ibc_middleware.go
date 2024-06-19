package denommetadata

import (
	. "slices"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"

	"github.com/dymensionxyz/dymension-rdk/utils"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

type IBCSendMiddleware struct {
	porttypes.ICS4Wrapper

	hubKeeper  types.HubKeeper
	bankKeeper types.BankKeeper
}

// NewIBCSendMiddleware creates a new ICS4Wrapper.
// It intercepts outgoing IBC packets and adds token metadata to the memo if the hub doesn't have it.
// This is a solution for adding token metadata to fungible tokens transferred over IBC,
// targeted at hubs that don't have the token metadata for the token being transferred.
// More info here: https://www.notion.so/dymension/ADR-x-IBC-Denom-Metadata-Transfer-From-Rollapp-to-Hub-54e74e50adeb4d77b1f8777c05a73390?pvs=4
func NewIBCSendMiddleware(
	ics porttypes.ICS4Wrapper,
	hubKeeper types.HubKeeper,
	bankKeeper types.BankKeeper,
) *IBCSendMiddleware {
	return &IBCSendMiddleware{
		ICS4Wrapper: ics,
		hubKeeper:   hubKeeper,
		bankKeeper:  bankKeeper,
	}
}

// SendPacket wraps IBC ChannelKeeper's SendPacket function
func (m *IBCSendMiddleware) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	destinationPort string, destinationChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	packet := new(transfertypes.FungibleTokenPacketData)
	if err = types.ModuleCdc.UnmarshalJSON(data, packet); err != nil {
		return 0, errorsmod.Wrapf(errortypes.ErrJSONUnmarshal, "unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	if types.MemoAlreadyHasPacketMetadata(packet.Memo) {
		return 0, types.ErrMemoDenomMetadataAlreadyExists
	}

	hub, err := m.hubKeeper.ExtractHubFromChannel(ctx, destinationPort, destinationChannel)
	if err != nil {
		return 0, errorsmod.Wrapf(errortypes.ErrInvalidRequest, "extract hub from channel: %s", err.Error())
	}

	if hub == nil {
		return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
	}

	if transfertypes.ReceiverChainIsSource(destinationPort, destinationChannel, packet.Denom) {
		return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
	}

	// Check if the hub already contains the denom metadata by matching the base of the denom metadata.
	// At the first match, we assume that the hub already contains the metadata.
	// It would be technically possible to have a race condition where the denom metadata is added to the hub
	// from another packet before this packet is acknowledged.
	if Contains(hub.RegisteredDenoms, packet.Denom) {
		return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
	}

	// get the denom metadata from the bank keeper, if it doesn't exist, move on to the next middleware in the chain
	denomMetadata, ok := m.bankKeeper.GetDenomMetaData(ctx, packet.Denom)
	if !ok {
		return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
	}

	packet.Memo, err = types.AddDenomMetadataToMemo(packet.Memo, denomMetadata)
	if err != nil {
		return 0, errorsmod.Wrapf(errortypes.ErrUnauthorized, "add denom metadata to memo: %s", err.Error())
	}

	data, err = types.ModuleCdc.MarshalJSON(packet)
	if err != nil {
		return 0, errorsmod.Wrapf(errortypes.ErrJSONMarshal, "marshal ICS-20 transfer packet data: %s", err.Error())
	}

	return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
}

var _ porttypes.IBCModule = &IBCRecvMiddleware{}

// IBCRecvMiddleware implements the ICS26 callbacks for the transfer middleware
type IBCRecvMiddleware struct {
	porttypes.IBCModule
	bankKeeper     types.BankKeeper
	transferKeeper types.TransferKeeper
	hubKeeper      types.HubKeeper
	hooks          types.MultiDenomMetadataHooks
}

// NewIBCRecvMiddleware creates a new IBCRecvMiddleware given the keeper and underlying application
func NewIBCRecvMiddleware(
	app porttypes.IBCModule,
	bankKeeper types.BankKeeper,
	transferKeeper types.TransferKeeper,
	hubKeeper types.HubKeeper,
	hooks types.MultiDenomMetadataHooks,
) IBCRecvMiddleware {
	return IBCRecvMiddleware{
		IBCModule:      app,
		bankKeeper:     bankKeeper,
		transferKeeper: transferKeeper,
		hubKeeper:      hubKeeper,
		hooks:          hooks,
	}
}

// OnAcknowledgementPacket adds the token metadata to the hub if it doesn't exist
func (im IBCRecvMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(errortypes.ErrJSONUnmarshal, "unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}

	if !ack.Success() {
		return im.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	var data transfertypes.FungibleTokenPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return errorsmod.Wrapf(errortypes.ErrJSONUnmarshal, "unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	dm := types.ParsePacketMetadata(data.Memo)
	if dm == nil {
		return im.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	hub, err := im.hubKeeper.ExtractHubFromChannel(ctx, packet.SourcePort, packet.SourceChannel)
	if err != nil {
		return errorsmod.Wrapf(errortypes.ErrInvalidRequest, "extract hub from channel: %s", err.Error())
	}
	if hub == nil {
		return errorsmod.Wrapf(errortypes.ErrNotFound, "hub not found")
	}

	if !Contains(hub.RegisteredDenoms, dm.Base) {
		// add the new token denom base to the list of hub's registered denoms
		hub.RegisteredDenoms = append(hub.RegisteredDenoms, dm.Base)

		im.hubKeeper.SetHub(ctx, *hub)
	}

	return im.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

// OnRecvPacket registers the denom metadata if it does not exist.
// It will intercept an incoming packet and check if the denom metadata exists.
// If it does not, it will register the denom metadata.
// The handler will expect a 'denom_metadata' object in the memo field of the packet.
// If the memo is not an object, or does not contain the metadata, it moves on to the next handler.
func (im IBCRecvMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	packetData := new(transfertypes.FungibleTokenPacketData)
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), packetData); err != nil {
		err = errorsmod.Wrapf(errortypes.ErrJSONUnmarshal, "unmarshal ICS-20 transfer packet data")
		return channeltypes.NewErrorAcknowledgement(err)
	}

	if packetData.Memo == "" {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	// at this point it's safe to assume that we are not handling a native token of the hub
	denomTrace := utils.GetForeignDenomTrace(packet.GetDestChannel(), packetData.Denom)
	ibcDenom := denomTrace.IBCDenom()

	if _, exist := im.bankKeeper.GetDenomMetaData(ctx, ibcDenom); exist {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	dm := types.ParsePacketMetadata(packetData.Memo)
	if dm == nil {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	if err := dm.Validate(); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	dm.Base = ibcDenom
	dm.DenomUnits[0].Denom = dm.Base

	im.bankKeeper.SetDenomMetaData(ctx, *dm)
	// set hook after denom metadata creation
	if err := im.hooks.AfterDenomMetadataCreation(ctx, *dm); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	if !im.transferKeeper.HasDenomTrace(ctx, denomTrace.Hash()) {
		im.transferKeeper.SetDenomTrace(ctx, denomTrace)
	}

	return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
}
