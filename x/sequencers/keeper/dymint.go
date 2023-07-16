package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	tmcrypto "github.com/tendermint/tendermint/crypto/encoding"
)

// set dymint sequencers from InitChain
func (k Keeper) SetDymintSequencers(ctx sdk.Context, validators []abci.ValidatorUpdate) {
	for _, val := range validators {
		tmkey, err := tmcrypto.PubKeyFromProto(val.PubKey)
		if err != nil {
			panic(err)
		}
		pubKey, err := cryptocodec.FromTmPubKeyInterface(tmkey)
		if err != nil {
			panic(err)
		}
		k.SetDymintSequencerByAddr(ctx, sdk.ConsAddress(pubKey.Address()), uint64(val.Power))
	}
}

// get a single sequencer registered on dymint by consensus address
func (k Keeper) GetDymintSequencerByAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (power uint64, found bool) {
	store := ctx.KVStore(k.storeKey)
	powerByte := store.Get(types.GetDymintSeqKey(consAddr))
	if powerByte == nil {
		return 0, false
	}

	return binary.LittleEndian.Uint64(powerByte), true
}

// set a single sequencer registered on dymint by consensus address
func (k Keeper) SetDymintSequencerByAddr(ctx sdk.Context, consAddr sdk.ConsAddress, power uint64) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(power))

	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDymintSeqKey(consAddr), b)
	return nil
}
