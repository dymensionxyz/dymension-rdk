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
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

var transferTimeout = (time.Duration(24*365) * time.Hour) // very long timeout

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
		return nil
	}

	seq, err := w.SubmitGenesisBridgeData(ctx, portID, channelID)
	if err != nil {
		return errorsmod.Wrap(err, "submit genesis bridge data")
	}

	// Mark this channel as having a pending genesis bridge submission
	portChan := types.PortAndChannel{Port: portID, Channel: channelID}
	if err := w.k.SetPendingChannel(ctx, portChan, types.WaitingForAck); err != nil {
		return errorsmod.Wrap(err, "add pending channel")
	}

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

	// if canonical channel is set, we past the genesis phase. nothing to do here.
	if state.CanonicalHubTransferChannelHasBeenSet() {
		return w.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	var ack channeltypes.Acknowledgement
	err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack)
	if err != nil {
		return err
	}

	// validate it's a channel we are expecting
	portChannel := types.PortAndChannel{Port: packet.SourcePort, Channel: packet.SourceChannel}
	isPending, err := w.k.IsPendingChannel(ctx, portChannel)
	if err != nil {
		return errorsmod.Wrap(err, "check pending channel")
	}
	if !isPending {
		// Not supposed to happen
		w.logger(ctx).Error("genesis bridge acknowledgement for unknown channel", "packet", packet, "ack", ack)
		return errors.New("acknowledgement failed for unknown channel")
	}

	// Mark the channel as failed
	if !ack.Success() {
		w.logger(ctx).Error("acknowledgement failed for genesis transfer", "packet", packet, "ack", ack)
		portChan := types.PortAndChannel{Port: packet.SourcePort, Channel: packet.SourceChannel}
		if err := w.k.SetPendingChannel(ctx, portChan, types.Failed); err != nil {
			return errorsmod.Wrap(err, "set channel retry required")
		}

		return nil
	}

	var gbData types.GenesisBridgeData
	err = json.Unmarshal(packet.Data, &gbData)
	if err != nil {
		return errorsmod.Wrap(err, "unmarshal genesis bridge data")
	}

	if gbData.GenesisTransfer != nil {
		// As we don't use the `ibc/transfer` module, we need to handle the funds escrow ourselves
		err = w.k.EscrowGenesisTransferFunds(ctx, packet.SourcePort, packet.SourceChannel, gbData.GenesisInfo.BaseCoinSupply())
		if err != nil {
			return errorsmod.Wrap(err, "escrow genesis transfer funds")
		}
	}

	w.k.enableBridge(ctx, state, packet.SourcePort, packet.SourceChannel)

	if err := w.k.ClearPendingChannels(ctx); err != nil {
		return errorsmod.Wrap(err, "clear pending channels")
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeOutboundTransfersEnabled))
	w.logger(ctx).Info("genesis bridge phase completed successfully")

	return nil
}

// OnTimeoutPacket handles IBC packet timeouts
func (w IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	state := w.k.GetState(ctx)

	// if outbound transfers are enabled, no-op
	if state.CanonicalHubTransferChannelHasBeenSet() {
		return w.IBCModule.OnTimeoutPacket(ctx, packet, relayer)
	}

	return errorsmod.Wrapf(gerrc.ErrUnknown, "unexpected packet timeout: %s", packet.String())
}
