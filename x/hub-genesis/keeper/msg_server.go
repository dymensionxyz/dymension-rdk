package keeper

import (
	"context"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

const port = transfertypes.PortID

var _ types.MsgServer = msgServer{}

type msgServer struct{ Keeper }

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// note that this is only allowed by the whitelisted relayer (enforced in ante)
func (m msgServer) SendTransfer(goCtx context.Context, msg *types.MsgSendTransfer) (*types.MsgSendTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := m.SendGenesisTransfer(ctx, msg.ChannelId)
	if err != nil {
		return nil, err
	}
	return &types.MsgSendTransferResponse{}, nil
}

const expectedChan = "channel-0" // tokenless only

func (k Keeper) SendGenesisTransfer(ctx sdk.Context, channelID string) error {
	if k.Tokenless(ctx) && channelID != expectedChan {
		return gerrc.ErrInvalidArgument.Wrapf("tokenless chain: wrong channel id, expect: %s", expectedChan)
	}
	state := k.GetState(ctx)
	if state.InFlight {
		return gerrc.ErrFailedPrecondition.Wrap("sent transfer is already in flight")
	}
	if state.OutboundTransfersEnabled {
		return gerrc.ErrInvalidArgument.Wrap("bridge already open")
	}
	c, ok := k.channelKeeper.GetChannel(ctx, port, channelID)
	if !ok {
		return gerrc.ErrNotFound.Wrap("channel")
	}
	if c.State != channeltypes.OPEN {
		return gerrc.ErrFailedPrecondition.Wrap("channel not open")
	}
	state.SetCanonicalTransferChannel(port, channelID)
	state.InFlight = true
	k.SetState(ctx, state)

	if err := k.SubmitGenesisBridgeData(ctx, channelID); err != nil {
		return errorsmod.Wrap(err, "submit genesis bridge data")
	}
	return nil
}

// SubmitGenesisBridgeData sends the genesis bridge data over the channel.
// The genesis bridge data includes the genesis info, the native denom metadata, and the genesis transfer packet.
// It uses the channel keeper to send the packet, instead of transfer keeper, as we are not sending fungible token directly.
func (w Keeper) SubmitGenesisBridgeData(ctx sdk.Context, channelID string) (err error) {
	_, chanCap, err := w.channelKeeper.LookupModuleByChannel(ctx, port, channelID)
	if err != nil {
		return
	}

	data, err := w.PrepareGenesisBridgeData(ctx)
	if err != nil {
		return errorsmod.Wrap(err, "prepare genesis bridge data")
	}

	bz, err := json.Marshal(data)
	if err != nil {
		return errorsmod.Wrap(err, "marshal genesis bridge data")
	}

	timeoutTimestamp := ctx.BlockTime().Add(transferTimeout).UnixNano()
	_, err = w.channelKeeper.SendPacket(ctx, chanCap, port, channelID, clienttypes.ZeroHeight(), uint64(timeoutTimestamp), bz)
	return err
}
