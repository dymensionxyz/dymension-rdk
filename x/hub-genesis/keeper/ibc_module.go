package keeper

import (
	"context"
	"errors"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

const (
	transferTimeout = time.Hour * 24 * 365
)

type IBCModule struct {
	porttypes.IBCModule
	k        Keeper
	transfer TransferKeeper
	bank     BankKeeper
}

type TransferKeeper interface {
	Transfer(goCtx context.Context, msg *transfertypes.MsgTransfer) (*transfertypes.MsgTransferResponse, error)
}

type BankKeeper interface {
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

func NewIBCModule(next porttypes.IBCModule, t TransferKeeper, k Keeper, bank BankKeeper) *IBCModule {
	return &IBCModule{next, k, t, bank}
}

func (w IBCModule) logger(ctx sdk.Context) log.Logger {
	return w.k.Logger(ctx).With("module", types.ModuleName, "component", "ibc module")
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
	l := ctx.Logger().With("module", "hubgenesis OnChanOpenConfirm middleware", "port id", portID, "channelID", channelID)

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

	state.SetCanonicalTransferChannel(portID, channelID)

	state.NumUnackedTransfers = uint64(len(state.GetGenesisAccounts()))

	w.k.SetState(ctx, state)

	if len(state.GetGenesisAccounts()) == 0 {
		// we want to handle the case where the rollapp doesn't have genesis transfers
		// normally we would enable outbound transfers on an ack, but in this case we won't have an ack
		w.k.enableOutboundTransfers(ctx)
		return nil
	}

	srcAccount := w.k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	srcAddr := srcAccount.GetAddress().String()

	for i, a := range state.GetGenesisAccounts() {
		if err := w.mintAndTransfer(ctx, a, srcAddr, portID, channelID); err != nil {
			// there is no feasible way to recover
			panic(fmt.Errorf("mint and transfer: %w", err))
		}
		l.Info("Sent genesis transfer.", "index", i, "receiver", a.GetAddress(), "coin", a)
	}

	l.Info("Sent all genesis transfers.", "n", len(state.GetGenesisAccounts()))

	return nil
}

func (w IBCModule) mintAndTransfer(ctx sdk.Context, account types.GenesisAccount, srcAddr string, portID string, channelID string) error {
	coin := account.GetAmount()
	err := w.bank.MintCoins(ctx, types.ModuleName, sdk.Coins{coin})
	if err != nil {
		return errorsmod.Wrap(err, "mint coins")
	}

	// NOTE: for simplicity we don't optimize to avoid sending duplicate metadata
	// we assume the hub will deduplicate. We expect to eventually get a timeout
	// or commit anyway, so the packet will be cleared up.
	// (Actually, since transfers may arrive out of order, we must include the
	// denom metadata anyway).
	memo, err := w.createMemo(ctx, account.Amount.Denom)
	if err != nil {
		return errorsmod.Wrap(err, "create memo")
	}

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
		return errorsmod.Wrap(err, "transfer")
	}
	w.k.saveUnackedTransferSeqNum(ctx, res.Sequence)
	return nil
}

func (w IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	l := w.logger(ctx)
	state := w.k.GetState(ctx)

	if !state.IsCanonicalHubTransferChannel(packet.SourcePort, packet.SourceChannel) {
		return w.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	if state.OutboundTransfersEnabled {
		return w.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	var data transfertypes.FungibleTokenPacketData
	err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data)
	if err != nil {
		return err
	}

	var ack channeltypes.Acknowledgement
	err = types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack)
	if err != nil {
		return err
	}

	err = w.k.ackTransferSeqNum(ctx, packet.Sequence, ack)
	if err != nil {
		l.Error("Processing ack from transfer.", "err", err)
		return errorsmod.Wrap(errors.Join(err, gerrc.ErrUnknown), "ack transfer seq num")
	}
	return w.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}
