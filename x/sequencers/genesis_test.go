package sequencers_test

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/nullify"
	utils "github.com/dymensionxyz/dymension-rdk/testutil/utils"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/testutils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFailedInitGenesis(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	pks := utils.CreateTestPubKeys(1)
	addr := sdk.ValAddress(pks[0].Address())
	val := utils.NewValidator(t, addr, pks[0])

	genesisState := types.GenesisState{
		Params:     types.DefaultParams(),
		Sequencers: []stakingtypes.Validator{val},
		Exported:   false,
	}

	//mess with the pubkey value
	pkAny, err := codectypes.NewAnyWithValue(&types.Params{})
	assert.NoError(t, err)
	val.ConsensusPubkey = pkAny
	genesisState.Sequencers = append(genesisState.Sequencers, val)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	sequencers.InitGenesis(ctx, *k, genesisState)
}

func TestGenesis(t *testing.T) {
	k, ctx := testutils.NewTestSequencerKeeper(t)

	pks := utils.CreateTestPubKeys(10)
	pks = pks[5:]
	addr1 := sdk.ValAddress(pks[0].Address())
	addr2 := sdk.ValAddress(pks[1].Address())

	genesisState := types.GenesisState{
		Params:     types.DefaultParams(),
		Sequencers: []stakingtypes.Validator{},
		Exported:   false,
	}

	//Init dymint sequencers
	err := k.SetDymintSequencerByAddr(ctx, sdk.GetConsAddress(pks[0]), 0)
	require.NoError(t, err)
	err = k.SetDymintSequencerByAddr(ctx, sdk.GetConsAddress(pks[1]), 0)
	require.NoError(t, err)

	//Init rollapp sequencers
	genesisState.Sequencers = append(genesisState.Sequencers, utils.NewValidator(t, addr1, pks[0]))
	genesisState.Sequencers = append(genesisState.Sequencers, utils.NewValidator(t, addr2, pks[1]))

	sequencers.InitGenesis(ctx, *k, genesisState)
	got := sequencers.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.Sequencers, got.Sequencers)
	require.EqualValues(t, genesisState.Params, got.Params)
}
