package types_test

import (
	"testing"

	"github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/stretchr/testify/require"
)

func TestGenesisState(t *testing.T) {

	testParams := types.NewParams("mock", 1, "5f8393904fb1e9c616fe89f013cafe7501a63f86")
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
			name: "missing commit",
			params: func() types.Params {
				p := testParams
				p.Commit = ""
				return p
			},
			expectedErr: true,
		},
		{
			name: "wrong length commit",
			params: func() types.Params {
				p := testParams
				p.Commit = "fdasfewkq102382w523"
				return p
			},
			expectedErr: true,
		},
		{
			name: "version not alphanumeric",
			params: func() types.Params {
				p := testParams
				p.Commit = "3a19edd887_9b576a866750bc9d480ada53d2c0d"
				return p
			},
			expectedErr: true,
		},
		{
			name: "wrong drs version",
			params: func() types.Params {
				p := testParams
				p.Version = 0
				return p
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		state := types.NewGenesisState(tc.params())
		err := types.ValidateGenesis(state)
		if tc.expectedErr {
			require.ErrorIs(t, err, gerrc.ErrInvalidArgument)
		} else {
			require.NoError(t, err)
		}
	}
}
