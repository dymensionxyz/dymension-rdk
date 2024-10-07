package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) CreateSequencer(goCtx context.Context, msg *types.MsgCreateSequencer) (*types.MsgCreateSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	operator := msg.MustOperatorAddr() // checked in validate basic
	if _, ok := m.GetSequencer(ctx, operator); ok {
		return nil, gerrc.ErrAlreadyExists
	}

	v := msg.Validator()
	m.SetSequencer(ctx, v)

	consAddr, err := v.GetConsAddr()
	if err != nil {
		panic(err) // it must be ok because we used it to check sig
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventCreateSequencer,
		sdk.NewAttribute(types.AttributeKeyConsAddr, consAddr.String()),
		sdk.NewAttribute(types.AttributeKeyOperatorAddr, v.OperatorAddress),
	))

	return &types.MsgCreateSequencerResponse{}, nil
}

func (m msgServer) UpsertSequencer(goCtx context.Context, msg *types.MsgUpsertSequencer) (*types.MsgUpsertSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	v := msg.Validator()
	m.SetSequencer(ctx, v)

	consAddr, err := v.GetConsAddr()
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventUpsertSequencer,
		sdk.NewAttribute(types.AttributeKeyConsAddr, consAddr.String()),
		sdk.NewAttribute(types.AttributeKeyOperatorAddr, v.OperatorAddress),
		sdk.NewAttribute(types.AttributeKeyRewardAddr, msg.MustRewardAddr()),
	))

	return &types.MsgUpsertSequencerResponse{}, nil
}

func (m msgServer) UpdateSequencer(goCtx context.Context, msg *types.MsgUpdateSequencer) (*types.MsgUpdateSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	operator := msg.MustOperatorAddr() // checked in validate basic

	seq, ok := m.GetSequencer(ctx, operator)
	if !ok {
		return nil, errorsmod.Wrap(gerrc.ErrNotFound, "sequencer")
	}

	m.SetRewardAddr(ctx, seq, msg.MustRewardAcc()) // checked in validate basic

	consAddr, err := seq.GetConsAddr()
	if err != nil {
		return nil, errorsmod.Wrap(gerrc.ErrInternal, "expected to get valid cons addr")
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventUpdateSequencer,
		sdk.NewAttribute(types.AttributeKeyConsAddr, consAddr.String()),
		sdk.NewAttribute(types.AttributeKeyOperatorAddr, seq.OperatorAddress),
		sdk.NewAttribute(types.AttributeKeyRewardAddr, msg.MustRewardAcc().String()),
	))
	return &types.MsgUpdateSequencerResponse{}, nil
}
