package types

import (
	"testing"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/require"
)

func TestValidateGenesis(t *testing.T) {
	tests := []struct {
		name    string
		genesis GenesisState
		wantErr bool
	}{
		{
			name:    "valid genesis with plan and timestamp",
			genesis: GenesisState{Plan: upgradetypes.Plan{Name: "test", Height: 1}, Timestamp: &prototypes.Timestamp{Seconds: 1}},
			wantErr: false,
		},
		{
			name:    "invalid genesis with empty plan",
			genesis: GenesisState{Timestamp: &prototypes.Timestamp{Seconds: 1}},
			wantErr: true,
		},
		{
			name:    "invalid genesis with empty timestamp",
			genesis: GenesisState{Plan: upgradetypes.Plan{Name: "test", Height: 1}},
			wantErr: true,
		},
		{
			name:    "valid empty genesis",
			genesis: GenesisState{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.genesis.ValidateGenesis()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
