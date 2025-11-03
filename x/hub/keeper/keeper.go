package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/utils/collcompat"
	"github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

type Keeper struct {
	cdc                   codec.BinaryCodec
	storeKey              storetypes.StoreKey
	registeredHubDenoms   collections.KeySet[string]
	decimalConversionPair collections.Item[types.DecimalConversionPair]
	bankKeeper            types.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(collcompat.NewKVStoreService(storeKey))
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		bankKeeper: bankKeeper,
		registeredHubDenoms: collections.NewKeySet(
			sb,
			collections.NewPrefix(types.RegisteredHubDenomsKeyPrefix),
			"registered_hub_denoms",
			collections.StringKey,
		),
		decimalConversionPair: collections.NewItem(
			sb,
			collections.NewPrefix(types.DecimalConversionPairKeyPrefix),
			"decimal_conversion_pair",
			collcompat.ProtoValue[types.DecimalConversionPair](cdc),
		),
	}
}

func (k Keeper) SetHubDenom(ctx sdk.Context, denom string) error {
	return k.registeredHubDenoms.Set(ctx, denom)
}

func (k Keeper) HasHubDenom(ctx sdk.Context, denom string) (bool, error) {
	return k.registeredHubDenoms.Has(ctx, denom)
}

func (k Keeper) GetAllHubDenoms(ctx sdk.Context) (denoms []string, _ error) {
	return denoms, k.registeredHubDenoms.Walk(ctx, new(collections.Range[string]), func(d string) (_ bool, _ error) {
		denoms = append(denoms, d)
		return
	})
}

// SetDecimalConversionPair sets the decimal conversion pair
func (k Keeper) SetDecimalConversionPair(ctx sdk.Context, pair types.DecimalConversionPair) error {
	return k.decimalConversionPair.Set(ctx, pair)
}

// GetDecimalConversionPair retrieves the decimal conversion pair
func (k Keeper) GetDecimalConversionPair(ctx sdk.Context) (types.DecimalConversionPair, error) {
	pair, err := k.decimalConversionPair.Get(ctx)
	if err != nil {
		return types.DecimalConversionPair{}, err
	}
	return pair, nil
}

// ConversionRequired checks if a conversion is required for a given denom
func (k Keeper) ConversionRequired(ctx sdk.Context, denom string) (bool, error) {
	pair, err := k.GetDecimalConversionPair(ctx)
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
	pair, err := k.GetDecimalConversionPair(ctx)
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
	pair, err := k.GetDecimalConversionPair(ctx)
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
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}
	// Burn coins from module
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
}

// MintCoins mints coins to an account
func (k Keeper) MintCoins(ctx sdk.Context, addr sdk.AccAddress, coin sdk.Coin) error {
	// Mint coins to module
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return err
	}
	// Send coins from module to account
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(coin))
}
