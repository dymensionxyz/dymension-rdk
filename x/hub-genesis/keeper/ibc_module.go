package keeper

import (
	"errors"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

var transferTimeout = time.Duration(24*365) * time.Hour // very long timeout

type IBCModule struct {
	porttypes.IBCModule
	k    Keeper
	bank types.BankKeeper
}

func NewIBCModule(next porttypes.IBCModule, k Keeper, bank types.BankKeeper) *IBCModule {
	return &IBCModule{next, k, bank}
}

func (w IBCModule) logger(ctx sdk.Context) log.Logger {
	return w.k.Logger(ctx)
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

	expect := state.GetHubPortAndChannel()
	got := types.PortAndChannel{
		Port:    packet.SourcePort,
		Channel: packet.SourceChannel,
	}
	if expect == nil || *expect != got {
		err := errorsmod.Wrap(gerrc.ErrInvalidArgument, "unexpected non genesis transfer packet before genesis bridge open")
		w.logger(ctx).Error("OnAcknowledgementPacket", "error", err)
		return err
	}

	state.InFlight = false
	w.k.SetState(ctx, state)

	var ack channeltypes.Acknowledgement
	err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack)
	if err != nil {
		return err
	}

	if !ack.Success() {
		w.logger(ctx).Error("acknowledgement failed for genesis transfer", "packet", packet, "ack", ack)
		return errors.New("acknowledgement failed for genesis transfer")
	}

	gfo := w.k.GetGenesisInfo(ctx)
	if !gfo.Amt().IsZero() {
		// As we don't use the `ibc/transfer` module, we need to handle the funds escrow ourselves
		err = w.k.EscrowGenesisTransferFunds(ctx, port, packet.SourceChannel, gfo.BaseCoinSupply())
		if err != nil {
			return errorsmod.Wrap(errors.Join(err, gerrc.ErrInternal), "escrow genesis transfer funds : rollapp is corrupted")
		}
	}

	// open the bridge
	state.OutboundTransfersEnabled = true
	w.k.SetState(ctx, state)

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeOutboundTransfersEnabled))
	w.logger(ctx).Info("genesis bridge phase completed successfully")
	return nil
}
