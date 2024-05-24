package keeper

import (
	"context"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

func (m msgServer) TriggerGenesisEvent(context.Context, *types.MsgHubGenesisEvent) (*types.MsgHubGenesisEventResponse, error) {
	return &types.MsgHubGenesisEventResponse{}, nil
}
