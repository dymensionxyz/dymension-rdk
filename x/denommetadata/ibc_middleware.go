package denommetadata

import (
	"errors"
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
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/dymensionxyz/sdk-utils/utils/uibc"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
	hgtypes "github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	hubtypes "github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

// ICS4Wrapper intercepts outgoing IBC packets and adds token metadata to the memo if the hub doesn't have it.
// This is a solution for adding token metadata to fungible tokens transferred over IBC,
// in case the Hub doesn't have the token metadata for the token being transferred.
// More info here: https://www.notion.so/dymension/ADR-x-IBC-Denom-Metadata-Transfer-From-Rollapp-to-Hub-54e74e50adeb4d77b1f8777c05a73390?pvs=4
type ICS4Wrapper struct {
	porttypes.ICS4Wrapper

	hubKeeper      types.HubKeeper
	bankKeeper     types.BankKeeper
	getHubGenState func(ctx sdk.Context) hgtypes.State
}

// NewICS4Wrapper creates a new ICS4Wrapper.
func NewICS4Wrapper(
	ics porttypes.ICS4Wrapper,
	hubKeeper types.HubKeeper,
	bankKeeper types.BankKeeper,
	getState func(ctx sdk.Context) hgtypes.State,
) *ICS4Wrapper {
	return &ICS4Wrapper{
		ICS4Wrapper:    ics,
		hubKeeper:      hubKeeper,
		bankKeeper:     bankKeeper,
		getHubGenState: getState,
	}
}

// SendPacket wraps IBC ChannelKeeper's SendPacket function
func (m *ICS4Wrapper) SendPacket(
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

	if hubGenState := m.getHubGenState(ctx); !hubGenState.IsCanonicalHubTransferChannel(destinationPort, destinationChannel) {
		return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
	}

	if types.MemoHasPacketMetadata(packet.Memo) {
		return 0, gerrc.ErrAlreadyExists
	}

	if transfertypes.ReceiverChainIsSource(destinationPort, destinationChannel, packet.Denom) {
		return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
	}

	state := m.hubKeeper.GetState(ctx)

	// Check if the hub already contains the denom metadata by matching the base of the denom metadata.
	// At the first match, we assume that the hub already contains the metadata.
	// It would be technically possible to have a race condition where the denom metadata is added to the hub
	// from another packet before this packet is acknowledged.
	// If the denom metadata exists but is either PENDING or ACTIVE, proceed to the next middleware in the chain.
	if ContainsFunc(state.Hub.RegisteredDenoms, func(denom *hubtypes.RegisteredDenom) bool {
		return denom.Base == packet.Denom && denom.Status != hubtypes.RegisteredDenom_INACTIVE
	}) {
		return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
	}

	denomMetadata, ok := m.bankKeeper.GetDenomMetaData(ctx, packet.Denom)
	if !ok {
		return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
	}

	packet.Memo, err = types.AddDenomMetadataToMemo(packet.Memo, denomMetadata)
	if err != nil {
		return 0, errorsmod.Wrapf(errors.Join(errortypes.ErrInvalidRequest, err), "add denom metadata to memo")
	}

	data, err = types.ModuleCdc.MarshalJSON(packet)
	if err != nil {
		return 0, errorsmod.Wrapf(errors.Join(errortypes.ErrJSONMarshal, err), "marshal ICS-20 transfer packet data")
	}

	registeredDenom := &hubtypes.RegisteredDenom{
		Base:   denomMetadata.Base,
		Status: hubtypes.RegisteredDenom_PENDING,
	}
	state.Hub.RegisteredDenoms = append(state.Hub.RegisteredDenoms, registeredDenom)
	m.hubKeeper.SetState(ctx, state)

	return m.ICS4Wrapper.SendPacket(ctx, chanCap, destinationPort, destinationChannel, timeoutHeight, timeoutTimestamp, data)
}

var _ porttypes.IBCModule = &IBCModule{}

// IBCModule implements the ICS26 callbacks for the transfer middleware
type IBCModule struct {
	porttypes.IBCModule
	bankKeeper     types.BankKeeper
	transferKeeper types.TransferKeeper
	hubKeeper      types.HubKeeper
	hooks          types.MultiDenomMetadataHooks
}

// NewIBCModule creates a new IBCModule given the keepers and underlying application
func NewIBCModule(
	app porttypes.IBCModule,
	bankKeeper types.BankKeeper,
	transferKeeper types.TransferKeeper,
	hubKeeper types.HubKeeper,
	hooks types.MultiDenomMetadataHooks,
) IBCModule {
	return IBCModule{
		IBCModule:      app,
		bankKeeper:     bankKeeper,
		transferKeeper: transferKeeper,
		hubKeeper:      hubKeeper,
		hooks:          hooks,
	}
}

// OnAcknowledgementPacket adds the token metadata to the hub if it doesn't exist
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(errors.Join(errortypes.ErrJSONUnmarshal, err), "unmarshal ICS-20 transfer acknowledgement")
	}

	var data transfertypes.FungibleTokenPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return errorsmod.Wrapf(errors.Join(errortypes.ErrJSONUnmarshal, err), "unmarshal ICS-20 transfer packet data")
	}

	dm := types.ParsePacketMetadata(data.Memo)
	if dm == nil {
		return im.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	state := im.hubKeeper.GetState(ctx)

	// find the denom from the list by matching the base
	// if it exists, update the status
	// if the ack is error or timeout, update the status to INACTIVE
	// otherwise, update the status to ACTIVE
	for i, denom := range state.Hub.RegisteredDenoms {
		if denom.Base == dm.Base {
			state.Hub.RegisteredDenoms[i].Status = map[bool]hubtypes.RegisteredDenom_Status{
				true:  hubtypes.RegisteredDenom_ACTIVE,
				false: hubtypes.RegisteredDenom_INACTIVE,
			}[ack.Success()]
			im.hubKeeper.SetState(ctx, state)
			break
		}
	}

	return im.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

// OnRecvPacket registers the denom metadata if it does not exist.
// It will intercept an incoming packet and check if the denom metadata exists.
// If it does not, it will register the denom metadata.
// The handler will expect a 'denom_metadata' object in the memo field of the packet.
// If the memo is not an object, or does not contain the metadata, it moves on to the next handler.
func (im IBCModule) OnRecvPacket(
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

	dm := types.ParsePacketMetadata(packetData.Memo)
	if dm == nil {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	if err := dm.Validate(); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	// at this point it's safe to assume that we are not handling a native token of the rollapp,
	// as the Hub, before including the denom metadata in the packet, should have checked if the receiver chain is the source.
	denomTrace := uibc.GetForeignDenomTrace(packet.GetDestChannel(), packetData.Denom)
	ibcDenom := denomTrace.IBCDenom()

	if _, ok := im.bankKeeper.GetDenomMetaData(ctx, ibcDenom); ok {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	dm.Base = ibcDenom
	dm.DenomUnits[0].Denom = dm.Base

	im.bankKeeper.SetDenomMetaData(ctx, *dm)
	// set hook after denom metadata creation
	if err := im.hooks.AfterDenomMetadataCreation(ctx, *dm); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
}
