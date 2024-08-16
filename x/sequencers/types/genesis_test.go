package types_test

import (
	"testing"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

// TODO: check

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState types.GenesisState
		valid    bool
	}{
		{
			desc: "valid",
			genState: types.GenesisState{
				Params: types.DefaultParams(),
			},
			valid: true,
		},
		{
			desc:     "default - missing operator address",
			genState: *types.DefaultGenesis(),
			valid:    false,
		},
		{
			desc: "not a val address",
			genState: types.GenesisState{
				Params: types.DefaultParams(),
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
