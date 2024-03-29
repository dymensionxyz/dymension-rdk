package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/dymensionxyz/dymension-rdk/testutil/app"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	utils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/mint/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	defaultEpochIdentifier                            = "minute"
	defaultMintingRewardsDistributionStartEpoch int64 = 1
	defaultFeeCollectorName                           = "fee_collector"
	defaultInflationRate                              = "1000.0"
	defaultBalanceAmt                                 = int64(1000000000)
)

var (
	defaultReductionFactor = sdk.NewDec(2).Quo(sdk.NewDec(3))
)

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
	mintKeeper.GetMinter(ctx).CurrentInflationRate.AddMut(sdk.MustNewDecFromStr(defaultInflationRate))
	// fund the fee collector account
	utils.FundModuleAccount(app, ctx, app.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName).GetName(), sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(defaultBalanceAmt))))

	// Get fee collector balance
	feeCollectorBalance := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(defaultFeeCollectorName), sdk.DefaultBondDenom)

	// For this hook to be called, identifier must be minute and epoch number must be > than 1
	testCases := []struct {
		name               string
		epoch              int64
		expectDistribution bool
	}{
		{
			name:               "before start epoch - no distributions",
			epoch:              defaultMintingRewardsDistributionStartEpoch - 1,
			expectDistribution: false,
		},
		{
			name:               "at start epoch - distributes",
			epoch:              defaultMintingRewardsDistributionStartEpoch,
			expectDistribution: true,
		},
		{
			name:               "after start epoch - distributes",
			epoch:              defaultMintingRewardsDistributionStartEpoch + 1,
			expectDistribution: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Get the current epoch
			epoch, found := epochKeeper.GetEpochInfo(ctx, defaultEpochIdentifier)
			suite.Require().True(found)
			epoch.CurrentEpoch = tc.epoch

			// Mint coin and distribute
			mintHook.AfterEpochEnd(ctx, epoch)

			// Check that the hook was called correctly
			newFeeCollectorBalance := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(defaultFeeCollectorName), sdk.DefaultBondDenom)
			if tc.expectDistribution {
				require.True(suite.T(), newFeeCollectorBalance.Amount.GT(feeCollectorBalance.Amount))
			} else {
				require.Equal(suite.T(), feeCollectorBalance.Amount, newFeeCollectorBalance.Amount)
			}
		})
	}
}
