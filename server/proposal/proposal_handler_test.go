package proposal

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/proto/tendermint/types"

	rollappparamstypes "github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"
)

func TestCustomParameterChangeProposalHandler(t *testing.T) {
	logger := log.NewNopLogger()
	ctx := sdk.NewContext(nil, tmtypes.Header{}, false, logger)
	hasDenom := func(ctx sdk.Context, denom string) bool { return true }
	hasNotDenom := func(ctx sdk.Context, denom string) bool { return false }

	tests := []struct {
		name        string
		changes     []proposal.ParamChange
		getSubspace getSubspaceFn
		hasDenom    hasDenomMetaDataFn
		expectErr   bool
		errContains string
	}{
		{
			name: "non rollappparams subspace (no special checks)",
			changes: []proposal.ParamChange{
				{
					Subspace: "othersubspace",
					Key:      "someKey",
					Value:    "someValue",
				},
			},
			getSubspace: func(subspace string) (paramsSubspace, bool) {
				return mockSubspace{}, true
			},
			hasDenom:  nil,
			expectErr: false,
		},
		{
			name: "rollappparams but different key than minGasPrices",
			changes: []proposal.ParamChange{
				{
					Subspace: rollappparamstypes.ModuleName,
					Key:      "someOtherKey",
					Value:    "something",
				},
			},
			getSubspace: func(subspace string) (paramsSubspace, bool) {
				return mockSubspace{}, true
			},
			hasDenom:  nil,
			expectErr: false,
		},
		{
			name: "rollappparams minGasPrices - valid denoms",
			changes: []proposal.ParamChange{
				{
					Subspace: rollappparamstypes.ModuleName,
					Key:      string(rollappparamstypes.KeyMinGasPrices),
					Value:    `[{"denom":"adenom","amount":"20000000000.0"}]`,
				},
			},
			getSubspace: func(subspace string) (paramsSubspace, bool) {
				return mockSubspace{}, true
			},
			hasDenom:  hasDenom,
			expectErr: false,
		},
		{
			name: "rollappparams minGasPrices - unknown denom",
			changes: []proposal.ParamChange{
				{
					Subspace: rollappparamstypes.ModuleName,
					Key:      string(rollappparamstypes.KeyMinGasPrices),
					Value:    `[{"denom":"baddenom","amount":"20000000000.0"}]`,
				},
			},
			getSubspace: func(subspace string) (paramsSubspace, bool) {
				return mockSubspace{}, true
			},
			hasDenom:    hasNotDenom,
			expectErr:   true,
			errContains: "denom baddenom does not exist",
		},
		{
			name: "rollappparams minGasPrices - invalid JSON",
			changes: []proposal.ParamChange{
				{
					Subspace: rollappparamstypes.ModuleName,
					Key:      string(rollappparamstypes.KeyMinGasPrices),
					Value:    `not-json`,
				},
			},
			getSubspace: func(subspace string) (paramsSubspace, bool) {
				return mockSubspace{}, true
			},
			hasDenom:    nil,
			expectErr:   true,
			errContains: "failed to unmarshal minGasPrices",
		},
		{
			name: "unknown subspace",
			changes: []proposal.ParamChange{
				{
					Subspace: "unknownsubspace",
					Key:      "someKey",
					Value:    "someValue",
				},
			},
			getSubspace: func(subspace string) (paramsSubspace, bool) {
				return mockSubspace{}, false
			},
			hasDenom:    hasDenom,
			expectErr:   true,
			errContains: "unknown subspace",
		},
		{
			name: "error on update",
			changes: []proposal.ParamChange{
				{
					Subspace: rollappparamstypes.ModuleName,
					Key:      "someKey",
					Value:    "someValue",
				},
			},
			getSubspace: func(subspace string) (paramsSubspace, bool) {
				return mockSubspace{
					updateErr: errors.New("update failed"),
				}, true
			},
			hasDenom:    hasDenom,
			expectErr:   true,
			errContains: "update failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := handleCustomParameterChangeProposal(ctx, tc.getSubspace, logger.Info, tc.hasDenom, tc.changes)
			if tc.expectErr {
				require.Error(t, err)
				if tc.errContains != "" {
					require.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type mockSubspace struct {
	updateErr error
}

func (m mockSubspace) Update(sdk.Context, []byte, []byte) error {
	return m.updateErr
}
