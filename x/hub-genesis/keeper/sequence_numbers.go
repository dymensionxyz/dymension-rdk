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

func seqNumKey(port, channel string, seq uint64) []byte {
	bz := []byte(fmt.Sprintf("seqnumval/%s/%s/", port, channel))
	bz = append(bz, sdk.Uint64ToBigEndian(seq)...)
	return bz
}

func (k Keeper) saveSeqNum(ctx sdk.Context, port, channel string, seq uint64) {
	ctx.KVStore(k.storeKey).Set(seqNumKey(port, channel, seq), []byte{})
}

func (k Keeper) delSeqNum(ctx sdk.Context, port, channel string, seq uint64) {
	ctx.KVStore(k.storeKey).Delete(seqNumKey(port, channel, seq))
}

// ackSeqNum handles the inbound acknowledgement of an outbound genesis transfer
func (k Keeper) ackSeqNum(ctx sdk.Context, port, channel string, seq uint64, success bool) {
	if !success {
		panic(fmt.Sprintf("genesis transfer unsuccessful, port: %s, channel: %s: seq: %d", port, channel, seq))
	}
	k.delSeqNum(ctx, port, channel, seq)
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
