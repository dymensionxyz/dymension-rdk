package keeper

import (
	errorsmod "cosmossdk.io/errors"
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

	amount := gInfo.Amt()

	// no genesis accounts defined => no genesis transfer needed
	if gInfo.Amt().IsZero() {
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
// The code is copied from the `transfer` module's `Keeper.sendTransfer` method.
func (k Keeper) EscrowGenesisTransferFunds(ctx sdk.Context, portID, channelID string, token sdk.Coin) error {
	escrowAddress := transfertypes.GetEscrowAddress(portID, channelID)
	sender := k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress()
	return k.bk.SendCoins(ctx, sender, escrowAddress, sdk.NewCoins(token))
}
