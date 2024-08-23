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
	if _, err := m.AccAddr(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "acc addr")
	}
	if m.Operator != m.GetPayload().GetOperatorAddr() {
		return errorsmod.Wrap(gerrc.ErrInvalidArgument, "signer operator must match payload operator")
	}
	if err := m.KeyAndSig.Valid(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "key and sig")
	}
	return nil
}

func (m *MsgCreateSequencer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustAccAddr()}
}

func (m *MsgCreateSequencer) AccAddr() (sdk.AccAddress, error) {
	oper, err := m.OperatorAddr()
	if err != nil {
		return nil, errorsmod.Wrap(err, "operator addr")
	}
	return sdk.AccAddress(oper), nil
}

func (m *MsgCreateSequencer) MustAccAddr() sdk.AccAddress {
	return sdk.AccAddress(m.MustOperatorAddr())
}

func (m *MsgCreateSequencer) OperatorAddr() (sdk.ValAddress, error) {
	return sdk.ValAddressFromBech32(m.GetOperator())
}

func (m *MsgCreateSequencer) MustOperatorAddr() sdk.ValAddress {
	addr, err := m.OperatorAddr()
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *MsgUpdateSequencer) ValidateBasic() error {
	if _, err := m.AccAddr(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "acc addr")
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
	return []sdk.AccAddress{m.MustAccAddr()}
}

func (m *MsgUpdateSequencer) AccAddr() (sdk.AccAddress, error) {
	oper, err := m.OperatorAddr()
	if err != nil {
		return nil, errorsmod.Wrap(err, "operator addr")
	}
	return sdk.AccAddress(oper), nil
}

func (m *MsgUpdateSequencer) MustAccAddr() sdk.AccAddress {
	return sdk.AccAddress(m.MustOperatorAddr())
}

func (m *MsgUpdateSequencer) OperatorAddr() (sdk.ValAddress, error) {
	return sdk.ValAddressFromBech32(m.GetOperator())
}

func (m *MsgUpdateSequencer) MustOperatorAddr() sdk.ValAddress {
	addr, err := m.OperatorAddr()
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *MsgUpdateSequencer) RewardAcc() (sdk.AccAddress, error) {
	return sdk.AccAddressFromBech32(m.GetPayload().GetRewardAddr())
}

func (m *MsgUpdateSequencer) MustRewardAcc() sdk.AccAddress {
	s := m.GetPayload().GetRewardAddr()
	return sdk.MustAccAddressFromBech32(s)
}

func BuildMsgCreateSequencer(
	signingData SigningData,
	payload *CreateSequencerPayload,
) (*MsgCreateSequencer, error) {
	keyAndSig, err := createKeyAndSig(signingData, payload)
	if err != nil {
		return nil, fmt.Errorf("create key and sig: %w", err)
	}
	return &MsgCreateSequencer{
		Operator:  signingData.Operator.String(),
		KeyAndSig: keyAndSig,
		Payload:   payload,
	}, nil
}

func BuildMsgUpdateSequencer(
	signingData SigningData,
	payload *UpdateSequencerPayload,
) (*MsgUpdateSequencer, error) {
	keyAndSig, err := createKeyAndSig(signingData, payload)
	if err != nil {
		return nil, fmt.Errorf("create key and sig: %w", err)
	}
	return &MsgUpdateSequencer{
		Operator:  signingData.Operator.String(),
		KeyAndSig: keyAndSig,
		Payload:   payload,
	}, nil
}

// utility to create the key and sig argument needed in messages
func createKeyAndSig(signingData SigningData, payload codec.ProtoMarshaler) (*KeyAndSig, error) {
	toSign, err := CreateBytesToSign(signingData.ChainID, signingData.Account.GetAccountNumber(), payload)
	if err != nil {
		return nil, fmt.Errorf("create payload to sign: %w", err)
	}

	sig, pubKey, err := signingData.Signer(toSign)
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}
	if pubKey == nil {
		return nil, errorsmod.Wrap(gerrc.ErrInvalidArgument, "signer returned nil pub key")
	}

	pubKeyAny, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return nil, errorsmod.Wrap(err, "pubkey to any")
	}

	return &KeyAndSig{
		PubKey:    pubKeyAny,
		Signature: sig,
	}, nil
}
