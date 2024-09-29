package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

/*
Whenever a genesis transfer is sent, we record the sequence number. We do not allow transfers until
all acks have been received with success.

We use an O(1) access pattern, because we don't place a limit on the number of genesis accounts.
*/

var UnackedTransferSeqNumsPrefix = []byte("unacked_seqs")

func seqNumKey(seq uint64) []byte {
	bz := make([]byte, len(UnackedTransferSeqNumsPrefix))
	copy(bz, UnackedTransferSeqNumsPrefix)
	bz = append(bz, sdk.Uint64ToBigEndian(seq)...)
	return bz
}

func seqNumFromKey(key []byte) uint64 {
	return sdk.BigEndianToUint64(key[len(UnackedTransferSeqNumsPrefix):])
}

func (k Keeper) saveUnackedTransferSeqNum(ctx sdk.Context, seq uint64) {
	ctx.KVStore(k.storeKey).Set(seqNumKey(seq), []byte{})
}

func (k Keeper) delUnackedTransferSeqNum(ctx sdk.Context, seq uint64) {
	ctx.KVStore(k.storeKey).Delete(seqNumKey(seq))
}

func (k Keeper) hasUnackedTransferSeqNum(ctx sdk.Context, seq uint64) bool {
	return ctx.KVStore(k.storeKey).Has(seqNumKey(seq))
}

// returns all seq nums, only intended for genesis export
func (k Keeper) getAllUnackedTransferSeqNums(ctx sdk.Context) []uint64 {
	state := k.GetState(ctx)
	n := state.NumUnackedTransfers
	ret := make([]uint64, 0, n)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, UnackedTransferSeqNumsPrefix)
	defer iterator.Close() // nolint: errcheck
	for ; iterator.Valid(); iterator.Next() {
		ret = append(ret, seqNumFromKey(iterator.Key()))
	}
	return ret
}

// FIXME: refactor after proto change
// ackTransferSeqNum handles the inbound acknowledgement of an outbound genesis transfer
func (k Keeper) ackTransferSeqNum(ctx sdk.Context, seq uint64) {
	if !k.hasUnackedTransferSeqNum(ctx, seq) {
		panic("genesis transfer sequence number not found")
	}

	k.delUnackedTransferSeqNum(ctx, seq)
	state := k.GetState(ctx)
	state.NumUnackedTransfers--
	state.OutboundTransfersEnabled = true
	k.SetState(ctx, state)

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeOutboundTransfersEnabled))
	k.Logger(ctx).With("module", types.ModuleName).Debug("Enabled outbound transfers.")
}
