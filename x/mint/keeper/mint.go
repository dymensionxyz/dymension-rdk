package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/mint/types"
)

func (k Keeper) HandleMintingEpoch(ctx sdk.Context) (sdk.Coins, error) {
	var mintedCoins sdk.Coins
	params := k.GetParams(ctx)

	// calculate coins
	total := k.bankKeeper.GetSupply(ctx, params.MintDenom)
	mintAmount := k.CalcMintedCoins(ctx, total.Amount)
	if mintAmount.IsZero() {
		return mintedCoins, nil
	}

	// mint coins, update supply
	mintedCoins = sdk.NewCoins(sdk.NewCoin(params.MintDenom, mintAmount.TruncateInt()))
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintedCoins)
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

func (k Keeper) CalcMintedCoins(ctx sdk.Context, totalAmt math.Int) sdk.Dec {
	params := k.GetParams(ctx)
	minter := k.GetMinter(ctx)

	epoch, _ := k.epochKeeper.GetEpochInfo(ctx, params.MintEpochIdentifier)
	spreadFactor := types.InflationAnnualDuration / epoch.Duration
	mintAmount := minter.CurrentInflationRate.MulInt(totalAmt.QuoRaw(int64(spreadFactor)))
	return mintAmount
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
