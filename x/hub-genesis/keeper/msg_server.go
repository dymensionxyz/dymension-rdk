package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

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

	// Get the sender and validate they are in the Allowlist
	if !m.IsAddressInGenesisTriggererAllowList(ctx, msg.Address) {
		return nil, sdkerrors.ErrUnauthorized
	}

	_, clientState, err := m.channelKeeper.GetChannelClientState(ctx, "transfer", msg.ChannelId)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedGetClientState, "failed to get client state for channel %s: %v", msg.ChannelId, err)
	}

	tmClientState, ok := clientState.(*tenderminttypes.ClientState)
	if !ok {
		return nil, errorsmod.Wrapf(types.ErrFailedGetClientState, "expected tendermint client state, got %T", clientState)
	}

	if tmClientState.GetChainID() != msg.HubId {
		return nil, errorsmod.Wrapf(types.ErrChainIDMismatch, "channel %s is connected to chain ID %s",
			msg.ChannelId, tmClientState.GetChainID())
	}

	// check if genesis event was already triggered
	state := m.GetState(ctx)
	if state.IsLocked {
		return nil, types.ErrGenesisEventAlreadyTriggered
	}

	if err := m.lockRollappGenesisTokens(ctx, msg.ChannelId, state.GenesisTokens); err != nil {
		return nil, errorsmod.Wrapf(types.ErrLockingGenesisTokens, "failed to lock tokens: %v", err)
	}

	state.IsLocked = true
	m.SetState(ctx, state)

	return &types.MsgHubGenesisEventResponse{}, nil
}
