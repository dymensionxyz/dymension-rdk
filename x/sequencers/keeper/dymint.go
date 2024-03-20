package keeper

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	tmcrypto "github.com/tendermint/tendermint/crypto/encoding"
)

// set dymint sequencers from InitChain
func (k Keeper) SetDymintSequencers(ctx sdk.Context, sequencers []abci.ValidatorUpdate) {
	if len := len(sequencers); len != 2 {
		switch len {
		case 0:
			panic(types.ErrNoSequencerOnInitChain)
		case 1:
			panic(types.ErrMissingOperatorAddrsOnInitChain)
		default:
			panic(types.ErrMultipleDymintSequencers)
		}
	}

	var (
		operatorAddr    sdk.ValAddress
		power           int64
		consensusPubKey cryptotypes.PubKey
		operatorPubkey  cryptotypes.PubKey
	)

	for _, seq := range sequencers {
		tmkey, err := tmcrypto.PubKeyFromProto(seq.PubKey)
		if err != nil {
			panic(err)
		}
		pubKey, err := cryptocodec.FromTmPubKeyInterface(tmkey)
		if err != nil {
			panic(err)
		}

		if pubKey.Type() == ed25519.KeyType {
			consensusPubKey = pubKey
			power = seq.Power
		} else {
			operatorPubkey = pubKey
			operatorAddr = sdk.ValAddress(pubKey.Address())
		}
	}

	if operatorAddr.Empty() || consensusPubKey == nil {
		panic(types.ErrFailedInitChain)
	}

	sequencer, err := types.NewSequencer(operatorAddr, consensusPubKey, power)
	if err != nil {
		panic(err)
	}
	k.SetSequencer(ctx, sequencer)
	err = k.SetSequencerByConsAddr(ctx, sequencer)
	if err != nil {
		panic(err)
	}

	// Required code as the cosmos sdk validates that the InitChain request is equal to the result
	// we pass this dummy sequencer to the InitGenesis function so it could be returned in the ValidatorUpdate response
	dummySequencer, err := types.NewSequencer(sdk.ValAddress(types.InitChainStubAddr), operatorPubkey, power)
	if err != nil {
		panic(err)
	}
	k.SetSequencer(ctx, dummySequencer)
}
