package types_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

func TestGovernorTestEquivalent(t *testing.T) {
	val1 := newGovernor(t, valAddr1)
	val2 := newGovernor(t, valAddr1)
	require.Equal(t, val1.String(), val2.String())

	val2 = newGovernor(t, valAddr2)
	require.NotEqual(t, val1.String(), val2.String())
}

func TestUpdateDescription(t *testing.T) {
	d1 := types.Description{
		Website: "https://governor.cosmos",
		Details: "Test governor",
	}

	d2 := types.Description{
		Moniker:  types.DoNotModifyDesc,
		Identity: types.DoNotModifyDesc,
		Website:  types.DoNotModifyDesc,
		Details:  types.DoNotModifyDesc,
	}

	d3 := types.Description{
		Moniker:  "",
		Identity: "",
		Website:  "",
		Details:  "",
	}

	d, err := d1.UpdateDescription(d2)
	require.Nil(t, err)
	require.Equal(t, d, d1)

	d, err = d1.UpdateDescription(d3)
	require.Nil(t, err)
	require.Equal(t, d, d3)
}

func TestShareTokens(t *testing.T) {
	governor := mkGovernor(100, sdk.NewDec(100))
	assert.True(sdk.DecEq(t, sdk.NewDec(50), governor.TokensFromShares(sdk.NewDec(50))))

	governor.Tokens = sdk.NewInt(50)
	assert.True(sdk.DecEq(t, sdk.NewDec(25), governor.TokensFromShares(sdk.NewDec(50))))
	assert.True(sdk.DecEq(t, sdk.NewDec(5), governor.TokensFromShares(sdk.NewDec(10))))
}

func TestRemoveTokens(t *testing.T) {
	governor := mkGovernor(100, sdk.NewDec(100))

	// remove tokens and test check everything
	governor = governor.RemoveTokens(sdk.NewInt(10))
	require.Equal(t, int64(90), governor.Tokens.Int64())

	// update governor to from bonded -> unbonded
	governor = governor.UpdateStatus(types.Unbonded)
	require.Equal(t, types.Unbonded, governor.Status)

	governor = governor.RemoveTokens(sdk.NewInt(10))
	require.Panics(t, func() { governor.RemoveTokens(sdk.NewInt(-1)) })
	require.Panics(t, func() { governor.RemoveTokens(sdk.NewInt(100)) })
}

func TestAddTokensGovernorBonded(t *testing.T) {
	governor := newGovernor(t, valAddr1)
	governor = governor.UpdateStatus(types.Bonded)
	governor, delShares := governor.AddTokensFromDel(sdk.NewInt(10))

	assert.True(sdk.DecEq(t, sdk.NewDec(10), delShares))
	assert.True(sdk.IntEq(t, sdk.NewInt(10), governor.BondedTokens()))
	assert.True(sdk.DecEq(t, sdk.NewDec(10), governor.DelegatorShares))
}

func TestAddTokensGovernorUnbonding(t *testing.T) {
	governor := newGovernor(t, valAddr1)
	governor = governor.UpdateStatus(types.Unbonding)
	governor, delShares := governor.AddTokensFromDel(sdk.NewInt(10))

	assert.True(sdk.DecEq(t, sdk.NewDec(10), delShares))
	assert.Equal(t, types.Unbonding, governor.Status)
	assert.True(sdk.IntEq(t, sdk.NewInt(10), governor.Tokens))
	assert.True(sdk.DecEq(t, sdk.NewDec(10), governor.DelegatorShares))
}

func TestAddTokensGovernorUnbonded(t *testing.T) {
	governor := newGovernor(t, valAddr1)
	governor = governor.UpdateStatus(types.Unbonded)
	governor, delShares := governor.AddTokensFromDel(sdk.NewInt(10))

	assert.True(sdk.DecEq(t, sdk.NewDec(10), delShares))
	assert.Equal(t, types.Unbonded, governor.Status)
	assert.True(sdk.IntEq(t, sdk.NewInt(10), governor.Tokens))
	assert.True(sdk.DecEq(t, sdk.NewDec(10), governor.DelegatorShares))
}

// TODO refactor to make simpler like the AddToken tests above
func TestRemoveDelShares(t *testing.T) {
	valA := types.Governor{
		OperatorAddress: valAddr1.String(),
		Status:          types.Bonded,
		Tokens:          sdk.NewInt(100),
		DelegatorShares: sdk.NewDec(100),
	}

	// Remove delegator shares
	valB, coinsB := valA.RemoveDelShares(sdk.NewDec(10))
	require.Equal(t, int64(10), coinsB.Int64())
	require.Equal(t, int64(90), valB.DelegatorShares.RoundInt64())
	require.Equal(t, int64(90), valB.BondedTokens().Int64())

	// specific case from random tests
	governor := mkGovernor(5102, sdk.NewDec(115))
	_, tokens := governor.RemoveDelShares(sdk.NewDec(29))

	require.True(sdk.IntEq(t, sdk.NewInt(1286), tokens))
}

func TestAddTokensFromDel(t *testing.T) {
	governor := newGovernor(t, valAddr1)

	governor, shares := governor.AddTokensFromDel(sdk.NewInt(6))
	require.True(sdk.DecEq(t, sdk.NewDec(6), shares))
	require.True(sdk.DecEq(t, sdk.NewDec(6), governor.DelegatorShares))
	require.True(sdk.IntEq(t, sdk.NewInt(6), governor.Tokens))

	governor, shares = governor.AddTokensFromDel(sdk.NewInt(3))
	require.True(sdk.DecEq(t, sdk.NewDec(3), shares))
	require.True(sdk.DecEq(t, sdk.NewDec(9), governor.DelegatorShares))
	require.True(sdk.IntEq(t, sdk.NewInt(9), governor.Tokens))
}

