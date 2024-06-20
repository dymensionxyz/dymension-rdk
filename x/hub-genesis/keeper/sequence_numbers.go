package keeper

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hubgentypes "github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

func seqNumKey(port, channel string) []byte {
	return []byte(fmt.Sprintf("seqnum/%s/%s", port, channel))
}

// NOTE: assumes monotonically increasing
func (k Keeper) saveLastSequenceNumber(ctx sdk.Context, port, channel string, seq uint64) {
	if existing := k.getLastSequenceNumber(ctx, port, channel); seq < existing {
		/*
			Proof that this will never happen:
			All transfers
		*/
		panic(
			fmt.Sprintf(
				"%s: a higher sequence number has already been set: existing: %d: new: %d: port: %s, channel: %s", hubgentypes.ModuleName, existing, seq, port, channel))
	}
	seqBz := make([]byte, 8)
	binary.BigEndian.PutUint64(seqBz, seq)
	ctx.KVStore(k.storeKey).Set(seqNumKey(port, channel), seqBz)
}

// NOTE: assumes monotonically increasing
func (k Keeper) getLastSequenceNumber(ctx sdk.Context, port, channel string) uint64 {
	bz := ctx.KVStore(k.storeKey).Get(seqNumKey(port, channel))
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) genesisIsFinished(ctx sdk.Context, port, channel string) bool {
	for seq := range k.getLastSequenceNumber(ctx, port, channel) + 1 {
		bz := k.channelKeeper.GetPacketCommitment(ctx, port, channel, seq)
		if len(bz) != 0 {
			// there is still a sequence number for which we didn't get an ack back yet
			// so we still need to wait some more.
			return false
		}
	}
	return true
}
