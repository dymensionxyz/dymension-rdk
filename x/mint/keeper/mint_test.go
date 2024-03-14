package keeper_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	utils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/mint/types"
)

func TestMinting(t *testing.T) {
	/* ---------------------------------- setup --------------------------------- */
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestMintKeeperFromApp(t, app)

	params := k.GetParams(ctx)
	params.MintEpochSpreadFactor = 100
	minter := types.Minter{
		CurrentInflationRate: sdk.NewDecWithPrec(15, 2), // 15%
	}
	k.SetParams(ctx, params)
	k.SetMinter(ctx, minter)

	// set expectations
	totalSupplyAmt := sdk.NewInt(100000000) // 100M
	totalSupplyCoin := sdk.NewCoin(params.MintDenom, totalSupplyAmt)
	expectedMintedAmt := sdk.NewInt(150000) // 150K (15% of 100M / spread_factor)

	/* ---------------------------------- test ---------------------------------- */
	//assert initial state
	recipientAcc := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	initialBalance := app.BankKeeper.GetBalance(ctx, recipientAcc, params.MintDenom)
	require.True(t, initialBalance.IsZero())

	//mint supply
	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(totalSupplyCoin))
	require.NoError(t, err)

	// mint coins, update supply
	mintedCoins, err := k.HandleMintingEpoch(ctx)
	require.NoError(t, err)
	require.False(t, mintedCoins.IsZero())

	mintedCoin := mintedCoins[0]

	// assert minted coins
	require.Equal(t, expectedMintedAmt, mintedCoin.Amount)

	// assert new supply
	distrBalance := app.BankKeeper.GetBalance(ctx, recipientAcc, params.MintDenom)
	require.True(t, mintedCoins.IsEqual(sdk.NewCoins(distrBalance)))

	newSupply := app.BankKeeper.GetSupply(ctx, params.MintDenom)
	assert.True(t, newSupply.IsEqual(totalSupplyCoin.Add(mintedCoin)))
}

//TODO: test start time

func TestCalcMintedCoins(t *testing.T) {
	var DymDecimals = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

	testCases := []struct {
		name                 string
		totalSupply          sdk.Int
		currentInflationRate sdk.Dec
		spreadFactor         int64
		expectedAmount       sdk.Int
	}{
		{
			name:                 "Test Default Params",
			totalSupply:          sdk.NewInt(1000000),
			currentInflationRate: sdk.NewDecWithPrec(15, 2), // 15%
			spreadFactor:         100,
			expectedAmount:       sdk.NewInt(1500),
		},
		{
			name:                 "Test dymension decimals",
			totalSupply:          sdk.NewInt(1000000).Mul(DymDecimals),
			currentInflationRate: sdk.NewDecWithPrec(15, 2), // 15%
			spreadFactor:         100,
			expectedAmount:       sdk.NewInt(1500).Mul(DymDecimals),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := utils.Setup(t, false)
			k, ctx := testkeepers.NewTestMintKeeperFromApp(t, app)

			// Set minter
			minter := types.Minter{
				CurrentInflationRate: tc.currentInflationRate,
			}
			k.SetMinter(ctx, minter)

			params := k.GetParams(ctx)
			params.MintEpochSpreadFactor = tc.spreadFactor
			k.SetParams(ctx, params)

			mintedCoins := k.CalcMintedCoins(ctx, tc.totalSupply)
			require.False(t, mintedCoins.IsZero())
			assert.Equal(t, tc.expectedAmount, mintedCoins.TruncateInt())
		})

	}
}
