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

		m := transfertypes.MsgTransfer{
			SourcePort:       portID,
			SourceChannel:    channelID,
			Token:            a.Amount,
			Sender:           srcAddr,
			Receiver:         a.GetAddress(),
			TimeoutHeight:    clienttypes.Height{},
			TimeoutTimestamp: uint64(ctx.BlockTime().Add(time.Hour * 24).UnixNano()),
			Memo:             "special",
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
