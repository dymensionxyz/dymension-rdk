package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	utils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
)

func TestMintDistribution(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestMintKeeperFromApp(t, app)

	// // fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	//assert initial state
	recipientAcc := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	initialBalance := app.BankKeeper.GetBalance(ctx, recipientAcc, params.MintDenom)
	require.True(t, initialBalance.IsZero())

	initialSupply := app.BankKeeper.GetSupply(ctx, params.MintDenom)

	// mint coins, update supply
	mintedCoin, err := k.HandleMintingEpoch(ctx)
	require.NoError(t, err)

	// TODO: assert amounts minted
	_ = minter
	_ = initialSupply

	distrBalance := app.BankKeeper.GetBalance(ctx, recipientAcc, params.MintDenom)
	require.Equal(t, mintedCoin, sdk.NewCoins(distrBalance))
}
