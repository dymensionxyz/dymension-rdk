package keeper

import (
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

const (
	transferTimeout = time.Hour * 24 * 14
)

type IBCModule struct {
	porttypes.IBCModule
	k         Keeper
	transfer  Transfer
	getDenom  GetDenomMetaData
	mintCoins MintCoins
}

type (
	Transfer         func(ctx sdk.Context, transfer *transfertypes.MsgTransfer) error
	GetDenomMetaData func(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	MintCoins        func(ctx sdk.Context, moduleName string, amt sdk.Coins) error
)

func NewIBCModule(next porttypes.IBCModule, t Transfer, k Keeper, d GetDenomMetaData, m MintCoins) *IBCModule {
	return &IBCModule{next, k, t, d, m}
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
		l.Error("Next middleware.", "err", err)
		return err
	}

	state := w.k.GetState(ctx)

	srcAccount := w.k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	srcAddr := srcAccount.GetAddress().String()

	for i, a := range state.GetGenesisAccounts() {
		if err := w.mintAndTransfer(ctx, i, len(state.GetGenesisAccounts()), a, srcAddr, portID, channelID); err != nil {
			// there is no feasible way to recover
			panic(fmt.Errorf("mint and transfer: %w", err))
		}
		l.Info("Sent genesis transfer.", "index", i, "receiver", a.GetAddress(), "coin", a)
	}

	state.GenesisAccounts = nil

	w.k.SetState(ctx, state)

	l.Info("Sent all genesis transfers.")

	return nil
}

func (w IBCModule) mintAndTransfer(
	ctx sdk.Context,
	i, n int,
	a types.GenesisAccount,
	srcAddr string,
	portID string,
	channelID string,
) error {
	coin := a.GetAmount()
	err := w.mintCoins(ctx, types.ModuleName, sdk.Coins{coin})
	if err != nil {
		return errorsmod.Wrap(err, "mint coins")
	}

	// NOTE: for simplicity we don't optimize to avoid sending duplicate metadata
	// we assume the hub will deduplicate
	memo, err := w.createMemo(ctx, a.Amount.Denom, i, n)
	if err != nil {
		return errorsmod.Wrap(err, "create memo")
	}

	m := transfertypes.MsgTransfer{
		SourcePort:       portID,
		SourceChannel:    channelID,
		Token:            a.Amount,
		Sender:           srcAddr,
		Receiver:         a.GetAddress(),
		TimeoutHeight:    clienttypes.Height{},
		TimeoutTimestamp: uint64(ctx.BlockTime().Add(transferTimeout).UnixNano()),
		Memo:             memo,
	}

	err = w.transfer(skipContext(ctx), &m)
	if err != nil {
		return errorsmod.Wrap(err, "transfer")
	}

	return nil
}
