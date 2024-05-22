package foo

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	transferkeeper "github.com/evmos/evmos/v12/x/ibc/transfer/keeper"
)

type OnChanOpenConfirmInterceptor struct {
	porttypes.IBCModule
	transferK transferkeeper.Keeper
}

func NewOnChanOpenConfirmInterceptor(
	keeper transferkeeper.Keeper,
	next porttypes.IBCModule,
) *OnChanOpenConfirmInterceptor {
	return &OnChanOpenConfirmInterceptor{next, keeper}
}

func (i OnChanOpenConfirmInterceptor) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	/*
		TODO: send a transfer on the channelID
	*/

	ctx.Logger().Info("OnChanOpenConfirm interceptor!", "port id", portID, "channelID", channelID)

	m := types.MsgTransfer{
		SourcePort:       portID,
		SourceChannel:    channelID,
		Token:            sdk.Coin{},
		Sender:           "",
		Receiver:         "",
		TimeoutHeight:    clienttypes.Height{},
		TimeoutTimestamp: 0,
		Memo:             "",
	}
	res, err := i.transferK.Transfer(ctx.Context(), &m)
	if err != nil {
		ctx.Logger().Info("OnChanOpenConfirm transfer", "err", err)
	}
	_ = res

	return i.IBCModule.OnChanOpenConfirm(ctx, portID, channelID)
}
