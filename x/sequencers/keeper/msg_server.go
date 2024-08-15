package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

type msgServer struct {
	Keeper
}

func (m msgServer) CreateSequencer(goCtx context.Context, msg *types.MsgCreateSequencer) (*types.MsgCreateSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr := msg.MustGetSigner()

	// TODO implement me
	panic("implement me")
}

func (m msgServer) UpdateSequencer(goCtx context.Context, msg *types.MsgUpdateSequencer) (*types.MsgUpdateSequencerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	m.cdc.MustMarshal()

	// TODO implement me
	panic("implement me")
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k Keeper) CheckSig(ctx sdk.Context, addr sdk.AccAddress, keyAndSig types.KeyAndSig, payloadApp codec.ProtoMarshaler) (bool, error) {
	payloadAppBz, err := payloadApp.Marshal()
	if err != nil {
		return false, err
	}
	acc := k.authAccountKeeper.GetAccount(ctx, addr)

	payload := &types.PayloadToSign{
		PayloadApp:    payloadAppBz,
		ChainId:       ctx.ChainID(),
		AccountNumber: acc.GetAccountNumber(),
	}

	acc.GetSequence() // TODO: is sequence necessary, is addr necessary?
	acc.GetAccountNumber()
	acc.GetAddress()
	ctx.ChainID()
}

/*
Design
	MsgCreate
		creator (msg signer)
		SignedPart
		optional Update
		operator addr

	MsgUpdate
		creator (msg signer)
		SignedPart
		Update

	Update
		reward addr

	SignedPart
		pub key
		signature
			payload bz
			chain id
			account number
			account addr
			sequence

	On receipt
		You can do msg.GetSigners to get the sdk.AccAddress
		You can get this account from the SDK with accountKeeper.GetAccount
		This gives you
			accAddr
			accNum
			seqNum
			pubKey
		Then you compare against the signature in the SignedPart
			payload bz
			chain id
			accNum
		Then you know payload was signed by the pub key
*/
