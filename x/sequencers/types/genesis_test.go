package types_test

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisState_Validate(t *testing.T) {
	pk := ed25519.GenPrivKey().PubKey()
	accAddr := sdk.AccAddress(pk.Address())

	for _, tc := range []struct {
		desc     string
		genState types.GenesisState
		valid    bool
	}{
		{
			desc: "valid",
			genState: types.GenesisState{
				Params:                 types.DefaultParams(),
				AddressPermissions:     types.DefaultAddressPermissions,
				GenesisOperatorAddress: sdk.ValAddress(pk.Address()).String(),
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
				Params:                 types.DefaultParams(),
				GenesisOperatorAddress: accAddr.String(),
			},
			valid: false,
		},
		{
			desc: "empty operator address",
			genState: types.GenesisState{
				Params:                 types.DefaultParams(),
				GenesisOperatorAddress: "",
			},
			valid: false,
		},
		{
			desc: "empty address in address permissions",
			genState: types.GenesisState{
				Params: types.DefaultParams(),
				AddressPermissions: []types.AddressPermissions{
					{
						Address: "",
						PermissionList: types.PermissionList{
							Permissions: []string{"test"},
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid permission list",
			genState: types.GenesisState{
				Params: types.DefaultParams(),
				AddressPermissions: []types.AddressPermissions{
					{
						Address:        sdk.MustBech32ifyAddressBytes(sdk.Bech32PrefixAccAddr, accAddr),
						PermissionList: types.EmptyPermissionList(),
					},
				},
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
