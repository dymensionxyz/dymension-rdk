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
	// Most values here are taken from mainnet genesis to mimic real-world behavior:
	// https://github.com/osmosis-labs/networks/raw/main/osmosis-1/genesis.json
	defaultGenesisEpochProvisions = "821917808219.178082191780821917"
	defaultEpochIdentifier        = "minute"
	// actual value taken from mainnet for sanity checking calculations.
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

// TODO: Add tests for checking the minting amount, other cases when coin should not be minted, etc
func TestAfterDistributeMintedCoin(t *testing.T) {
	// Setup your test context and keeper
	app := utils.Setup(t, false)
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
	// Get the current epoch
	epoch, found := epochKeeper.GetEpochInfo(ctx, defaultEpochIdentifier)
	require.True(t, found)
	epoch.CurrentEpoch = defaultMintingRewardsDistributionStartEpoch + 1

	// Mint coin and distribute
	mintHook.AfterEpochEnd(ctx, epoch)

	// Check that the hook was called (fee collector balance should be updated)
	newFeeCollectorBalance := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(defaultFeeCollectorName), sdk.DefaultBondDenom)
	require.True(t, newFeeCollectorBalance.Amount.GT(feeCollectorBalance.Amount))
}
