package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	evmostransferkeeper "github.com/evmos/evmos/v12/x/ibc/transfer/keeper"

	"github.com/dymensionxyz/dymension-rdk/x/convertor/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

// Keeper wraps the Evmos IBC transfer keeper to perform decimal conversion
// before tokens are moved to escrow.
//
// The Evmos keeper is embedded, so all its methods are automatically available.
// Only the Transfer method is overridden to add decimal conversion logic.
type Keeper struct {
	evmostransferkeeper.Keeper
	hubKeeper  types.HubKeeper
	bankKeeper types.BankKeeper
}

// NewTransferKeeper creates a new TransferKeeper wrapper around the Evmos transfer keeper.
func NewTransferKeeper(
	transferKeeper evmostransferkeeper.Keeper,
	hubKeeper types.HubKeeper,
	bankKeeper types.BankKeeper,
) Keeper {
	return Keeper{
		Keeper:     transferKeeper,
		hubKeeper:  hubKeeper,
		bankKeeper: bankKeeper,
	}
}

// Transfer overrides the transfer keeper's Transfer method to perform decimal conversion
// before the tokens are moved to escrow.
func (w Keeper) Transfer(
	goCtx context.Context,
	msg *transfertypes.MsgTransfer,
) (*transfertypes.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if there's a decimal conversion required for this denom
	required, err := w.ConversionRequired(ctx, msg.Token.Denom)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "check if conversion required")
	}

	// If no conversion is needed, pass through to the underlying keeper
	if !required {
		return w.Keeper.Transfer(goCtx, msg)
	}

	// Make sure we're not trying to send the bridge denom itself
	pair, err := w.hubKeeper.GetDecimalConversionPair(ctx)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "get decimal conversion pair")
	}
	if pair.FromToken == msg.Token.Denom {
		return nil, errorsmod.Wrapf(gerrc.ErrInvalidArgument, "cannot send bridge denom itself")
	}

	// Parse sender address
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "invalid sender address")
	}

	// Convert the coin from rollapp token (18 decimals) to bridge token (custom decimals)
	convertedCoin, err := w.ConvertToBridgeCoin(ctx, msg.Token)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "convert coin to bridge token")
	}

	// Check if there's any truncation (remainder after conversion)
	// When converting from 18 decimals to lower decimals, we might lose precision
	reconvertedCoin, err := w.ConvertFromBridgeCoin(ctx, convertedCoin)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "reconvert coin to check truncation")
	}

	// Burn the original tokens from the sender (since we'll be sending converted tokens instead)
	if err := w.BurnCoins(ctx, sender, reconvertedCoin); err != nil {
		return nil, errorsmod.Wrapf(err, "burn original tokens from sender")
	}

	// Mint the converted tokens to the sender (so the transfer keeper can then move them to escrow)
	if err := w.MintCoins(ctx, sender, convertedCoin); err != nil {
		return nil, errorsmod.Wrapf(err, "mint converted tokens to sender")
	}

	// Create a new message with the converted token
	convertedMsg := &transfertypes.MsgTransfer{
		SourcePort:       msg.SourcePort,
		SourceChannel:    msg.SourceChannel,
		Token:            convertedCoin,
		Sender:           msg.Sender,
		Receiver:         msg.Receiver,
		TimeoutHeight:    msg.TimeoutHeight,
		TimeoutTimestamp: msg.TimeoutTimestamp,
		Memo:             msg.Memo,
	}

	// Call the underlying transfer keeper with the converted message
	return w.Keeper.Transfer(goCtx, convertedMsg)
}
