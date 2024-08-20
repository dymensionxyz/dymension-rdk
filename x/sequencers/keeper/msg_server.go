package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) CreateSequencer(goCtx context.Context, msg *types.MsgCreateSequencer) (*types.MsgCreateSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr := msg.MustSigner() // ensured in validate basic
	allow, err := m.IsSigned(ctx, accAddr, msg.GetKeyAndSig(), msg.GetPayload())
	if err != nil {
		return nil, errorsmod.Wrap(err, "check sig")
	}
	if !allow {
		return nil, gerrc.ErrUnauthenticated
	}
	operator := msg.MustOperator() // checked in validate basic
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
	accAddr := msg.MustSigner()
	allow, err := m.IsSigned(ctx, accAddr, msg.GetKeyAndSig(), msg.GetPayload())
	if err != nil {
		return nil, errorsmod.Wrap(err, "check sig")
	}
	if !allow {
		return nil, gerrc.ErrUnauthenticated
	}
	consAddr, err := msg.GetKeyAndSig().Validator().GetConsAddr()
	if err != nil {
		panic(err) // it must be ok because we used it to check sig
	}
	seq, ok := m.GetSequencerByConsAddr(ctx, consAddr)
	if !ok {
		return nil, errorsmod.Wrap(gerrc.ErrNotFound, "sequencer by cons addr")
	}
	m.SetRewardAddr(ctx, seq, msg.MustRewardAcc()) // We can Must because it's checked in validate basic

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventUpdateSequencer,
		sdk.NewAttribute(types.AttributeKeyConsAddr, consAddr.String()),
		sdk.NewAttribute(types.AttributeKeyOperatorAddr, seq.OperatorAddress),
		sdk.NewAttribute(types.AttributeKeyRewardAddr, msg.MustRewardAcc().String()),
	))
	return &types.MsgUpdateSequencerResponse{}, nil
}

var _ types.MsgServer = msgServer{}

// IsSigned return true iff the key and sig contains a key and signature where the signature was produced by the key, and the signature
// is over the account from the provided address, and the app payload data.
//
// The reasoning is as follows:
// We know that the TX containing the Msg was signed by addr, because it has passed the sdk signature verification ante.
// Therefore, if we require that the private key for the consensus address was used to sign off over this addr AND this chain ID then
// we know that the owner of the private key really intended this payload to be included in this transaction, and it is not man in the middle or replay.
func (k Keeper) IsSigned(ctx sdk.Context, addr sdk.AccAddress, keyAndSig *types.KeyAndSig, payloadApp codec.ProtoMarshaler) (bool, error) {
	acc := k.authAccountKeeper.GetAccount(ctx, addr)

	v := keyAndSig.Validator()

	payloadBz, err := types.CreateBytesToSign(
		ctx.ChainID(),
		acc.GetAccountNumber(),
		payloadApp,
	)
	if err != nil {
		return false, errorsmod.Wrap(err, "create bytes to sign")
	}

	pubKey, err := v.ConsPubKey()
	if err != nil {
		return false, errorsmod.Wrap(err, "get cons pubkey")
	}

	return pubKey.VerifySignature(payloadBz, keyAndSig.GetSignature()), nil
}
