package keeper

import (
	"context"

	"github.com/dymensionxyz/rollapp/x/sequencers/types"
)

// this line is used by starport scaffolding # proto/tx/rpc

// CreateValidator defines a method for creating a new validator.
func (m msgServer) CreateSequencer(_ context.Context, _ *types.MsgCreateSequencer) (*types.MsgCreateSequencerResponse, error) {
	panic("not implemented") // TODO: Implement
}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
