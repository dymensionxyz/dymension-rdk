package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	tmcrypto "github.com/tendermint/tendermint/crypto/encoding"
)

// set dymint sequencers from InitChain
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
