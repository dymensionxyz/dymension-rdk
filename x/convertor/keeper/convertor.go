package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/dymensionxyz/dymension-rdk/x/convertor/types"
)

// ConversionRequired checks if a conversion is required for a given denom
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

// ConvertFromBridgeAmt converts an amount from one denom to another using the decimal conversion pair
func (k Keeper) ConvertFromBridgeAmt(ctx sdk.Context, amount math.Int) (math.Int, error) {
	pair, err := k.hubKeeper.GetDecimalConversionPair(ctx)
	if err != nil {
		return math.Int{}, err
	}

	newAmt, err := types.ConvertAmount(amount, pair.FromDecimals, 18)
	if err != nil {
		return math.Int{}, err
	}
	return newAmt, nil
}

// ConvertToBridgeAmt converts an amount to another denom using the decimal conversion pair
func (k Keeper) ConvertToBridgeAmt(ctx sdk.Context, amount math.Int) (math.Int, error) {
	pair, err := k.hubKeeper.GetDecimalConversionPair(ctx)
	if err != nil {
		return math.Int{}, err
	}

	newAmt, err := types.ConvertAmount(amount, 18, pair.FromDecimals)
	if err != nil {
		return math.Int{}, err
	}
	return newAmt, nil
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
