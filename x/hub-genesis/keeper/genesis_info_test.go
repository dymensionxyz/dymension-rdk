package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis_HappyFlow(t *testing.T) {
	genesisBridgeFunds := sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100_000))
	genAccounts := []types.GenesisAccount{
		{
			Address: utils.AccAddress().String(),
			Amount:  genesisBridgeFunds.Amount.QuoRaw(2),
		},
		{
			Address: utils.AccAddress().String(),
			Amount:  genesisBridgeFunds.Amount.QuoRaw(2),
		},
	}

	app := utils.SetupWithGenesisBridge(t, genesisBridgeFunds, genAccounts)
	k, ctx := testkeepers.NewTestHubGenesisKeeperFromApp(app)

	// change some state values post the genesis, make sure it doesn't affect the genesis info
	utils.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000))
	require.NotEqual(t, app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom), genesisBridgeFunds)

	gInfo := k.GetGenesisInfo(ctx)
	assert.Equal(t, genAccounts, gInfo.GenesisAccounts)
	// assert native denom
	assert.Equal(t, genesisBridgeFunds.Denom, gInfo.NativeDenom.Base)
	// assert initial supply
	assert.Equal(t, genesisBridgeFunds.Amount, gInfo.InitialSupply)
}

func TestInitGenesis_MissingGenesisFundsOnGenesis(t *testing.T) {
	genesisBridgeFunds := sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100_000))
	genAccounts := []types.GenesisAccount{
		{
			Address: utils.AccAddress().String(),
			Amount:  genesisBridgeFunds.Amount.MulRaw(2), // genesis account has more funds than the module account
		},
	}
	assert.Panics(t, func() {
		utils.SetupWithGenesisBridge(t, genesisBridgeFunds, genAccounts)
	})
}

func TestInitGenesis_MissingDenomMetadata(t *testing.T) {
	genesisBridgeFunds := sdk.NewCoin("newdenom", math.NewInt(100_000))
	genAccounts := []types.GenesisAccount{
		{
			Address: utils.AccAddress().String(),
			Amount:  genesisBridgeFunds.Amount.QuoRaw(2),
		},
		{
			Address: utils.AccAddress().String(),
			Amount:  genesisBridgeFunds.Amount.QuoRaw(2),
		},
	}
	assert.Panics(t, func() {
		utils.SetupWithGenesisBridge(t, genesisBridgeFunds, genAccounts)
	})
}

func TestInitGenesis_NoNativeDenom(t *testing.T) {
	app := utils.SetupWithNoNativeDenom(t)
	k, ctx := testkeepers.NewTestHubGenesisKeeperFromApp(app)

	gInfo := k.GetGenesisInfo(ctx)
	// assert native denom
	assert.Equal(t, "", gInfo.NativeDenom.Base)
	// assert initial supply
	assert.Equal(t, math.ZeroInt(), gInfo.InitialSupply)
}
