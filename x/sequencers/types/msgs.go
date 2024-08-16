package types

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

var (
	_ sdk.Msg = (*MsgCreateSequencer)(nil)
	_ sdk.Msg = (*MsgUpdateSequencer)(nil)
)

func (m *MsgCreateSequencer) ValidateBasic() error {
	if _, err := m.GetSigner(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get signer")
	}
	// TODO implement me
	panic("implement me")
}

func (m *MsgCreateSequencer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

func (m *MsgCreateSequencer) MustGetSigner() sdk.AccAddress {
	addr, err := m.GetSigner()
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *MsgCreateSequencer) GetSigner() (sdk.AccAddress, error) {
	addr, err := sdk.AccAddressFromBech32(m.Creator)
	return addr, errorsmod.Wrap(err, "acc addr from bech32")
}

func (m *MsgUpdateSequencer) ValidateBasic() error {
	if _, err := m.GetSigner(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get signer")
	}
	// TODO implement me
	panic("implement me")
}

func (m *MsgUpdateSequencer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

func (m *MsgUpdateSequencer) MustGetSigner() sdk.AccAddress {
	addr, err := m.GetSigner()
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *MsgUpdateSequencer) GetSigner() (sdk.AccAddress, error) {
	addr, err := sdk.AccAddressFromBech32(m.Creator)
	return addr, errorsmod.Wrap(err, "acc addr from bech32")
}

func (m *MsgUpdateSequencer) MustRewardAccAddr() sdk.AccAddress {
	s := m.GetPayload().GetRewardAddr()
	return sdk.MustAccAddressFromBech32(s)
}

// Validator is a convenience method - it returns a validator object which already
// has implementations of various useful methods like obtaining various type conversions
// for the public key.
func (m *KeyAndSig) Validator() stakingtypes.Validator {
	return stakingtypes.Validator{ConsensusPubkey: m.PubKey}
}

type CreatorAccount interface {
	GetAddress() sdk.AccAddress
	GetAccountNumber() uint64
}

type SigningData struct {
	Account CreatorAccount
	ChainID string
	PubKey  cryptotypes.PubKey
}

func BuildMsgCreateSequencer(
	signingData SigningData,
	payload *CreateSequencerPayload,
) (*MsgCreateSequencer, error) {
	keyAndSig, creator, err := createKeyAndSigAndCreator(signingData, payload)
	if err != nil {
		return nil, fmt.Errorf("create key and sig: %w", err)
	}
	return &MsgCreateSequencer{
		Creator:   creator.String(),
		KeyAndSig: keyAndSig,
		Payload:   payload,
	}, nil
}

func BuildMsgUpdateSequencer(
	signingData SigningData,
	payload *UpdateSequencerPayload,
) (*MsgUpdateSequencer, error) {
	keyAndSig, creator, err := createKeyAndSigAndCreator(signingData, payload)
	if err != nil {
		return nil, fmt.Errorf("create key and sig: %w", err)
	}
	return &MsgUpdateSequencer{
		Creator:   creator.String(),
		KeyAndSig: keyAndSig,
		Payload:   payload,
	}, nil
}

func createKeyAndSigAndCreator(
	signingData SigningData,
	payload codec.ProtoMarshaler,
) (*KeyAndSig, sdk.AccAddress, error) {
	toSign, err := CreateBytesToSign(signingData.ChainID, signingData.Account.GetAccountNumber(), payload)
	if err != nil {
		return nil, sdk.AccAddress{}, fmt.Errorf("create payload to sign: %w", err)
	}

	var sig []byte
	// TODO: sign

	pubKeyAny, err := codectypes.NewAnyWithValue(signingData.PubKey)
	if err != nil {
		return nil, sdk.AccAddress{}, errorsmod.Wrap(err, "pubkey to any")
	}

	return &KeyAndSig{
		PubKey:    pubKeyAny,
		Signature: sig,
	}, signingData.Account.GetAddress(), nil
}

// CreateBytesToSign creates the bytes which must be signed
// Used to do the initial signing, and then also to verify signature of original data
func CreateBytesToSign(
	chainID string,
	accountNumber uint64,
	payload codec.ProtoMarshaler,
) ([]byte, error) {
	payloadBz, err := payload.Marshal()
	if err != nil {
		return nil, err
	}
	toSign := &PayloadToSign{
		PayloadApp:    payloadBz,
		ChainId:       chainID,
		AccountNumber: accountNumber,
	}
	return toSign.Marshal()
}
