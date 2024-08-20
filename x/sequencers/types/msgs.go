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
	_ sdk.Msg                            = (*MsgCreateSequencer)(nil)
	_ sdk.Msg                            = (*MsgUpdateSequencer)(nil)
	_ codectypes.UnpackInterfacesMessage = (*MsgCreateSequencer)(nil)
	_ codectypes.UnpackInterfacesMessage = (*MsgUpdateSequencer)(nil)
)

func (m *MsgCreateSequencer) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return m.GetKeyAndSig().UnpackInterfaces(unpacker)
}

func (m *MsgUpdateSequencer) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return m.GetKeyAndSig().UnpackInterfaces(unpacker)
}

func (m *KeyAndSig) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return unpacker.UnpackAny(m.PubKey, new(cryptotypes.PubKey))
}

func (m *KeyAndSig) Valid() error {
	if m.GetPubKey() == nil {
		return errors.New("pub key is nil")
	}
	if m.GetPubKey().GetCachedValue() == nil {
		return errors.New("pub key cached value is nil")
	}
	v := stakingtypes.Validator{
		ConsensusPubkey: m.GetPubKey(),
	}
	tm, err := v.TmConsPublicKey()
	if err != nil {
		return errorsmod.Wrap(err, "tm cons pub key")
	}
	if tm.GetEd25519() == nil {
		return errors.New("not ed5519")
	}
	return nil
}

func (m *MsgCreateSequencer) ValidateBasic() error {
	if _, err := m.Signer(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get signer")
	}
	if err := m.KeyAndSig.Valid(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "key and sig")
	}
	if _, err := m.Operator(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "operator")
	}
	return nil
}

func (m *MsgCreateSequencer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustSigner()}
}

func (m *MsgCreateSequencer) MustSigner() sdk.AccAddress {
	addr, err := m.Signer()
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *MsgCreateSequencer) Signer() (sdk.AccAddress, error) {
	addr, err := sdk.AccAddressFromBech32(m.Creator)
	return addr, errorsmod.Wrap(err, "acc addr from bech32")
}

func (m *MsgCreateSequencer) Operator() (sdk.ValAddress, error) {
	return sdk.ValAddressFromBech32(m.GetPayload().GetOperatorAddr())
}

func (m *MsgCreateSequencer) MustOperator() sdk.ValAddress {
	addr, err := m.Operator()
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *MsgUpdateSequencer) ValidateBasic() error {
	if _, err := m.Signer(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get signer")
	}
	if err := m.KeyAndSig.Valid(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "key and sig")
	}
	if _, err := m.RewardAcc(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "reward addr")
	}
	return nil
}

func (m *MsgUpdateSequencer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustSigner()}
}

func (m *MsgUpdateSequencer) MustSigner() sdk.AccAddress {
	addr, err := m.Signer()
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *MsgUpdateSequencer) Signer() (sdk.AccAddress, error) {
	addr, err := sdk.AccAddressFromBech32(m.Creator)
	return addr, errorsmod.Wrap(err, "acc addr from bech32")
}

func (m *MsgUpdateSequencer) RewardAcc() (sdk.AccAddress, error) {
	return sdk.AccAddressFromBech32(m.GetPayload().GetRewardAddr())
}

func (m *MsgUpdateSequencer) MustRewardAcc() sdk.AccAddress {
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
	Signer  func(msg []byte) ([]byte, cryptotypes.PubKey, error) // implemented with a wrapper around keyring
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

	sig, pubKey, err := signingData.Signer(toSign)
	if err != nil {
		return nil, sdk.AccAddress{}, fmt.Errorf("sign: %w", err)
	}
	if pubKey == nil {
		return nil, sdk.AccAddress{}, errorsmod.Wrap(gerrc.ErrInvalidArgument, "signer returned nil pub key")
	}

	pubKeyAny, err := codectypes.NewAnyWithValue(pubKey)
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
