package keeper

import (
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// this line is used by starport scaffolding # proto/tx/rpc

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
