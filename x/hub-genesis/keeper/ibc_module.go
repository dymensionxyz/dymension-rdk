package keeper

import (
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

var (
	transferTimeout = uint64((time.Duration(24*365) * time.Hour).Nanoseconds()) // very long timeout
)

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

// OnChanOpenConfirm will send any unsent genesis account transfers over the channel.
// It is ASSUMED that the channel is for the Hub. This can be ensured by not exposing
// the sequencer API until after genesis is complete.
// Since transfers are only sent once, it does not matter if someone else tries to open
// a channel in future (it will no-op).
func (w IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	l := w.logger(ctx).With("method", "OnChanOpenConfirm", "port id", portID, "channelID", channelID)

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

	seq, err := w.SubmitGenesisBridgeData(ctx, portID, channelID)
	if err != nil {
		return errorsmod.Wrap(err, "submit genesis bridge data")
	}

	state.SetCanonicalTransferChannel(portID, channelID)

	l.Info("genesis bridge data submitted", "sequence", seq, "port", portID, "channel", channelID)
	return nil
}

func (w IBCModule) SubmitGenesisBridgeData(ctx sdk.Context, portID string, channelID string) (seq uint64, err error) {
	_, chanCap, err := w.channelKeeper.LookupModuleByChannel(ctx, portID, channelID)
	if err != nil {
		return
	}

	// prepare genesis info
	gInfo := w.k.GetGenesisInfo(ctx)

	// prepare the denom metadata
	d, ok := w.bank.GetDenomMetaData(ctx, gInfo.BaseDenom())
	if !ok {
		return 0, errorsmod.Wrap(gerrc.ErrInternal, "denom metadata not found")
	}

	// prepare the genesis transfer
	genesisTransferPacket, err := w.k.PrepareGenesisTransfer(ctx, portID, channelID)
	if err != nil {
		return 0, errorsmod.Wrap(err, "genesis transfer")
	}

	data := &types.GenesisBridgeData{
		GenesisInfo:     gInfo,
		NativeDenom:     d,
		GenesisTransfer: genesisTransferPacket,
	}

	bz, err := types.ModuleCdc.MarshalJSON(data)
	if err != nil {
		return 0, errorsmod.Wrap(err, "marshal genesis bridge data")
	}

	return w.channelKeeper.SendPacket(ctx, chanCap, portID, channelID, clienttypes.ZeroHeight(), transferTimeout, bz)
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

	if !ack.Success() {
		w.logger(ctx).Error("acknowledgement failed for genesis transfer", "packet", packet, "ack", ack)
		return errors.New("acknowledgement failed for genesis transfer")
	}

	err = w.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	if err != nil {
		return err
	}

	w.k.enableBridge(ctx, state)
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeOutboundTransfersEnabled))
	w.logger(ctx).Info("genesis bridge phase completed successfully")

	return nil
}

// enableBridge enables the bridge after successful genesis bridge phase.
func (k Keeper) enableBridge(ctx sdk.Context, state types.State) {
	state.OutboundTransfersEnabled = true
	k.SetState(ctx, state)
}
