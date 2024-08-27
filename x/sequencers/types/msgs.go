package types

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
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
)

func (m *MsgCreateSequencer) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return unpacker.UnpackAny(m.PubKey, new(cryptotypes.PubKey))
}

func (m *MsgCreateSequencer) ValidateBasic() error {
	if _, err := m.OperatorAccAddr(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "acc addr")
	}
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
	operator, err := m.OperatorAddr()
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "acc addr")
	}
	pubKey, _ := v.ConsPubKey()

	// We return OK only if the key and sig contains a key and signature where the signature was produced by the key, and the signature
	// is over the operator account.
	//
	// The reasoning is as follows:
	//
	// We know from the SDK TX signing mechanism that the account originates from the operator, and on this chain ID.
	// Therefore, we just need to check that the consensus private key also over this operator. Then we know that
	// the private key holder of the operator and the consensus keys is the same actor.
	if !pubKey.VerifySignature(operator, m.GetSignature()) {
		return errorsmod.Wrap(gerrc.ErrUnauthenticated, "priv key of pub cons key was not used to sign operator addr")
	}

	return nil
}

// Validator is a convenience method - it returns a validator object which already
// has implementations of various useful methods like obtaining various type conversions
// for the public key.
func (m *MsgCreateSequencer) Validator() stakingtypes.Validator {
	return stakingtypes.Validator{ConsensusPubkey: m.PubKey, OperatorAddress: m.MustOperatorAddr().String()}
}

func (m *MsgCreateSequencer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustAccAddr()}
}

func (m *MsgCreateSequencer) OperatorAccAddr() (sdk.AccAddress, error) {
	operator, err := m.OperatorAddr()
	if err != nil {
		return nil, errorsmod.Wrap(err, "operator addr")
	}
	return sdk.AccAddress(operator), nil
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

func BuildMsgCreateSequencer(
	// Produces a signature over a message
	signer func(msg []byte) ([]byte, cryptotypes.PubKey, error), // implemented with a wrapper around keyring
	// Operator, will be signed over by signer
	operator sdk.ValAddress,
) (*MsgCreateSequencer, error) {
	sig, pubKey, err := signer(operator)
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
	return &MsgCreateSequencer{
		Operator:  operator.String(),
		PubKey:    pubKeyAny,
		Signature: sig,
	}, nil
}

func (m *MsgUpdateSequencer) ValidateBasic() error {
	if _, err := m.OperatorAccAddr(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "acc addr")
	}
	if _, err := m.RewardAcc(); err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "reward addr")
	}
	return nil
}

func (m *MsgUpdateSequencer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustAccAddr()}
}

func (m *MsgUpdateSequencer) OperatorAccAddr() (sdk.AccAddress, error) {
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
	return sdk.AccAddressFromBech32(m.GetRewardAddr())
}

func (m *MsgUpdateSequencer) MustRewardAcc() sdk.AccAddress {
	ret, err := m.RewardAcc()
	if err != nil {
		panic(err)
	}
	return ret
}
