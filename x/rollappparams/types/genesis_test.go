package types_test

import (
	"testing"

	"github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/stretchr/testify/require"
)

func TestGenesisState(t *testing.T) {

	testParams := types.NewParams("mock", 1)
	testCases := []struct {
		name        string
		params      func() types.Params
		expectedErr bool
	}{
		{
			name: "default",
			params: func() types.Params {
				p := testParams
				return p
			},
		},
		{
			name: "wrong drs version",
			params: func() types.Params {
				p := testParams
				p.DrsVersion = 0
				return p
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		state := types.NewGenesisState(tc.params())
		err := types.ValidateGenesis(state)
		if tc.expectedErr {
			require.Error(t, err, gerrc.ErrInvalidArgument)
		} else {
			require.NoError(t, err)
		}
	}
}
