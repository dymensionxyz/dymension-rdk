package keeper

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func seqNumKey(port, channel string) []byte {
	return []byte(fmt.Sprintf("seqnum/%s/%s", port, channel))
}

// NOTE: assumes monotonically increasing
func (k Keeper) saveLastSequenceNumber(ctx sdk.Context, port, channel string, seq uint64) {
	seqBz := make([]byte, 8)
	binary.BigEndian.PutUint64(seqBz, seq)
	ctx.KVStore(k.storeKey).Set(seqNumKey(port, channel), seqBz)
}

func (k Keeper) getLastSequenceNumber(ctx sdk.Context, port, channel string) uint64 {
	bz := ctx.KVStore(k.storeKey).Get(seqNumKey(port, channel))
	return sdk.BigEndianToUint64(bz)
}

// genesisIsFinished returns if the genesis bridge protocol phase is finished. It is finished
// when all genesis transfers sent from the RA to the Hub have been acked. After this you're
// allowed to send regular transfers. The first regular transfer received on the Hub marks
// the end of the protocol from the Hub's perspective.
func (k Keeper) genesisIsFinished(ctx sdk.Context, port, channel string) bool {
	state := k.GetState(ctx)
	if state.GetFinished() {
		return true
	}
	// This operation may not be super cheap, but once the genesis phase is finished, it won't be necessary.
	// Much simpler than using a map to check off each seq num.
	for seq := range k.getLastSequenceNumber(ctx, port, channel) + 1 {
		bz := k.channelKeeper.GetPacketCommitment(ctx, port, channel, seq)
		if len(bz) != 0 {
			// there is still a sequence number for which we didn't get an ack back yet
			// so we still need to wait some more.
			return false
		}
	}
	state.Finished = true
	k.SetState(ctx, state)
	return true
}
