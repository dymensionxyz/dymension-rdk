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

func TestState_Validate(t *testing.T) {
	tests := []struct {
		name    string
		state   State
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid state",
			state: State{
				HubPortAndChannel: &PortAndChannel{
					Port:    "transfer",
					Channel: "channel-0",
				},
				GenesisAccounts:          []GenesisAccount{{Address: sdk.MustBech32ifyAddressBytes("dym", sampleAddr()), Amount: sdk.NewInt(100)}},
				OutboundTransfersEnabled: true,
			},
			wantErr: false,
		},
		{
			name: "invalid amount",
			state: State{
				GenesisAccounts: []GenesisAccount{
					{Address: sdk.MustBech32ifyAddressBytes("dym", sampleAddr()), Amount: sdk.NewInt(0)},
				},
			},
			wantErr: true,
			errMsg:  "invalid amount",
		},
		{
			name: "invalid address",
			state: State{
				GenesisAccounts: []GenesisAccount{
					{Address: sdk.MustBech32ifyAddressBytes("hub", sampleAddr()), Amount: sdk.NewInt(100)},
				},
			},
			wantErr: true,
			errMsg:  "invalid address",
		},
		{
			name: "multiple valid accounts",
			state: State{
				GenesisAccounts: []GenesisAccount{
					{Address: sdk.MustBech32ifyAddressBytes("dym", sampleAddr()), Amount: sdk.NewInt(100)},
					{Address: sdk.MustBech32ifyAddressBytes("dym", sampleAddr()), Amount: sdk.NewInt(200)},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty accounts",
			state:   State{},
			wantErr: false,
		},
		{
			name: "invalid hub port and channel",
			state: State{
				HubPortAndChannel: &PortAndChannel{Port: "invalid/port", Channel: "invalid/channel"},
			},
			wantErr: true,
			errMsg:  "invalid port Id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.state.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
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
				State: State{
					GenesisAccounts: []GenesisAccount{
						{Address: sdk.MustBech32ifyAddressBytes("dym", sampleAddr()), Amount: sdk.NewInt(100)},
					},
					OutboundTransfersEnabled: false,
					HubPortAndChannel:        nil,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid state",
			genesisState: GenesisState{
				State: State{
					GenesisAccounts: []GenesisAccount{
						{Address: "invalid_address", Amount: sdk.NewInt(100)},
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid address",
		},
		{
			name: "outbound transfers enabled",
			genesisState: GenesisState{
				State: State{
					OutboundTransfersEnabled: true,
				},
			},
			wantErr: true,
			errMsg:  "outbound transfers should be disabled in genesis",
		},
		{
			name: "hub port and channel set",
			genesisState: GenesisState{
				State: State{
					HubPortAndChannel: &PortAndChannel{
						Port:    "transfer",
						Channel: "channel-0",
					},
				},
			},
			wantErr: true,
			errMsg:  "hub port and channel should not be set in genesis",
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
