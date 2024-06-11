package keeper

import (
	"context"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) TriggerGenesisEvent(context.Context, *types.MsgHubGenesisEvent) (*types.MsgHubGenesisEventResponse, error) {
	return &types.MsgHubGenesisEventResponse{}, nil
}
