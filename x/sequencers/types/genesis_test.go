package types_test

import (
	"testing"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState types.GenesisState
		valid    bool
	}{
		{
			desc:     "default",
			genState: *types.DefaultGenesis(),
			valid:    true,
		},
		//TODO: bad params
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
