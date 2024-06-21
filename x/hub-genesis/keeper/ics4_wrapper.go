package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/tendermint/tendermint/libs/log"
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

func (w ICS4Wrapper) logger(ctx sdk.Context) log.Logger {
	return w.k.Logger(ctx).With("module", types.ModuleName, "component", "ics4 middleware")
}

// SendPacket does two things:
//  1. It stops anyone from sending a packet with the special memo. Only the module itself is allowed to do so.
//  2. It stops anyone from sending a regular transfer until the genesis phase is finished. To help with this,
//     it tracks all acks which arrive from genesis transfers.
func (w ICS4Wrapper) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	l := w.logger(ctx)

	state := w.k.GetState(ctx)
	if !state.OutboundTransfersEnabled {
		l.Debug("Transfer rejected: outbound transfers are disabled.")
		return 0, errorsmod.Wrap(gerrc.ErrFailedPrecondition, "genesis phase not finished")
	}

	var transfer transfertypes.FungibleTokenPacketData
	_ = transfertypes.ModuleCdc.UnmarshalJSON(data, &transfer)

	saveSeq := false
	if memoHasKey(transfer.GetMemo()) && state.IsCanonicalHubTransferChannel(sourcePort, sourceChannel) {
		if !skipAuthorizationCheck(ctx) {
			return 0, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "cannot use transfer genesis memo")
		}
		// This is a genesis transfer, we record the sequence number.
		// Record the sequence number because we need to tick them off as they get acked.
		saveSeq = true
	}

	seq, err := w.ICS4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	if saveSeq && err == nil {
		w.k.saveUnackedTransferSeqNum(ctx, seq)
	}
	return seq, err
}
