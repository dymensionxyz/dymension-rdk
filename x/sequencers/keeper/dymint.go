package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	tmcrypto "github.com/tendermint/tendermint/crypto/encoding"
)

// set dymint sequencers from InitChain
func (k Keeper) SetDymintSequencers(ctx sdk.Context, validators []abci.ValidatorUpdate) {
	if len(validators) > 1 {
		panic("more than one sequencer is not supported")
	}
	val := validators[0]

	tmkey, err := tmcrypto.PubKeyFromProto(val.PubKey)
	if err != nil {
		panic(err)
	}
	pubKey, err := cryptocodec.FromTmPubKeyInterface(tmkey)
	if err != nil {
		panic(err)
	}

	sequencer, err := types.NewSequencer(sdk.ValAddress(types.GenesisOperatorAddrStub), pubKey, uint64(val.Power))
	if err != nil {
		panic(err)
	}

	k.SetValidator(ctx, sequencer)
}
