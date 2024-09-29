package keeper

import (
	"errors"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

const (
	transferTimeout = time.Hour * 24 * 365
)

type IBCModule struct {
	porttypes.IBCModule
	k             Keeper
	transfer      types.TransferKeeper
	bank          types.BankKeeper
	channelKeeper types.ChannelKeeper
}

func NewIBCModule(next porttypes.IBCModule, t types.TransferKeeper, k Keeper, bank types.BankKeeper, chanKeeper types.ChannelKeeper) *IBCModule {
	return &IBCModule{next, k, t, bank, chanKeeper}
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

		// FIXME: return error if it's already set? why we need to support additional channels after we have the canonical one?
		return nil
	}
	state.SetCanonicalTransferChannel(portID, channelID)

	// create the memo with the genesis data
	memo, err := w.CreateGenesisMemo(ctx)
	if err != nil {
		return errorsmod.Wrap(err, "create genesis memo")
	}

	// FIXME: refactor
	// we always wait for genesis ack
	state.NumUnackedTransfers = 1
	w.k.SetState(ctx, state)

	var sequence uint64
	fundIRO := len(state.GetGenesisAccounts()) > 0
	if !fundIRO {
		sequence, err = w.submitMemoOnly(ctx, portID, channelID, memo)
		if err != nil {
			return errorsmod.Wrap(err, "submit memo only")
		}
		l.Info("Sent genesis memo only.")
	} else {
		srcAccount := w.k.ak.GetModuleAccount(ctx, types.ModuleName)
		srcAddr := srcAccount.GetAddress().String()

		genAcc := state.GetGenesisAccounts()[0]
		sequence, err = w.SubmitGenesisFunds(ctx, genAcc, srcAddr, portID, channelID, memo)
		if err != nil {
			return errorsmod.Wrap(err, "mint and transfer")
		}
		l.Info("Sent genesis transfer.", "receiver", genAcc.GetAddress(), "coin", genAcc.Amount)
	}

	w.k.saveUnackedTransferSeqNum(ctx, sequence)

	return nil
}

// sendMemoOnly
func (w IBCModule) submitMemoOnly(ctx sdk.Context, portID string, channelID string, memo string) (seq uint64, err error) {
	_, chanCap, err := w.channelKeeper.LookupModuleByChannel(ctx, portID, channelID)
	if err != nil {
		return
	}

	// FIXME: wrap memo encoding into cutom strucr
	return w.channelKeeper.SendPacket(ctx, chanCap, portID, channelID, clienttypes.Height{0, 0}, uint64(transferTimeout.Nanoseconds()), []byte(memo))
}

func (w IBCModule) SubmitGenesisFunds(ctx sdk.Context, account types.GenesisAccount, srcAddr, portID, channelID, memo string) (seq uint64, err error) {
	m := transfertypes.MsgTransfer{
		SourcePort:       portID,
		SourceChannel:    channelID,
		Token:            account.Amount,
		Sender:           srcAddr,
		Receiver:         account.GetAddress(),
		TimeoutHeight:    clienttypes.Height{},
		TimeoutTimestamp: uint64(ctx.BlockTime().Add(transferTimeout).UnixNano()),
		Memo:             memo,
	}

	res, err := w.transfer.Transfer(sdk.WrapSDKContext(allowSpecialMemoCtx(ctx)), &m)
	if err != nil {
		return 0, errorsmod.Wrap(err, "transfer")
	}
	return res.Sequence, nil
}

func (w IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	state := w.k.GetState(ctx)
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

	w.k.ackTransferSeqNum(ctx, packet.Sequence)
	w.logger(ctx).Info("genesis transfer phase acked successfully")

	return w.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}
