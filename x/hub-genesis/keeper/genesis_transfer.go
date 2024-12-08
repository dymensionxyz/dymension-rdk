package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

const (
	// hubRecipient is the address of `x/rollapp` module's account on the hub chain.
	hubRecipient = "dym1mk7pw34ypusacm29m92zshgxee3yreums8avur"
)

// PrepareGenesisBridgeData prepares the genesis bridge data.
// Bridge data contains the genesis transfer packet data if the genesis accounts are defined, otherwise it's nil.
// Additionally, the method returns the packet coin (if any) that will be used for the escrow.
func (k Keeper) PrepareGenesisBridgeData(ctx sdk.Context) (types.GenesisBridgeData, error) {
	gInfo := k.GetGenesisInfo(ctx)

	denomMeta, ok := k.bk.GetDenomMetaData(ctx, gInfo.BaseDenom())
	if !ok {
		return types.GenesisBridgeData{}, errorsmod.Wrap(gerrc.ErrInternal, "denom metadata not found")
	}

	amount := math.ZeroInt()
	for _, acc := range gInfo.GenesisAccounts {
		amount = amount.Add(acc.Amount)
	}

	// no genesis accounts defined => no genesis transfer needed
	if amount.IsZero() {
		return types.GenesisBridgeData{
			GenesisInfo:     gInfo,
			NativeDenom:     denomMeta,
			GenesisTransfer: nil,
		}, nil
	}

	var (
		sender = k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress().String()
		denom  = gInfo.BaseDenom()
		packet = transfertypes.NewFungibleTokenPacketData(denom, amount.String(), sender, hubRecipient, "")
	)

	return types.GenesisBridgeData{
		GenesisInfo:     gInfo,
		NativeDenom:     denomMeta,
		GenesisTransfer: &packet,
	}, nil
}

// EscrowGenesisTransferFunds escrows the genesis transfer funds.
func (k Keeper) EscrowGenesisTransferFunds(ctx sdk.Context, portID, channelID string, token sdk.Coin) error {
	escrowAddress := transfertypes.GetEscrowAddress(portID, channelID)
	sender := k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress()
	return k.bk.SendCoins(ctx, sender, escrowAddress, sdk.NewCoins(token))
}

// enableBridge enables the bridge after successful genesis bridge phase.
func (k Keeper) enableBridge(ctx sdk.Context, state types.State, portID, channelID string) {
	state.SetCanonicalTransferChannel(portID, channelID)
	state.OutboundTransfersEnabled = true
	k.SetState(ctx, state)
}

// ResubmitPendingGenesisBridges attempts to resubmit genesis bridge data for all pending channels
func (k Keeper) ResubmitPendingGenesisBridges(ctx sdk.Context) {
	state := k.GetState(ctx)
	// If canonical channel is set, we don't need to resubmit
	if state.CanonicalHubTransferChannelHasBeenSet() {
		return
	}

	// Iterate over all pending channels
	err := k.PendingChannels.Walk(ctx, nil, func(portChannel string, retryRequired uint64) (stop bool, err error) {
		// Skip if channel is not failed yet
		if types.ChannelState(retryRequired) != types.Failed {
			return false, nil
		}

		portChan, err := types.FromPortAndChannelKey(portChannel)
		if err != nil {
			k.Logger(ctx).Error("invalid port/channel key", "portChannel", portChannel)
			return false, nil
		}

		seq, err := k.gb.SubmitGenesisBridgeData(ctx, portChan.Port, portChan.Channel)
		if err != nil {
			k.Logger(ctx).Error("failed to resubmit genesis bridge data", "port", portChan.Port, "channel", portChan.Channel, "error", err)
			return false, nil
		}

		// disable further retries
		err = k.SetPendingChannel(ctx, portChan, types.WaitingForAck)
		if err != nil {
			k.Logger(ctx).Error("failed to disable further retries", "port", portChan.Port, "channel", portChan.Channel, "error", err)
			return false, nil
		}

		k.Logger(ctx).Info("resubmitted genesis bridge data", "sequence", seq, "port", portChan.Port, "channel", portChan.Channel)
		return false, nil
	})
	if err != nil {
		k.Logger(ctx).Error("failed to resubmit genesis bridge data", "error", err)
	}
}
