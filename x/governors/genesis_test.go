package governors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	governors "github.com/dymensionxyz/dymension-rdk/x/governors"
	"github.com/dymensionxyz/dymension-rdk/x/governors/teststaking"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

func TestValidateGenesis(t *testing.T) {
	genGovernors1 := make([]types.Governor, 1, 5)
	pk := ed25519.GenPrivKey().PubKey()
	genGovernors1[0] = teststaking.NewGovernor(t, sdk.ValAddress(pk.Address()))
	genGovernors1[0].Tokens = sdk.OneInt()
	genGovernors1[0].DelegatorShares = sdk.OneDec()

	tests := []struct {
		name    string
		mutate  func(*types.GenesisState)
		wantErr bool
	}{
		{"default", func(*types.GenesisState) {}, false},
		// validate genesis governors
		{"duplicate governor", func(data *types.GenesisState) {
			data.Governors = genGovernors1
			data.Governors = append(data.Governors, genGovernors1[0])
		}, true},
		{"no delegator shares", func(data *types.GenesisState) {
			data.Governors = genGovernors1
			data.Governors[0].DelegatorShares = sdk.ZeroDec()
		}, true},
		{"jailed and bonded governor", func(data *types.GenesisState) {
			data.Governors = genGovernors1
			data.Governors[0].Status = types.Bonded
		}, true},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			genesisState := types.DefaultGenesisState()
			tt.mutate(genesisState)

			if tt.wantErr {
				assert.Error(t, governors.ValidateGenesis(genesisState))
			} else {
				assert.NoError(t, governors.ValidateGenesis(genesisState))
			}
		})
	}
}
