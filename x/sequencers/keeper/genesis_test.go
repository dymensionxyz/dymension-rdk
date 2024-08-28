package keeper_test

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"
)

func TestGetCoinbaseAddressAlt(t *testing.T) {
	_ = common.BytesToAddress(nil)
}

// A regression test to make sure we are compatible with ethermint, which requires a 'coinbase'
// addr for compat with EVM op codes.
// https://github.com/dymensionxyz/ethermint/blob/b1506ae83050d2361857251766d93253e317900c/x/evm/keeper/state_transition.go#L41-L44
func TestGetCoinbaseAddress(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	val, ok := k.GetValidatorByConsAddr(ctx, ctx.BlockHeader().ProposerAddress)
	require.True(t, ok)
	_ = common.BytesToAddress(val.GetOperator())
}

func TestInitAndExportGenesis(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	exp := types.GenesisState{
		Params: types.DefaultParams(),
		Sequencers: []types.Sequencer{
			{
				Validator:  &utils.Proposer,
				RewardAddr: utils.OperatorAcc().String(),
			},
		},
	}
	k.InitGenesis(ctx, exp)
	got := k.ExportGenesis(ctx)
	require.Equal(t, &exp, got)
}
