package teststaking

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

// NewGovernor is a testing helper method to create governors in tests
func NewGovernor(t testing.TB, operator sdk.ValAddress) types.Governor {
	v, err := types.NewGovernor(operator, types.Description{})
	require.NoError(t, err)
	return v
}
