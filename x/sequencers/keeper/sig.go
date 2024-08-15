package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

// a more easily testable helper for check sig
func checkSigAccNumber(ctx sdk.Context, acc uint64, keyAndSig *types.KeyAndSig, payloadApp codec.ProtoMarshaler) (bool, error) {
	pubKey, ok := keyAndSig.PubKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return false, errorsmod.WithType(errorsmod.Wrap(gerrc.ErrInvalidArgument, "could not assert cryptotypes pub key"), keyAndSig.PubKey.GetCachedValue())
	}

	payloadAppBz, err := payloadApp.Marshal()
	if err != nil {
		return false, err
	}

	payload := &types.PayloadToSign{
		PayloadApp:    payloadAppBz,
		ChainId:       ctx.ChainID(),
		AccountNumber: acc,
	}

	payloadBz, err := payload.Marshal()
	if err != nil {
		return false, err
	}

	ok = pubKey.VerifySignature(payloadBz, keyAndSig.GetSignature())
	return ok, nil
}
