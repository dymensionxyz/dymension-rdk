package types_test

import (
	"testing"

	"github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
	"github.com/stretchr/testify/require"
)

func TestGenesisState(t *testing.T) {
	testCases := []struct {
		name        string
		params      func() types.Params
		expectedErr bool
	}{
		{
			name:   "default",
			params: types.DefaultParams,
		},
		{
			name: "missing version",
			params: func() types.Params {
				p := types.DefaultParams()
				p.Version = ""
				return p
			},
			expectedErr: true,
		},
		{
			name: "wrong length version",
			params: func() types.Params {
				p := types.DefaultParams()
				p.Version = "fdasfewkq102382w523"
				return p
			},
			expectedErr: true,
		},
		{
			name: "version not alphanumeric",
			params: func() types.Params {
				p := types.DefaultParams()
				p.Version = "3a19edd887a9b576a866750bc9d480ada53d2c0d"
				return p
			},
			expectedErr: true,
		},
		{
			name: "block max gas too small",
			params: func() types.Params {
				p := types.DefaultParams()
				p.Blockmaxgas = 0
				return p
			},
			expectedErr: true,
		},
		{
			name: "block max size too small",
			params: func() types.Params {
				p := types.DefaultParams()
				p.Blockmaxsize = 50000
				return p
			},
			expectedErr: true,
		},
		{
			name: "block max size too big",
			params: func() types.Params {
				p := types.DefaultParams()
				p.Blockmaxsize = 1000000
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
