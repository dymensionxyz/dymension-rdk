package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/utils/whitelistedrelayer"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

var _ types.MsgServer = msgServer{}

type msgServer struct{ Keeper }

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) SendTransfer(goCtx context.Context, msg *types.MsgSendTransfer) (*types.MsgSendTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	wl, err := whitelistedrelayer.GetList(ctx, m.dk, m.sk)
	if err != nil {
		return nil, err
	}
	if !wl.Has(msg.Relayer) {
		return nil, gerrc.ErrPermissionDenied.Wrap("not whitelisted")
	}
	return &types.MsgSendTransferResponse{}, nil
}
