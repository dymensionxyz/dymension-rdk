package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

const (
	// hubRecipient is the address of `x/rollapp` module's account on the hub chain.
	hubRecipient = "dym1mk7pw34ypusacm29m92zshgxee3yreums8avur"
)

// PrepareGenesisTransfer prepares the genesis transfer packet.
// It returns the packet data if the genesis accounts are defined, otherwise it returns nil.
// The transfer funds are escrowed explicitly in this method.
// A memo is attaached with the genesis accounts info, to be validated against the genesis accounts defined on the hub chain.
func (k Keeper) PrepareGenesisTransfer(ctx sdk.Context, portID, channelID string) (*transfertypes.FungibleTokenPacketData, error) {
	gAccounts := k.GetGenesisInfo(ctx).GenesisAccounts
	amount := math.ZeroInt()
	for _, acc := range gAccounts {
		amount = amount.Add(acc.Amount)
	}

	// no genesis accounts defined => no genesis transfer needed
	if amount.IsZero() {
		return nil, nil
	}

	sender := k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress().String()
	denom := k.GetBaseDenom(ctx)

	// As we don't use the `ibc/transfer` module, we need to handle the funds escrow ourselves
	err := k.EscrowGenesisTransferFunds(ctx, portID, channelID, sdk.NewCoin(denom, amount))
	if err != nil {
		return nil, errorsmod.Wrap(err, "escrow genesis transfer funds")
	}

	packet := transfertypes.NewFungibleTokenPacketData(denom, amount.String(), sender, hubRecipient, "")
	return &packet, nil
}

// EscrowGenesisTransferFunds escrows the genesis transfer funds.
// The code is copied from the `transfer` module's `Keeper.sendTransfer` method.
func (k Keeper) EscrowGenesisTransferFunds(ctx sdk.Context, portID, channelID string, token sdk.Coin) error {
	escrowAddress := transfertypes.GetEscrowAddress(portID, channelID)
	sender := k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress()
	return k.bk.SendCoins(ctx, sender, escrowAddress, sdk.NewCoins(token))
}
