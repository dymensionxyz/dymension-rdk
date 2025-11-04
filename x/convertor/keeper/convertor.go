package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
	return pair.FromToken == denom || pair.ToToken == denom, nil
}

// ConvertCoin converts a coin from one denom to another using the decimal conversion pair
func (k Keeper) ConvertFromBridgeCoin(ctx sdk.Context, coin sdk.Coin) (sdk.Coin, error) {
	pair, err := k.hubKeeper.GetDecimalConversionPair(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	newAmt, err := types.ConvertAmount(coin.Amount, pair.FromDecimals, 18)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(pair.ToToken, newAmt), nil
}

// ConvertToBridgeCoin converts a coin to another denom using the decimal conversion pair
func (k Keeper) ConvertToBridgeCoin(ctx sdk.Context, coin sdk.Coin) (sdk.Coin, error) {
	pair, err := k.hubKeeper.GetDecimalConversionPair(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	newAmt, err := types.ConvertAmount(coin.Amount, 18, pair.FromDecimals)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(pair.FromToken, newAmt), nil
}

// BurnCoins burns coins from an account
func (k Keeper) BurnCoins(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error {
	// Send coins from account to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, banktypes.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}
	// Burn coins from module
	return k.bankKeeper.BurnCoins(ctx, banktypes.ModuleName, sdk.NewCoins(coin))
}

// MintCoins mints coins to an account
func (k Keeper) MintCoins(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error {
	// Mint coins to module
	if err := k.bankKeeper.MintCoins(ctx, banktypes.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}
	// Send coins from module to account
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, banktypes.ModuleName, addr, sdk.NewCoins(coin))
}
