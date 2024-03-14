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
	if len(sequencers) > 1 {
		panic(types.ErrMultipleDymintSequencers)
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

	sequencer, err := types.NewSequencer(sdk.ValAddress(types.GenesisOperatorAddrStub), pubKey, uint64(seq.Power))
	if err != nil {
		panic(err)
	}

	k.SetSequencer(ctx, sequencer)
}
