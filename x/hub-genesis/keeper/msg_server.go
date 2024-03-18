package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	tenderminttypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"
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

func (m msgServer) TriggerGenesisEvent(goCtx context.Context, msg *types.MsgHubGenesisEvent) (*types.MsgHubGenesisEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the sender and validate they are in the whitelist
	if !m.IsAddressInGenesisTriggererWhiteList(ctx, msg.Address) {
		return nil, sdkerrors.ErrUnauthorized
	}

	_, clientState, err := m.channelKeeper.GetChannelClientState(ctx, "transfer", msg.ChannelId)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidGenesisChannelId, "failed to get client state for channel %s", msg.ChannelId)
	}

	tmClientState, ok := clientState.(*tenderminttypes.ClientState)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrInvalidGenesisChannelId, "expected tendermint client state, got %T", clientState)
	}

	if tmClientState.GetChainID() != msg.HubId {
		return nil, sdkerrors.Wrapf(types.ErrInvalidGenesisChainId, "channel %s is connected to chain ID %s, expected %s",
			msg.ChannelId, tmClientState.GetChainID(), msg.HubId)
	}

	// if the hub is found, the genesis event was already triggered
	_, found := m.GetHub(ctx, msg.HubId)
	if found {
		return nil, types.ErrGenesisEventAlreadyTriggered
	}

	hub := types.NewHub(msg.HubId, msg.ChannelId)

	if err := m.lockRollappGenesisTokens(ctx, hub.ChannelId); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrLockingGenesisTokens, "failed to lock tokens: %v", err)
	}

	// we save the hub in order to prevent the genesis event from being triggered again
	m.SetHub(ctx, hub)

	return &types.MsgHubGenesisEventResponse{}, nil
}
