package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"

	"github.com/dymensionxyz/dymension-rdk/utils/uevent"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct{ Keeper }

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) UpdateRewardAddress(goCtx context.Context, msg *types.MsgUpdateRewardAddress) (*types.MsgUpdateRewardAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// all must-methods are safe to use since they're validated in ValidateBasic

	operator := msg.MustOperatorAddr()
	seq, ok := m.GetSequencer(ctx, operator)
	if !ok {
		return nil, errorsmod.Wrap(gerrc.ErrNotFound, "sequencer")
	}

	rewardAddr := msg.MustRewardAcc()
	m.SetRewardAddr(ctx, seq, rewardAddr)

	err := uevent.EmitTypedEvent(ctx, &types.EventUpdateRewardAddress{
		Operator:   operator.String(),
		RewardAddr: rewardAddr.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.MsgUpdateRewardAddressResponse{}, nil
}

func (m msgServer) UpdateWhitelistedRelayers(goCtx context.Context, msg *types.MsgUpdateWhitelistedRelayers) (*types.MsgUpdateWhitelistedRelayersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// all must-methods are safe to use since they're validated in ValidateBasic

	operator := msg.MustOperatorAddr()
	seq, ok := m.GetSequencer(ctx, operator)
	if !ok {
		return nil, errorsmod.Wrap(gerrc.ErrNotFound, "sequencer")
	}

	relayers := types.MustNewWhitelistedRelayers(msg.Relayers)
	err := m.SetWhitelistedRelayers(ctx, seq, relayers)
	if err != nil {
		return nil, fmt.Errorf("set whitelisted relayers: %w", err)
	}

	err = uevent.EmitTypedEvent(ctx, &types.EventUpdateWhitelistedRelayers{
		Operator: seq.OperatorAddress,
		Relayers: relayers.Relayers,
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.MsgUpdateWhitelistedRelayersResponse{}, nil
}

func (m msgServer) UpsertSequencer(goCtx context.Context, msg *types.ConsensusMsgUpsertSequencer) (*types.ConsensusMsgUpsertSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// all must-methods are safe to use since they're validated in ValidateBasic

	// save the validator
	v := msg.MustValidator()
	m.SetSequencer(ctx, v)

	// save the reward address
	rewardAddr := msg.MustRewardAddr()
	m.SetRewardAddr(ctx, v, rewardAddr)

	// save the whitelisted relayer list
	err := m.SetWhitelistedRelayers(ctx, v, types.MustNewWhitelistedRelayers(msg.Relayers))
	if err != nil {
		return nil, fmt.Errorf("set whitelisted relayers: %w", err)
	}

	consAddr, err := v.GetConsAddr()
	if err != nil {
		return nil, fmt.Errorf("get validator consensus addr: %w", err)
	}

	err = uevent.EmitTypedEvent(ctx, &types.EventUpsertSequencer{
		Operator:   msg.MustOperatorAddr().String(),
		ConsAddr:   consAddr.String(),
		RewardAddr: rewardAddr.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("emit event: %w", err)
	}

	return &types.ConsensusMsgUpsertSequencerResponse{}, nil
}
