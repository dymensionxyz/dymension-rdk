package keeper

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

const (
	// hubRecipient is the address of `x/rollapp` module's account on the hub chain.
	hubRecipient = "dym1fuckme" // FIXME
)

type GenesisTransferMemo struct {
	GenesisTransferData GenesisTransferData `json:"genesis_transfer"`
}

type GenesisTransferData struct {
	GenesisAccounts []types.GenesisAccount `json:"genesis_accounts"`
}

func (k Keeper) PrepareGenesisTransfer(ctx sdk.Context, portID, channelID string) (*transfertypes.FungibleTokenPacketData, error) {
	state := k.GetState(ctx)
	amount := math.ZeroInt()
	for _, acc := range state.GenesisAccounts {
		amount = amount.Add(acc.Amount)
	}

	// no genesis accounts defined => no genesis transfer needed
	if amount.IsZero() {
		return nil, nil
	}

	sender := k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress().String()
	denom := k.GetBaseDenom(ctx)

	// prepare memo with the genesis accounts info
	memo, err := k.CreateGenesisAccountsMemo(ctx, state.GenesisAccounts)
	if err != nil {
		return nil, errorsmod.Wrap(err, "create genesis memo")
	}

	// As we don't use the `ibc/transfer` module, we need to handle the funds escrow ourselves
	err = k.EscrowGenesisTransferFunds(ctx, portID, channelID, sdk.NewCoin(denom, amount))
	if err != nil {
		return nil, errorsmod.Wrap(err, "escrow genesis transfer funds")
	}

	packet := transfertypes.NewFungibleTokenPacketData(denom, amount.String(), sender, hubRecipient, memo)
	return &packet, nil
}

// EscrowGenesisTransferFunds escrows the genesis transfer funds.
// The code is copied from the `transfer` module's `Keeper.sendTransfer` method.
func (k Keeper) EscrowGenesisTransferFunds(ctx sdk.Context, portID, channelID string, token sdk.Coin) error {
	escrowAddress := transfertypes.GetEscrowAddress(portID, channelID)
	sender := k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress()
	return k.bk.SendCoins(ctx, sender, escrowAddress, sdk.NewCoins(token))
}

// CreateGenesisAccountsMemo creates a memo with the genesis accounts info.
func (k Keeper) CreateGenesisAccountsMemo(ctx sdk.Context, genesisAccounts []types.GenesisAccount) (string, error) {
	memo := GenesisTransferMemo{
		GenesisTransferData: GenesisTransferData{
			GenesisAccounts: genesisAccounts,
		},
	}

	memoBytes, err := json.Marshal(memo)
	if err != nil {
		return "", err
	}

	return string(memoBytes), nil
}
