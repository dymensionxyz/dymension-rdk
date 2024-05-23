package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	transferkeeper "github.com/evmos/evmos/v12/x/ibc/transfer/keeper"
)

type OnChanOpenConfirmInterceptor struct {
	porttypes.IBCModule
	transferK transferkeeper.Keeper
	k         Keeper
}

func NewOnChanOpenConfirmInterceptor(next porttypes.IBCModule, transferK transferkeeper.Keeper, k Keeper) *OnChanOpenConfirmInterceptor {
	return &OnChanOpenConfirmInterceptor{next, transferK, k}
}

func (i OnChanOpenConfirmInterceptor) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	l := ctx.Logger().With("name", "OnChanOpenConfirm interceptor!", "port id", portID, "channelID", channelID)

	err := i.IBCModule.OnChanOpenConfirm(ctx, portID, channelID)
	if err != nil {
		l.Error("Passed on OnChanOpenConfirm", "err", err)
		return err
	}

	state := i.k.GetState(ctx)

	firstCoin := state.GenesisTokens[0] // TODO: send all transfers

	srcAccount := i.k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	srcAddr := srcAccount.GetAddress()

	dstAddr := sdk.AccAddress("dym13d2qrv402klpu6t6qk0uvd8eqxmrw6srmsm4yu")

	m := transfertypes.MsgTransfer{
		SourcePort:       portID,
		SourceChannel:    channelID,
		Token:            firstCoin,
		Sender:           srcAddr.String(),
		Receiver:         dstAddr.String(),
		TimeoutHeight:    clienttypes.Height{},
		TimeoutTimestamp: uint64(ctx.BlockTime().Add(time.Hour * 24).UnixNano()),
		Memo:             "special",
	}

	_, err = i.transferK.Transfer(sdk.WrapSDKContext(ctx), &m)
	if err != nil {
		ctx.Logger().Error("OnChanOpenConfirm transfer", "err", err)
	} else {
		ctx.Logger().Info("sent special transfer")
	}

	return nil
}