func TestUpdateStatus(t *testing.T) {
	governor := newGovernor(t, valAddr1)
	governor, _ = governor.AddTokensFromDel(sdk.NewInt(100))
	require.Equal(t, types.Unbonded, governor.Status)
	require.Equal(t, int64(100), governor.Tokens.Int64())

	// Unbonded to Bonded
	governor = governor.UpdateStatus(types.Bonded)
	require.Equal(t, types.Bonded, governor.Status)

	// Bonded to Unbonding
	governor = governor.UpdateStatus(types.Unbonding)
	require.Equal(t, types.Unbonding, governor.Status)

	// Unbonding to Bonded
	governor = governor.UpdateStatus(types.Bonded)
	require.Equal(t, types.Bonded, governor.Status)
}

func TestPossibleOverflow(t *testing.T) {
	delShares := sdk.NewDec(391432570689183511).Quo(sdk.NewDec(40113011844664))
	governor := mkGovernor(2159, delShares)
	newGovernor, _ := governor.AddTokensFromDel(sdk.NewInt(71))

	require.False(t, newGovernor.DelegatorShares.IsNegative())
	require.False(t, newGovernor.Tokens.IsNegative())
}

func TestGovernorMarshalUnmarshalJSON(t *testing.T) {
	governor := newGovernor(t, valAddr1)
	js, err := legacy.Cdc.MarshalJSON(governor)
	require.NoError(t, err)
	require.NotEmpty(t, js)
	require.Contains(t, string(js), "\"consensus_pubkey\":{\"type\":\"tendermint/PubKeyEd25519\"")
	got := &types.Governor{}
	err = legacy.Cdc.UnmarshalJSON(js, got)
	assert.NoError(t, err)
	assert.True(t, governor.Equal(got))
}

func TestGovernorSetInitialCommission(t *testing.T) {
	val := newGovernor(t, valAddr1)
	testCases := []struct {
		governor    types.Governor
		commission  types.Commission
		expectedErr bool
	}{
		{val, types.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()), false},
		{val, types.NewCommission(sdk.ZeroDec(), sdk.NewDecWithPrec(-1, 1), sdk.ZeroDec()), true},
		{val, types.NewCommission(sdk.ZeroDec(), sdk.NewDec(15000000000), sdk.ZeroDec()), true},
		{val, types.NewCommission(sdk.NewDecWithPrec(-1, 1), sdk.ZeroDec(), sdk.ZeroDec()), true},
		{val, types.NewCommission(sdk.NewDecWithPrec(2, 1), sdk.NewDecWithPrec(1, 1), sdk.ZeroDec()), true},
		{val, types.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.NewDecWithPrec(-1, 1)), true},
		{val, types.NewCommission(sdk.ZeroDec(), sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(2, 1)), true},
	}

	for i, tc := range testCases {
		val, err := tc.governor.SetInitialCommission(tc.commission)

		if tc.expectedErr {
			require.Error(t, err,
				"expected error for test case #%d with commission: %s", i, tc.commission,
			)
		} else {
			require.NoError(t, err,
				"unexpected error for test case #%d with commission: %s", i, tc.commission,
			)
			require.Equal(t, tc.commission, val.Commission,
				"invalid governor commission for test case #%d with commission: %s", i, tc.commission,
			)
		}
	}
}

// Check that sort will create deterministic ordering of governors
func TestGovernorsSortDeterminism(t *testing.T) {
	vals := make([]types.Governor, 10)
	sortedVals := make([]types.Governor, 10)

	// Create random governor slice
	for i := range vals {
		pk := ed25519.GenPrivKey().PubKey()
		vals[i] = newGovernor(t, sdk.ValAddress(pk.Address()))
	}

	// Save sorted copy
	sort.Sort(types.Governors(vals))
	copy(sortedVals, vals)

	// Randomly shuffle governors, sort, and check it is equal to original sort
	for i := 0; i < 10; i++ {
		rand.Shuffle(10, func(i, j int) {
			it := vals[i]
			vals[i] = vals[j]
			vals[j] = it
		})

		types.Governors(vals).Sort()
		require.Equal(t, sortedVals, vals, "Governor sort returned different slices")
	}
}

func TestBondStatus(t *testing.T) {
	require.False(t, types.Unbonded == types.Bonded)
	require.False(t, types.Unbonded == types.Unbonding)
	require.False(t, types.Bonded == types.Unbonding)
	require.Equal(t, types.BondStatus(4).String(), "4")
	require.Equal(t, types.BondStatusUnspecified, types.Unspecified.String())
	require.Equal(t, types.BondStatusUnbonded, types.Unbonded.String())
	require.Equal(t, types.BondStatusBonded, types.Bonded.String())
	require.Equal(t, types.BondStatusUnbonding, types.Unbonding.String())
}

func mkGovernor(tokens int64, shares sdk.Dec) types.Governor {
	return types.Governor{
		OperatorAddress: valAddr1.String(),
		Status:          types.Bonded,
		Tokens:          sdk.NewInt(tokens),
		DelegatorShares: shares,
	}
}

// Creates a new governors and asserts the error check.
func newGovernor(t *testing.T, operator sdk.ValAddress) types.Governor {
	v, err := types.NewGovernor(operator, types.Description{})
	require.NoError(t, err)
	return v
}
