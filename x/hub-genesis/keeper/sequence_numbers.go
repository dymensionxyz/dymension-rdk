package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

/*
Whenever a genesis transfer is sent, we record the sequence number. We do not allow transfers until
all acks have been received with success.
*/

var UnackedTransferSeqNumsPrefix = []byte("unacked_seqs")

func seqNumKey(seq uint64) []byte {
	bz := make([]byte, len(UnackedTransferSeqNumsPrefix))
	copy(bz, UnackedTransferSeqNumsPrefix)
	bz = append(bz, sdk.Uint64ToBigEndian(seq)...)
	return bz
}

func (k Keeper) saveUnackedTransferSeqNum(ctx sdk.Context, seq uint64) {
	ctx.KVStore(k.storeKey).Set(seqNumKey(seq), []byte{})
}

func (k Keeper) delUnackedTransferSeqNum(ctx sdk.Context, seq uint64) {
	ctx.KVStore(k.storeKey).Delete(seqNumKey(seq))
}

// returns all seq nums, only intended for genesis export
func (k Keeper) getAllUnackedTransferSeqNums(ctx sdk.Context) []uint64 {
	state := k.GetState(ctx)
	n := state.NumUnackedTransfers
	ret := make([]uint64, n)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, UnackedTransferSeqNumsPrefix)
	defer iterator.Close() // nolint: errcheck
	for ; iterator.Valid(); iterator.Next() {
		ret = append(ret, sdk.BigEndianToUint64(iterator.Key()))
	}
	return ret
}

// ackTransferSeqNum handles the inbound acknowledgement of an outbound genesis transfer
func (k Keeper) ackTransferSeqNum(ctx sdk.Context, seq uint64, success bool) {
	if !success {
		panic(fmt.Sprintf("genesis transfer unsuccessful seq: %d", seq))
	}
	k.delUnackedTransferSeqNum(ctx, seq)
	state := k.GetState(ctx)
	state.NumUnackedTransfers--
	if state.NumUnackedTransfers == 0 {
		// all acks have come back successfully
		state.OutboundTransfersEnabled = true
		ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeOutboundTransfersEnabled))
	}
	k.SetState(ctx, state)
}

func (k Keeper) outboundTransfersEnabled(ctx sdk.Context) bool {
	k.Logger(ctx).With("module", types.ModuleName).Debug("outbound transfers enabled")
	state := k.GetState(ctx)
	return state.OutboundTransfersEnabled
}
