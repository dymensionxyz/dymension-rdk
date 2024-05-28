package keeper

import (
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

	// firstCoin := state.GenesisTokens[0] // TODO: send all transfers
	var firstCoin sdk.Coin
	// dstStr := "dym13d2qrv402klpu6t6qk0uvd8eqxmrw6srmsm4yu"
	var dstStr string
	for i, a := range state.GenesisAccounts {
		l.Info("got genesis account", "acc", a)
		if i == 0 {
			firstCoin = a.GetAmount()
			dstStr = a.GetAddress()
		}
	}

	srcAccount := i.k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	srcAddr := srcAccount.GetAddress()

	m := transfertypes.MsgTransfer{
		SourcePort:       portID,
		SourceChannel:    channelID,
		Token:            firstCoin,
		Sender:           srcAddr.String(),
		Receiver:         dstStr,
		TimeoutHeight:    clienttypes.Height{},
		TimeoutTimestamp: uint64(ctx.BlockTime().Add(time.Hour * 24).UnixNano()),
		Memo:             "special",
	}

	err = i.transfer(ctx, &m)
	if err != nil {
		ctx.Logger().Error("OnChanOpenConfirm transfer", "err", err)
	} else {
		ctx.Logger().Info("sent special transfer")
	}

	return nil
}
