package keeper

import (
	"errors"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

type DenomMetadataKeeper interface {
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
}

type OnChanOpenConfirmInterceptor struct {
	porttypes.IBCModule
	transfer Transfer
	k        Keeper
	denomK   DenomMetadataKeeper
}

type Transfer func(ctx sdk.Context, transfer *transfertypes.MsgTransfer) error

func NewOnChanOpenConfirmInterceptor(next porttypes.IBCModule, t Transfer, k Keeper, denomK DenomMetadataKeeper) *OnChanOpenConfirmInterceptor {
	return &OnChanOpenConfirmInterceptor{next, t, k, denomK}
}

func (i OnChanOpenConfirmInterceptor) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	l := ctx.Logger().With("name", "OnChanOpenConfirm interceptor!", "port id", portID, "channelID", channelID)

	err := i.IBCModule.OnChanOpenConfirm(ctx, portID, channelID)
	if err != nil {
		l.Error("Next middleware: on OnChanOpenConfirm", "err", err)
		return err
	}

	state := i.k.GetState(ctx)

	srcAccount := i.k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	srcAddr := srcAccount.GetAddress().String()

	var errs []error

	for _, a := range state.GetGenesisAccounts() {

		// NOTE: for simplicity we don't optimize to avoid sending duplicate metadata
		// we assume the hub will deduplicate
		memo := i.createMemo(ctx, a.Amount.Denom)

		m := transfertypes.MsgTransfer{
			SourcePort:       portID,
			SourceChannel:    channelID,
			Token:            a.Amount,
			Sender:           srcAddr,
			Receiver:         a.GetAddress(),
			TimeoutHeight:    clienttypes.Height{},
			TimeoutTimestamp: uint64(ctx.BlockTime().Add(time.Hour * 24).UnixNano()),
			Memo:             memo,
		}

		err = i.transfer(ctx, &m)

		if err == nil {
			ctx.Logger().Info("sent special transfer")
			continue
		}

		err = fmt.Errorf("transfer: receiver: %s: amt: %s", a.GetAddress(), a.Amount.String())
		errs = append(errs, err)

		ctx.Logger().Error("OnChanOpenConfirm transfer", "err", err) // TODO: don't log(?)
	}

	return errors.Join(errs...)
}

// createMemo creates a memo to go with the transfer. It's used by the hub to confirm
// that the transfer originated from the chain itself, rather than a user of the chain.
// It may also contain token metadata.
func (i OnChanOpenConfirmInterceptor) createMemo(ctx sdk.Context, denom string) string {
	i.denomK.GetDenomMetaData(ctx, denom)
}
