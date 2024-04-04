package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	tmcrypto "github.com/tendermint/tendermint/crypto/encoding"
)

// SetDymintSequencers sets the sequencers set by dymint on InitChain.
// As currently we're using the abci InitChain, there are 2 obstacles we need to face, unlike when creating a validator:
// 1. InitChain expected the validatorUpdate it gets in return to be the same as it sends
// 2. We need someway to set the operator address, which is the address that will be used for rewards for the sequencer.
// To overcome those obstacles, we do the following:
//  1. Upon InitChain, call SetDymintSequencers and create a dummy sequencer object with the consensus pubkey and power we get from the validatorUpdate.
//  2. Afterwards, upon initGenesis, we build a validator-like object where the operator address is the one we set in the genesis file and the
//     consensus pubkey and power are the ones we set in the dummy sequencer object.
//
// At the end we delete the dummy sequencer object.
func (k Keeper) SetDymintSequencers(ctx sdk.Context, sequencers []abci.ValidatorUpdate) {
	if len := len(sequencers); len != 1 {
		switch len {
		case 0:
			panic(types.ErrNoSequencerOnInitChain)
		default:
			panic(types.ErrMultipleDymintSequencers)
		}
	}
	seq := sequencers[0]

	tmkey, err := tmcrypto.PubKeyFromProto(seq.PubKey)
	if err != nil {
		panic(err)
	}
	pubKey, err := cryptocodec.FromTmPubKeyInterface(tmkey)
	if err != nil {
		panic(err)
	}

	// On InitChain we only have the consesnsus pubkey, so we set the operator address to a dummy value
	sequencer, err := types.NewSequencer(sdk.ValAddress(types.InitChainStubAddr), pubKey, seq.Power)
	if err != nil {
		panic(err)
	}

	k.SetSequencer(ctx, sequencer)
	err = k.SetSequencerByConsAddr(ctx, sequencer)
	if err != nil {
		panic(err)
	}
}
