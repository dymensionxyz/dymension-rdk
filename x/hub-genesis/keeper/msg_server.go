package keeper

import (
	"context"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/dymensionxyz/dymension-rdk/utils/whitelistedrelayer"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

const port = transfertypes.PortID

var _ types.MsgServer = msgServer{}

type msgServer struct{ Keeper }

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) SendTransfer(goCtx context.Context, msg *types.MsgSendTransfer) (*types.MsgSendTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	relayer := msg.GetSigners()[0].String() // guaranteed length 1
	err := m.SendGenesisTransfer(ctx, relayer, msg.ChannelId)
	if err != nil {
		return nil, err
	}
	return &types.MsgSendTransferResponse{}, nil
}

func (k Keeper) SendGenesisTransfer(ctx sdk.Context, relayer, channelID string) error {
	wl, err := whitelistedrelayer.GetList(ctx, k.dk, k.sk)
	if err != nil {
		return errorsmod.Wrap(err, "get whitelisted relayers")
	}
	if !wl.Has(relayer) {
		return gerrc.ErrPermissionDenied.Wrap("not whitelisted")
	}
	state := k.GetState(ctx)
	if !state.CanonicalHubTransferChannelHasBeenSet() {
		state.SetCanonicalTransferChannel(port, channelID)
		k.SetState(ctx, state)
	}

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

	if data.GenesisTransfer != nil {
		// As we don't use the `ibc/transfer` module, we need to handle the funds escrow ourselves
		err = w.EscrowGenesisTransferFunds(ctx, port, channelID, data.GenesisInfo.BaseCoinSupply())
		if err != nil {
			return errorsmod.Wrap(err, "escrow genesis transfer funds")
		}
	}

	bz, err := json.Marshal(data)
	if err != nil {
		return errorsmod.Wrap(err, "marshal genesis bridge data")
	}

	timeoutTimestamp := ctx.BlockTime().Add(transferTimeout).UnixNano()
	_, err = w.channelKeeper.SendPacket(ctx, chanCap, port, channelID, clienttypes.ZeroHeight(), uint64(timeoutTimestamp), bz)
	return err
}
