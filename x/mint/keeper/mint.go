package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/mint/types"
)

func (k Keeper) HandleMintingEpoch(ctx sdk.Context) (sdk.Coins, error) {
	var mintedCoins sdk.Coins
	params := k.GetParams(ctx)

	// fetch stored minter & params
	minter := k.GetMinter(ctx)

	//calculate coins
	total := k.bankKeeper.GetSupply(ctx, params.MintDenom)
	mintAmount := minter.CurrentInflationRate.MulInt(total.Amount).QuoInt(sdk.NewInt(params.MintEpochSpreadFactor))

	// mint coins, update supply
	mintedCoins = sdk.NewCoins(sdk.NewCoin(params.MintDenom, mintAmount.TruncateInt()))
	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		return mintedCoins, err
	}

	// send the minted coins to the fee collector account
	err = k.DistributeMintedCoin(ctx, mintedCoins)
	if err != nil {
		return mintedCoins, err
	}

	return mintedCoins, nil
}

// ___________________________________________________________________________________________________

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// DistributeMintedCoins implements distribution of minted coins from mint to external modules.
func (k Keeper) DistributeMintedCoin(ctx sdk.Context, mintedCoins sdk.Coins) error {
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, mintedCoins)
	if err != nil {
		return err
	}

	// call a hook after the minting and distribution of new coins
	k.hooks.AfterDistributeMintedCoin(ctx, mintedCoins)

	return nil
}
