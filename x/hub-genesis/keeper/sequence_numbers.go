package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func seqNumKey(seq uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, seq)

	var key []byte
	key = append(key, []byte("seqnums/")...)
	key = append(key, bz...)
	return key
}

func seqCntKey() []byte {
	var key []byte
	key = append(key, []byte("seqcnt")...)
	return key
}

func (k Keeper) saveSequenceNumber(ctx sdk.Context, seq uint64) error {
	store := ctx.KVStore(k.storeKey)
	ctx.KVStore(k.storeKey).Set(seqNumKey(seq), []byte{})
	cntBz := store.Get(seqCntKey())
	cnt := binary.BigEndian.Uint64(cntBz)
	cnt++
	if !ok
}

func (k Keeper) delSequenceNumber(ctx sdk.Context, seq uint64) error {
	store := ctx.KVStore(k.storeKey).
		store.Set(seqNumKey(seq), nil)
}

func (k Keeper) genesisIsFinished(ctx sdk.Context) bool {
}
