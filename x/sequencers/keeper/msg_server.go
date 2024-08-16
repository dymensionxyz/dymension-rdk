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

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) CreateSequencer(goCtx context.Context, msg *types.MsgCreateSequencer) (*types.MsgCreateSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr := msg.MustGetSigner()
	allow, err := m.CheckSig(ctx, accAddr, msg.GetKeyAndSig(), msg.GetPayload())
	if err != nil {
		return nil, errorsmod.Wrap(err, "check sig")
	}
	if !allow {
		return nil, gerrc.ErrUnauthenticated
	}

	v := msg.GetKeyAndSig().Validator()
	v.OperatorAddress = msg.GetPayload().GetOperatorAddr()
	m.SetSequencer(ctx, v)
	return &types.MsgCreateSequencerResponse{}, nil
}

func (m msgServer) UpdateSequencer(goCtx context.Context, msg *types.MsgUpdateSequencer) (*types.MsgUpdateSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr := msg.MustGetSigner()
	allow, err := m.CheckSig(ctx, accAddr, msg.GetKeyAndSig(), msg.GetPayload())
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
	m.SetRewardAddr(ctx, seq, msg.MustRewardAccAddr())
	return &types.MsgUpdateSequencerResponse{}, nil
}

var _ types.MsgServer = msgServer{}

// CheckSig return true iff the key and sig contains a key and signature where the signature was produced by the key, and the signature
// is over the account from the provided address, and the app payload data.
//
// The reasoning is as follows:
// We know that the TX containing the Msg was signed by addr, because it has passed the sdk signature verification ante.
// Therefore, if we require that the private key for the consensus address was used to sign off over this addr AND this chain ID then
// we know that the owner of the private key really intended this payload to be included in this transaction, and it is not man in the middle or replay.
func (k Keeper) CheckSig(ctx sdk.Context, addr sdk.AccAddress, keyAndSig *types.KeyAndSig, payloadApp codec.ProtoMarshaler) (bool, error) {
	acc := k.authAccountKeeper.GetAccount(ctx, addr)

	v := keyAndSig.Validator()

	pubKey, err := v.ConsPubKey()
	if err != nil {
		return false, err
	}

	payloadAppBz, err := payloadApp.Marshal()
	if err != nil {
		return false, err
	}

	payload := &types.PayloadToSign{
		PayloadApp:    payloadAppBz,
		ChainId:       ctx.ChainID(),
		AccountNumber: acc.GetAccountNumber(),
	}

	payloadBz, err := payload.Marshal()
	if err != nil {
		return false, err
	}

	ok := pubKey.VerifySignature(payloadBz, keyAndSig.GetSignature())
	return ok, nil
}
