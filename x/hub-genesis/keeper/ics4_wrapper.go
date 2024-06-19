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

// SendPacket prevents anyone from sending a packet with the memo
// The app should be wired to allow the middleware to circumvent this
// It also prevents anyone from sending a transfer before all the genesis transfers
// have had successull acks come back.
func (w ICS4Wrapper) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	/*
		Rethinking how this can work:
		- block transfers if the genesis accounts list is not empty, that means they weren't all sent
		- record the highest seq num, always
		- on a transfer attempt, check for any of those seq nums, save the result when it's good, to amortize
	*/

	if !w.k.genesisIsFinished(ctx) {
		return 0, errorsmod.Wrap(gerrc.ErrFailedPrecondition, "genesis phase not finished")
	}

	var transfer transfertypes.FungibleTokenPacketData
	_ = transfertypes.ModuleCdc.UnmarshalJSON(data, &transfer)

	if memoHasKey(transfer.GetMemo()) {
		if !skipAuthorizationCheck(ctx) {
			return 0, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "cannot use transfer genesis memo")
		}

		// record the sequence number because we need to tick them off as they get acked

		seq, err := w.ICS4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
		if err != nil {
			return seq, err
		}
		return seq, w.k.saveLastSequenceNumber(ctx, seq)
	}
	return w.ICS4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}
