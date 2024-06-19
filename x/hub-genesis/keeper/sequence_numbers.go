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

// NOTE: assumes monotically increasing
func (k Keeper) saveLastSequenceNumber(ctx sdk.Context, seq uint64) {
	seqBz := make([]byte, 8)
	binary.BigEndian.PutUint64(seqBz, seq)
	ctx.KVStore(k.storeKey).Set([]byte("seqnum"), seqBz)
}

func (k Keeper) genesisIsFinished(ctx sdk.Context, port, channel string) bool {
	seqBz := ctx.KVStore(k.storeKey).Get([]byte("seqnum"))
	seq := binary.BigEndian.Uint64(seqBz)
	bz := k.channelKeeper.GetPacketCommitment(ctx, port, channel, seq)
	if len(bz) == 0 {
		// TODO: need to loop through all of them and save a convenience value
	}
	return false
}
