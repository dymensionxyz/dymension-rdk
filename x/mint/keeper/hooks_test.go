package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	utils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/mint/keeper"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/dymensionxyz/dymension-rdk/x/mint/types"
)

const (
	defaultEpochIdentifier                            = "hour"
	defaultMintingRewardsDistributionStartEpoch int64 = 1
	defaultFeeCollectorName                           = "fee_collector"
	defaultInflationRate                              = "0.15" // 15%
	defaultBalanceAmt                                 = int64(1000000000)
)

var defaultReductionFactor = sdk.NewDec(2).Quo(sdk.NewDec(3))

type MintKeeperTestSuite struct {
	suite.Suite

	app *app.App
	k   keeper.Keeper
	ctx sdk.Context
}

func TestHooksTestSuite(t *testing.T) {
	suite.Run(t, new(MintKeeperTestSuite))
}

func (suite *MintKeeperTestSuite) TestAfterDistributeMintedCoin() {
	// Setup your test context and keeper
	app := utils.Setup(suite.T(), false)
	mintKeeper, _ := testkeepers.NewTestMintKeeperFromApp(app)
	epochKeeper, ctx := testkeepers.NewTestEpochKeeperFromApp(app)

	// Get mint hook
	mintHook := mintKeeper.Hooks()
	// Set InflationRate for coin minting
	minter := minttypes.Minter{
		CurrentInflationRate: sdk.NewDecWithPrec(15, 2), // 15%
	}
	mintKeeper.SetMinter(ctx, minter)
	// fund the fee collector account
	utils.FundModuleAccount(app, ctx, app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName).GetName(), sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(defaultBalanceAmt))))

	// Get fee collector balance
	feeCollectorBalance := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(defaultFeeCollectorName), sdk.DefaultBondDenom)

	// For this hook to be called, identifier must be minute and epoch number must be > than 1
	testCases := []struct {
		name               string
		epoch              int64
		expectDistribution bool
		expectedMintedAmt  math.Int
	}{
		{
			name:               "before start epoch - no distributions",
			epoch:              defaultMintingRewardsDistributionStartEpoch - 1,
			expectDistribution: false,
			expectedMintedAmt:  sdk.NewInt(0),
		},
		{
			name:               "at start epoch - distributes",
			epoch:              defaultMintingRewardsDistributionStartEpoch,
			expectDistribution: true,
			expectedMintedAmt:  sdk.NewInt(17123), // 17123 (15% of 1000M / (365*24)) or (inflationRate * totalSupply / spreadFactor)
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Get the current epoch
			epoch, found := epochKeeper.GetEpochInfo(ctx, defaultEpochIdentifier)
			require.True(suite.T(), found)
			epoch.CurrentEpoch = tc.epoch

			// Mint coin and distribute
			mintHook.AfterEpochEnd(ctx, epoch)

			// Check that the hook was called correctly
			newFeeCollectorBalance := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(defaultFeeCollectorName), sdk.DefaultBondDenom)
			if tc.expectDistribution {
				require.True(suite.T(), newFeeCollectorBalance.Amount.GT(feeCollectorBalance.Amount))
				// Check the minting amount
				actualMintedAmt := newFeeCollectorBalance.Amount.Sub(feeCollectorBalance.Amount)
				require.Equal(suite.T(), tc.expectedMintedAmt, actualMintedAmt)
			} else {
				require.Equal(suite.T(), feeCollectorBalance.Amount, newFeeCollectorBalance.Amount)
			}
		})
	}
}
