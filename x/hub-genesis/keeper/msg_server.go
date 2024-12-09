package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
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
	relayer := msg.GetSigners()[0].String() // guaranteed length 1
	err := m.SendGenesisTransfer(ctx, relayer, msg.ChannelId)
	if err != nil {
		return nil, err
	}
	return &types.MsgSendTransferResponse{}, nil
}

func (k Keeper) SendGenesisTransfer(ctx sdk.Context, relayer, channelID string) error {
	wl, err := whitelistedrelayer.GetList(ctx, k.dk, k.sk)
	if err != nil {
		return errorsmod.Wrap(err, "get whitelisted relayers")
	}
	if !wl.Has(relayer) {
		return gerrc.ErrPermissionDenied.Wrap("not whitelisted")
	}
}
