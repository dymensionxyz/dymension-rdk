package types_test

import (
	"testing"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: *types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: types.GenesisState{
				Params: types.DefaultParams(),
				Sequencers: []stakingtypes.Validator{{
					OperatorAddress: "sequencer1",
				}, {
					OperatorAddress: "sequencer2",
				}},
				Exported: false,
			},
			valid: true,
		},
		{
			desc: "duplicated sequencer",
			genState: types.GenesisState{
				Params: types.DefaultParams(),
				Sequencers: []stakingtypes.Validator{{
					OperatorAddress: "sequencer1",
				}, {
					OperatorAddress: "sequencer1",
				}},
				Exported: false,
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.ValidateGenesis()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
