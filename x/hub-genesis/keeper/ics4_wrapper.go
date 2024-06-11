package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
)

type ICS4Wrapper struct {
	porttypes.ICS4Wrapper
}

func NewICS4Wrapper(next porttypes.ICS4Wrapper) *ICS4Wrapper {
	return &ICS4Wrapper{next}
}

// SendPacket prevents anyone from sending a packet with the memo
// The app should be wired to allow the middleware to circumvent this
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
		return 0, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "cannot use transfer genesis memo")
	}
	return w.ICS4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}
