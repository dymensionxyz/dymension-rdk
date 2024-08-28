package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
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
	v := msg.Validator()
	cons, err := v.GetConsAddr()
	if err != nil {
		panic(err) // it must be ok because we used it to check sig
	}

	if _, ok := m.GetSequencer(ctx, operator); ok {
		return nil, gerrc.ErrAlreadyExists
	}
	if _, ok := m.GetSequencerByConsAddr(ctx, cons); ok {
		return nil, gerrc.ErrAlreadyExists
	}

	m.SetSequencer(ctx, v)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventCreateSequencer,
		sdk.NewAttribute(types.AttributeKeyConsAddr, cons.String()),
		sdk.NewAttribute(types.AttributeKeyOperatorAddr, v.OperatorAddress),
	))

	return &types.MsgCreateSequencerResponse{}, nil
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
