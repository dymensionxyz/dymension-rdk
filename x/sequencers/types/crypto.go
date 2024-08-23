package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

type CreatorAccount interface {
	GetAccountNumber() uint64
}

type SigningData struct {
	Operator sdk.ValAddress
	Account  CreatorAccount
	ChainID  string
	Signer   func(msg []byte) ([]byte, cryptotypes.PubKey, error) // implemented with a wrapper around keyring
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

// Validator is a convenience method - it returns a validator object which already
// has implementations of various useful methods like obtaining various type conversions
// for the public key.
func (m *KeyAndSig) Validator() stakingtypes.Validator {
	return stakingtypes.Validator{ConsensusPubkey: m.PubKey}
}

// Ok return nil only if the key and sig contains a key and signature where the signature was produced by the key, and the signature
// is over the account from the provided address, and the app payload data.
//
// The reasoning is as follows:
// We know that the TX containing the Msg was signed by addr, because it has passed the sdk signature verification ante.
// Therefore, if we require that the private key for the consensus address was used to sign off over this addr AND this chain ID then
// we know that the owner of the private key really intended this payload to be included in this transaction, and it is not man in the middle or replay.
func (m *KeyAndSig) Ok(ctx sdk.Context, acc auth.AccountI, payloadApp codec.ProtoMarshaler) error {
	v := m.Validator()

	/*
		A simpler logic (1):
		We have the operator signature over chain ID and account number
		We have cons key signature over the operator signature
		Therefore the cons key owner did intend to use this operator

		A simpler logic (2):
		We have the operator signature over chain ID and account number
		We have cons key signature over the chain ID and account number
		Therefore implicitly we have acceptance of the operator addr too
		Therefore the cons key owner did intend to use this operator
	*/

	payloadBz, err := CreateBytesToSign(
		ctx.ChainID(),
		acc.GetAccountNumber(),
		payloadApp,
	)
	if err != nil {
		return errorsmod.Wrap(err, "create bytes to sign")
	}

	pubKey, err := v.ConsPubKey()
	if err != nil {
		return errorsmod.Wrap(err, "get cons pubkey")
	}

	ok := pubKey.VerifySignature(payloadBz, m.GetSignature())

	if !ok {
		return gerrc.ErrUnauthenticated
	}

	return nil
}
