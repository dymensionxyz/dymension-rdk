package keeper

import (
	"context"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) SoftwareUpgrade(ctx context.Context, upgrade *types.MsgSoftwareUpgrade) (*types.MsgSoftwareUpgradeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m msgServer) CancelUpgrade(ctx context.Context, upgrade *types.MsgCancelUpgrade) (*types.MsgCancelUpgradeResponse, error) {
	//TODO implement me
	panic("implement me")
}
