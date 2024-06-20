package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

type ctxKeySkip struct{}

// skipAuthorizationCheckContext returns a context which can be passed to ibc SendPacket
// if passed, the memo guard will not check that call
func skipAuthorizationCheckContext(ctx sdk.Context) sdk.Context {
	return ctx.WithValue(ctxKeySkip{}, true)
}

func skipAuthorizationCheck(ctx sdk.Context) bool {
	val, ok := ctx.Value(ctxKeySkip{}).(bool)
	return ok && val
}

type ICS4Wrapper struct {
	porttypes.ICS4Wrapper
	k Keeper
}

func NewICS4Wrapper(next porttypes.ICS4Wrapper, k Keeper) *ICS4Wrapper {
	return &ICS4Wrapper{next, k}
}

// SendPacket does two things:
// 1. It stops anyone from sending a packet with the special memo. Only the module itself is allowed to do so.
// 2. It stops anyone from sending a regular transfer until the genesis phase is finished.
func (w ICS4Wrapper) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	var transfer transfertypes.FungibleTokenPacketData
	_ = transfertypes.ModuleCdc.UnmarshalJSON(data, &transfer)

	if memoHasKey(transfer.GetMemo()) {
		if !skipAuthorizationCheck(ctx) {
			return 0, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "cannot use transfer genesis memo")
		}

		// This is a genesis transfer, we record the sequence number.
		// Record the sequence number because we need to tick them off as they get acked.

		seq, err := w.ICS4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
		if err != nil {
			return seq, err
		}
		w.k.saveSeqNum(ctx, sourcePort, sourceChannel, seq)
		return seq, nil
	} else if !w.k.genesisIsFinished(ctx, sourcePort, sourceChannel) {
		return 0, errorsmod.Wrap(gerrc.ErrFailedPrecondition, "genesis phase not finished")
	}
	return w.ICS4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}
