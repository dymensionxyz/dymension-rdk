package types_test

import (
	"testing"

	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	g := types.GenesisState{
		Params: types.DefaultParams(),
		Sequencers: []types.Sequencer{
			{
				Validator:  &utils.Proposer,
				RewardAddr: "",
			},
		},
	}

	t.Run("ok", func(t *testing.T) {
		c := g.MustClone()
		require.NoError(t, c.ValidateGenesis())
	})
	t.Run("nil seq", func(t *testing.T) {
		c := g.MustClone()
		c.Sequencers[0].Validator = nil
		require.Error(t, c.ValidateGenesis())
	})
	t.Run("bad operator", func(t *testing.T) {
		c := g.MustClone()
		c.Sequencers[0].Validator.OperatorAddress = "foo"
		require.Error(t, c.ValidateGenesis())
	})
	t.Run("bad cons", func(t *testing.T) {
		c := g.MustClone()
		c.Sequencers[0].Validator.ConsensusPubkey = nil
		require.Error(t, c.ValidateGenesis())
	})
	t.Run("bad reward addr", func(t *testing.T) {
		c := g.MustClone()
		c.Sequencers[0].RewardAddr = "foo"
		require.Error(t, c.ValidateGenesis())
	})
}
