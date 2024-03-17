package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	utils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/mint/types"
)

func TestMintDistribution(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestMintKeeperFromApp(t, app)

	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := types.DefaultParams()
	const mintAmt = 1000000000000000000
	params.GenesisEpochProvisions = sdk.NewDec(mintAmt)
	params.MintDenom = "mintDenom"
	k.SetParams(ctx, params)

	//assert initial state
	recipientAcc := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	initialBalance := app.BankKeeper.GetBalance(ctx, recipientAcc, params.MintDenom)
	require.True(t, initialBalance.IsZero())

	// mint coins, update supply
	mintedCoin := minter.EpochProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)
	err := k.MintCoins(ctx, mintedCoins)
	require.NoError(t, err)

	// send the minted coins to their respective module accounts (e.g. staking rewards to the feecollector)
	err = k.DistributeMintedCoin(ctx, mintedCoin)
	require.NoError(t, err)

	distrBalance := app.BankKeeper.GetBalance(ctx, recipientAcc, params.MintDenom)

	require.Equal(t, mintedCoin, distrBalance)
	require.Equal(t, mintedCoin.Denom, params.MintDenom)
}
