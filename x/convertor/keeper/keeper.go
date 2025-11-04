package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transferkeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"github.com/dymensionxyz/dymension-rdk/x/convertor/types"
)

// Keeper wraps the Evmos IBC transfer keeper to perform decimal conversion
// before tokens are moved to escrow.
//
// The Evmos keeper is embedded, so all its methods are automatically available.
// Only the Transfer method is overridden to add decimal conversion logic.
type Keeper struct {
	transferkeeper.Keeper
	transferOverride types.TransferKeeper
	hubKeeper        types.HubKeeper
	bankKeeper       types.BankKeeper
}

// NewTransferKeeper creates a new TransferKeeper wrapper around the Evmos transfer keeper.
func NewTransferKeeper(
	transferKeeper transferkeeper.Keeper,
	transferOverride types.TransferKeeper,
	hubKeeper types.HubKeeper,
	bankKeeper types.BankKeeper,
) Keeper {
	return Keeper{
		Keeper:           transferKeeper,
		transferOverride: transferOverride,
		hubKeeper:        hubKeeper,
		bankKeeper:       bankKeeper,
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
		return w.transferOverride.Transfer(goCtx, msg)
	}

	// Parse sender address
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "invalid sender address")
	}

	// Convert the coin from rollapp token (18 decimals) to bridge token (custom decimals)
	convertedAmt, err := w.ConvertToBridgeAmt(ctx, msg.Token.Amount)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "convert coin to bridge token")
	}

	// burn the original tokens from the sender
	delta := sdk.NewCoin(msg.Token.Denom, msg.Token.Amount.Sub(convertedAmt))
	if err := w.BurnCoins(ctx, sender, delta); err != nil {
		return nil, errorsmod.Wrapf(err, "burn original tokens from sender")
	}

	// Create a new message with the converted token
	convertedMsg := &transfertypes.MsgTransfer{
		SourcePort:       msg.SourcePort,
		SourceChannel:    msg.SourceChannel,
		Token:            sdk.NewCoin(msg.Token.Denom, convertedAmt),
		Sender:           msg.Sender,
		Receiver:         msg.Receiver,
		TimeoutHeight:    msg.TimeoutHeight,
		TimeoutTimestamp: msg.TimeoutTimestamp,
		Memo:             msg.Memo,
	}

	// Call the underlying transfer keeper with the converted message
	return w.Keeper.Transfer(goCtx, convertedMsg)
}
