package keeper

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
	k         Keeper
}

func NewOnChanOpenConfirmInterceptor(
	k Keeper,
	transferK transferkeeper.Keeper,
	next porttypes.IBCModule,
) *OnChanOpenConfirmInterceptor {
	return &OnChanOpenConfirmInterceptor{next, transferK, k}
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

	state := i.k.GetState(ctx)

	firstCoin := state.GenesisTokens[0]

	srcAccount := i.k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	srcAddr := srcAccount.GetAddress()

	dstAddr := sdk.AccAddress("dym13d2qrv402klpu6t6qk0uvd8eqxmrw6srmsm4yu")

	m := types.MsgTransfer{
		SourcePort:       portID,
		SourceChannel:    channelID,
		Token:            firstCoin,
		Sender:           srcAddr.String(),
		Receiver:         dstAddr.String(),
		TimeoutHeight:    clienttypes.Height{},
		TimeoutTimestamp: 0,
		Memo:             "special",
	}
	res, err := i.transferK.Transfer(ctx.Context(), &m)
	if err != nil {
		ctx.Logger().Info("OnChanOpenConfirm transfer", "err", err)
	}
	_ = res

	return i.IBCModule.OnChanOpenConfirm(ctx, portID, channelID)
}
