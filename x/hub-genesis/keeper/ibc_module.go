package keeper

import (
	"encoding/json"
	"errors"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

var transferTimeout = (time.Hour) // we use short timeout as we have retry mechanism

type IBCModule struct {
	porttypes.IBCModule
	k             Keeper
	bank          types.BankKeeper
	channelKeeper types.ChannelKeeper
}

func NewIBCModule(next porttypes.IBCModule, k Keeper, bank types.BankKeeper, chanKeeper types.ChannelKeeper) *IBCModule {
	return &IBCModule{next, k, bank, chanKeeper}
}

func (w IBCModule) logger(ctx sdk.Context) log.Logger {
	return w.k.Logger(ctx)
}

// On successful OnChanOpenConfirm for the canonical channel, the genesis bridge flow will be initiated.
// It will prepare the genesis bridge data and send it over the channel.
// The genesis bridge data includes the genesis info, the native denom metadata, and the genesis transfer packet.
// Since transfers are only sent once, it does not matter if someone else tries to open
// a channel in future (it will no-op).
func (w IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	err := w.IBCModule.OnChanOpenConfirm(ctx, portID, channelID)
	if err != nil {
		return err
	}

	state := w.k.GetState(ctx)
	if state.CanonicalHubTransferChannelHasBeenSet() {
		// We only set the canonical channel in this function, so if it's already been set, we don't need
		// to send the transfers again.
		return nil
	}

	// if already ongoing, do nothing (not sure it can happen. maybe on retries?)
	if w.k.IsPendingChannel(portID, channelID) {
		return nil
	}

	seq, err := w.SubmitGenesisBridgeData(ctx, portID, channelID)
	if err != nil {
		return errorsmod.Wrap(err, "submit genesis bridge data")
	}

	// Mark this channel as having a pending genesis bridge submission
	w.k.AddPendingChannel(portID, channelID)

	w.logger(ctx).Info("genesis bridge data submitted", "sequence", seq, "port", portID, "channel", channelID)
	return nil
}

// SubmitGenesisBridgeData sends the genesis bridge data over the channel.
// The genesis bridge data includes the genesis info, the native denom metadata, and the genesis transfer packet.
// It uses the channel keeper to send the packet, instead of transfer keeper, as we are not sending fungible token directly.
func (w IBCModule) SubmitGenesisBridgeData(ctx sdk.Context, portID string, channelID string) (seq uint64, err error) {
	_, chanCap, err := w.channelKeeper.LookupModuleByChannel(ctx, portID, channelID)
	if err != nil {
		return
	}

	data, err := w.k.PrepareGenesisBridgeData(ctx)
	if err != nil {
		return 0, errorsmod.Wrap(err, "prepare genesis bridge data")
	}

	if data.GenesisTransfer != nil {
		// As we don't use the `ibc/transfer` module, we need to handle the funds escrow ourselves
		err = w.k.EscrowGenesisTransferFunds(ctx, portID, channelID, data.GenesisInfo.BaseCoinSupply())
		if err != nil {
			return 0, errorsmod.Wrap(err, "escrow genesis transfer funds")
		}
	}

	bz, err := json.Marshal(data)
	if err != nil {
		return 0, errorsmod.Wrap(err, "marshal genesis bridge data")
	}

	timeoutTimestamp := ctx.BlockTime().Add(transferTimeout).UnixNano()
	return w.channelKeeper.SendPacket(ctx, chanCap, portID, channelID, clienttypes.ZeroHeight(), uint64(timeoutTimestamp), bz)
}

func (w IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	state := w.k.GetState(ctx)

	// if outbound transfers are enabled, we past the genesis phase. nothing to do here.
	if state.OutboundTransfersEnabled {
		return w.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	var ack channeltypes.Acknowledgement
	err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack)
	if err != nil {
		return err
	}

	if !w.k.IsPendingChannel(packet.SourcePort, packet.SourceChannel) {
		// error?????
		return w.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	if !ack.Success() {
		// need to mark the packet as acked, so it will be retried
		w.k.onGoingChannels[packet.SourcePort+"/"+packet.SourceChannel] = true
		w.logger(ctx).Error("acknowledgement failed for genesis transfer", "packet", packet, "ack", ack)
		return errors.New("acknowledgement failed for genesis transfer")
	}

	// If ack success, enable the bridge for the successful channel
	w.k.enableBridge(ctx, state, packet.SourcePort, packet.SourceChannel)

	// clean retry queue
	w.k.ClearPendingChannels()

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeOutboundTransfersEnabled))
	w.logger(ctx).Info("genesis bridge phase completed successfully")

	return nil
}

// ontimeout packet
func (w IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {

	state := w.k.GetState(ctx)

	// if outbound transfers are enabled, no-op
	if state.OutboundTransfersEnabled {
		return w.IBCModule.OnTimeoutPacket(ctx, packet, relayer)
	}

	if !w.k.IsPendingChannel(packet.SourcePort, packet.SourceChannel) {
		w.logger(ctx).Error("timeout for unknown channel", "packet", packet)
		return errors.New("timeout for unknown channel")
	}

	// set the channel as retry required
	w.k.onGoingChannels[packet.SourcePort+"/"+packet.SourceChannel] = true

	return w.IBCModule.OnTimeoutPacket(ctx, packet, relayer)
}
