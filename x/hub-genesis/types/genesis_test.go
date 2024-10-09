package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleAddr() []byte {
	return secp256k1.GenPrivKey().PubKey().Address().Bytes()
}

func TestGenesisState_ValidateBasic(t *testing.T) {
	tests := []struct {
		name         string
		genesisState GenesisState
		wantErr      bool
		errMsg       string
	}{
		{
			name: "valid genesis state",
			genesisState: GenesisState{
				GenesisAccounts: []GenesisAccount{
					{Address: sdk.MustBech32ifyAddressBytes("dym", sampleAddr()), Amount: sdk.NewInt(100)},
					{Address: sdk.MustBech32ifyAddressBytes("dym", sampleAddr()), Amount: sdk.NewInt(100)},
				},
			},
			wantErr: false,
		},
		{
			name:         "valid genesis - no genesis accounts",
			genesisState: GenesisState{},
			wantErr:      false,
		},
		{
			name: "invalid state - invalid address",
			genesisState: GenesisState{
				GenesisAccounts: []GenesisAccount{
					{Address: sdk.MustBech32ifyAddressBytes("dym", sampleAddr()), Amount: sdk.NewInt(100)},
					{Address: "invalid_address", Amount: sdk.NewInt(100)},
				},
			},
			wantErr: true,
			errMsg:  "invalid address",
		},
		{
			name: "invalid state - invalid amount",
			genesisState: GenesisState{
				GenesisAccounts: []GenesisAccount{
					{Address: sdk.MustBech32ifyAddressBytes("dym", sampleAddr()), Amount: sdk.NewInt(100)},
					{Address: sdk.MustBech32ifyAddressBytes("dym", sampleAddr()), Amount: sdk.NewInt(-100)},
				},
			},
			wantErr: true,
			errMsg:  "invalid amount",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.genesisState.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
