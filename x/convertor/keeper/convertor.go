package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/dymensionxyz/dymension-rdk/x/convertor/types"
)

// ConversionRequired checks if a conversion is required for a given denom.
// The denom parameter should be in the IBC hash format (ibc/XXX).
func (k Keeper) ConversionRequired(ctx sdk.Context, denom string) (bool, error) {
	pair, err := k.hubKeeper.GetDecimalConversionPair(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return pair.FromToken == denom, nil
}

// ConvertFromBridgeAmt converts an amount from bridge token decimals to rollapp decimals (18)
// and emits an event with the conversion details
func (k Keeper) ConvertFromBridgeAmt(
	ctx sdk.Context,
	amount math.Int,
) (math.Int, error) {
	pair, err := k.hubKeeper.GetDecimalConversionPair(ctx)
	if err != nil {
		return math.Int{}, err
	}

	convertedAmt, err := types.ConvertAmount(amount, pair.FromDecimals, 18)
	if err != nil {
		return math.Int{}, err
	}

	// Log conversion details
	k.Logger(ctx).Debug("ConvertFromBridgeAmt",
		"from_decimals", pair.FromDecimals,
		"to_decimals", 18,
		"input_amount", amount.String(),
		"output_amount", convertedAmt.String(),
	)

	// Emit conversion event
	k.emitConversionEvent(ctx, amount, convertedAmt)

	return convertedAmt, nil
}

// ConvertToBridgeAmt converts an amount from rollapp decimals (18) to bridge token decimals
// and emits an event with the conversion details
func (k Keeper) ConvertToBridgeAmt(
	ctx sdk.Context,
	amount math.Int,
) (math.Int, error) {
	pair, err := k.hubKeeper.GetDecimalConversionPair(ctx)
	if err != nil {
		return math.Int{}, err
	}

	convertedAmt, err := types.ConvertAmount(amount, 18, pair.FromDecimals)
	if err != nil {
		return math.Int{}, err
	}

	// Log conversion details
	k.Logger(ctx).Debug("ConvertToBridgeAmt",
		"from_decimals", 18,
		"to_decimals", pair.FromDecimals,
		"input_amount", amount.String(),
		"output_amount", convertedAmt.String(),
	)

	// Emit conversion event
	k.emitConversionEvent(ctx, amount, convertedAmt)

	return convertedAmt, nil
}

// emitConversionEvent emits a decimal conversion event
func (k Keeper) emitConversionEvent(
	ctx sdk.Context,
	originalAmt math.Int,
	convertedAmt math.Int,
) {
	attrs := []sdk.Attribute{
		sdk.NewAttribute(types.AttributeKeyOriginalAmount, originalAmt.String()),
		sdk.NewAttribute(types.AttributeKeyConvertedAmount, convertedAmt.String()),
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeDecimalConversion, attrs...))
}

// BurnCoins burns coins from an account
func (k Keeper) BurnCoins(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error {
	// Send coins from account to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, ibctransfertypes.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}
	// Burn coins from module
	return k.bankKeeper.BurnCoins(ctx, ibctransfertypes.ModuleName, sdk.NewCoins(coin))
}

// MintCoins mints coins to an account
func (k Keeper) MintCoins(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error {
	// Mint coins to module
	if err := k.bankKeeper.MintCoins(ctx, ibctransfertypes.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}
	// Send coins from module to account
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, ibctransfertypes.ModuleName, addr, sdk.NewCoins(coin))
}
