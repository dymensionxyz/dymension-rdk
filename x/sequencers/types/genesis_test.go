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

	for _, tc := range []struct {
		desc     string
		genState types.GenesisState
		valid    bool
	}{
		{
			desc: "valid",
			genState: types.GenesisState{
				Params:                 types.DefaultParams(),
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
				GenesisOperatorAddress: sdk.AccAddress(pk.Address()).String(),
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
