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
	acc := m.authAccountKeeper.GetAccount(ctx, msg.MustAccAddr()) // ensured in validate basic
	if err := msg.GetKeyAndSig().Ok(ctx, acc, msg.GetPayload()); err != nil {
		return nil, errorsmod.Wrap(err, "check sig ok")
	}
	operator := msg.MustOperatorAddr() // checked in validate basic
	if _, ok := m.GetSequencer(ctx, operator); ok {
		return nil, gerrc.ErrAlreadyExists
	}

	v := msg.GetKeyAndSig().Validator()
	v.OperatorAddress = operator.String()
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

func (m msgServer) UpdateSequencer(goCtx context.Context, msg *types.MsgUpdateSequencer) (*types.MsgUpdateSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	acc := m.authAccountKeeper.GetAccount(ctx, msg.MustAccAddr()) // ensured in validate basic
	if err := msg.GetKeyAndSig().Ok(ctx, acc, msg.GetPayload()); err != nil {
		return nil, errorsmod.Wrap(err, "check sig ok")
	}
	consAddr, err := msg.GetKeyAndSig().Validator().GetConsAddr()
	if err != nil {
		panic(err) // it must be ok because we used it to check sig
	}
	seq, ok := m.GetSequencerByConsAddr(ctx, consAddr)
	if !ok {
		return nil, errorsmod.Wrap(gerrc.ErrNotFound, "sequencer by cons addr")
	}
	m.SetRewardAddr(ctx, seq, msg.MustRewardAcc()) // checked in validate basic

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventUpdateSequencer,
		sdk.NewAttribute(types.AttributeKeyConsAddr, consAddr.String()),
		sdk.NewAttribute(types.AttributeKeyOperatorAddr, seq.OperatorAddress),
		sdk.NewAttribute(types.AttributeKeyRewardAddr, msg.MustRewardAcc().String()),
	))
	return &types.MsgUpdateSequencerResponse{}, nil
}
