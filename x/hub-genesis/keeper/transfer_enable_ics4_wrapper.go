package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/tendermint/tendermint/libs/log"
)

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

// SendPacket is a wrapper around the ICS4Wrapper.SendPacket method.
// It will reject outbound transfers until the genesis bridge phase is finished.
func (w ICS4Wrapper) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (sequence uint64, err error) {
	if !w.k.IsBridgeOpen(ctx) {
		w.logger(ctx).Info("Transfer rejected: outbound transfers are disabled.")
		return 0, errorsmod.Wrap(gerrc.ErrFailedPrecondition, "genesis phase not finished")
	}

	return w.ICS4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}
