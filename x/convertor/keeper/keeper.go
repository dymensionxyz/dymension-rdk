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
	transferStack types.TransferKeeper // allows to have transfer stack (e.g to support erc20 middleware)
	hubKeeper     types.HubKeeper
	bankKeeper    types.BankKeeper
}

// NewTransferKeeper creates a new TransferKeeper wrapper around the Evmos transfer keeper.
func NewTransferKeeper(
	transferKeeper transferkeeper.Keeper,
	transferStack types.TransferKeeper,
	hubKeeper types.HubKeeper,
	bankKeeper types.BankKeeper,
) Keeper {

	if transferStack == nil {
		transferStack = transferKeeper
	}

	return Keeper{
		Keeper:        transferKeeper,
		transferStack: transferStack,
		hubKeeper:     hubKeeper,
		bankKeeper:    bankKeeper,
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
		return w.transferStack.Transfer(goCtx, msg)
	}

	// Parse sender address
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "invalid sender address")
	}

	// Log the balance before conversion
	balanceBefore := w.bankKeeper.GetBalance(ctx, sender, msg.Token.Denom)
	ctx.Logger().Info("Account balance before conversion",
		"sender", msg.Sender,
		"denom", msg.Token.Denom,
		"balance", balanceBefore.Amount.String(),
	)

	// Convert the coin from rollapp token (18 decimals) to bridge token (custom decimals)
	convertedAmt, err := w.ConvertToBridgeAmt(ctx, msg.Token.Amount)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "convert coin to bridge token")
	}

	// Convert the bridge amount back to rollapp decimals to calculate the actual amount to transfer
	convertedAmtInRollappDecimals, err := w.ConvertFromBridgeAmt(ctx, convertedAmt)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "convert bridge amount back to rollapp decimals")
	}

	// Calculate the dust/precision loss (amount that will be lost due to decimal conversion)
	dust := msg.Token.Amount.Sub(convertedAmtInRollappDecimals)

	// burn the original tokens from the sender
	delta := sdk.NewCoin(msg.Token.Denom, msg.Token.Amount.Sub(convertedAmt))
	if err := w.BurnCoins(ctx, sender, delta); err != nil {
		return nil, errorsmod.Wrapf(err, "burn original tokens from sender")
	}

	// Log the balance after burning delta
	balanceAfterDeltaBurn := w.bankKeeper.GetBalance(ctx, sender, msg.Token.Denom)
	ctx.Logger().Info("Account balance after delta burn",
		"sender", msg.Sender,
		"denom", msg.Token.Denom,
		"balance", balanceAfterDeltaBurn.Amount.String(),
		"delta_burned", delta.Amount.String(),
	)

	// Log the conversion details for debugging
	ctx.Logger().Info("Token conversion on transfer",
		"sender", msg.Sender,
		"original_amount", msg.Token.Amount.String(),
		"converted_bridge_amount", convertedAmt.String(),
		"converted_rollapp_amount", convertedAmtInRollappDecimals.String(),
		"dust_to_burn", dust.String(),
		"delta_to_burn", delta.String(),
		"denom", msg.Token.Denom,
	)

	// Burn the dust from the sender (precision loss due to decimal conversion)
	if !dust.IsZero() {
		ctx.Logger().Error("dust", "dust", dust.String())
		// dustCoin := sdk.NewCoin(msg.Token.Denom, dust)
		// if err := w.BurnCoins(ctx, sender, dustCoin); err != nil {
		// 	return nil, errorsmod.Wrapf(err, "burn dust tokens from sender")
		// }
	}

	// Create a new message with the converted token (in bridge decimals)
	// The IBC packet will contain this amount, which will be converted back on the receiving chain
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
	resp, err := w.Keeper.Transfer(goCtx, convertedMsg)

	// Log the balance after the transfer completes
	if err == nil {
		balanceAfterTransfer := w.bankKeeper.GetBalance(ctx, sender, msg.Token.Denom)
		ctx.Logger().Info("Account balance after transfer",
			"sender", msg.Sender,
			"denom", msg.Token.Denom,
			"balance", balanceAfterTransfer.Amount.String(),
		)
	}

	return resp, err
}
