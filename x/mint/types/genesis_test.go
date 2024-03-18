package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/mint/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState(t *testing.T) {
	testCases := []struct {
		name        string
		params      func() types.Params
		minter      func() types.Minter
		expectedErr bool
	}{
		{
			name:   "default",
			params: types.DefaultParams,
			minter: types.InitialMinter,
		},
		{
			name: "missing epoch identifier",
			params: func() types.Params {
				p := types.DefaultParams()
				p.InflationChangeEpochIdentifier = ""
				return p
			},
			minter:      types.InitialMinter,
			expectedErr: true,
		},
		{
			name: "bad inflation rate",
			params: func() types.Params {
				p := types.DefaultParams()
				p.TargetInflationRate = sdk.MustNewDecFromStr("2.0")
				return p
			},
			minter:      types.InitialMinter,
			expectedErr: true,
		},
		{
			name:   "bad minter inflation rate",
			params: types.DefaultParams,
			minter: func() types.Minter {
				return types.Minter{
					CurrentInflationRate: sdk.MustNewDecFromStr("2.0"),
				}
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		state := types.NewGenesisState(tc.minter(), tc.params())
		err := types.ValidateGenesis(*state)
		if tc.expectedErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}
