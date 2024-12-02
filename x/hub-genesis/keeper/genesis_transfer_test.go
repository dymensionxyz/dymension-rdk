package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

func TestGenesisTransferCreation(t *testing.T) {
	genesisBridgeFunds := sdk.NewCoin("stake", math.NewInt(100_000))
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

	data, err := k.PrepareGenesisBridgeData(ctx)
	require.NoError(t, err)
	packet := data.GenesisTransfer
	require.NotNil(t, packet)

	assert.Equal(t, "stake", packet.Denom)
	assert.Equal(t, genesisBridgeFunds.Amount.String(), packet.Amount)
	assert.Equal(t, app.AccountKeeper.GetModuleAddress(types.ModuleName).String(), packet.Sender)
}

func TestGenesisTransferCreation_NoGenesisAccounts(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestHubGenesisKeeperFromApp(app)

	data, err := k.PrepareGenesisBridgeData(ctx)
	require.NoError(t, err)
	require.Nil(t, data.GenesisTransfer)
}
