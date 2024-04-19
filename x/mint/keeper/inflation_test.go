package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	utils "github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/mint/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInflationChangeTCs(t *testing.T) {
	// we use the same target inflation rate for all test cases
	targetInflationRate := sdk.NewDecWithPrec(5, 2) // 5%

	testCases := []struct {
		name                 string
		currentInflationRate sdk.Dec
		inflationRateChange  sdk.Dec
		expectedInflation    sdk.Dec
	}{
		{
			name:                 "Test Decrease",
			currentInflationRate: sdk.NewDecWithPrec(7, 2), // 7%
			inflationRateChange:  sdk.NewDecWithPrec(1, 2), // 1%
			expectedInflation:    sdk.NewDecWithPrec(6, 2), // 6%
		},
		{
			name:                 "Test Increase",
			currentInflationRate: sdk.NewDecWithPrec(1, 2), // 1%
			inflationRateChange:  sdk.NewDecWithPrec(1, 2), // 1%
			expectedInflation:    sdk.NewDecWithPrec(2, 2), // 2%
		},
		{
			name:                 "Test Decrease - Max",
			currentInflationRate: sdk.NewDecWithPrec(6, 2),  // 6%
			inflationRateChange:  sdk.NewDecWithPrec(15, 2), // 1.5%
			expectedInflation:    targetInflationRate,       // 5%
		},
		{
			name:                 "Test Increase - Max",
			currentInflationRate: sdk.NewDecWithPrec(4, 2),  // 4%
			inflationRateChange:  sdk.NewDecWithPrec(15, 2), // 1.5%
			expectedInflation:    targetInflationRate,       // 5%
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := utils.Setup(t, false)
			k, ctx := testkeepers.NewTestMintKeeperFromApp(app)

			// Set initial minter & params
			minter := types.Minter{
				CurrentInflationRate: tc.currentInflationRate,
			}
			k.SetMinter(ctx, minter)

			params := k.GetParams(ctx)
			params.TargetInflationRate = targetInflationRate
			params.InflationRateChange = tc.inflationRateChange
			k.SetParams(ctx, params)

			_, err := k.HandleInflationChange(ctx)
			require.NoError(t, err)

			minter = k.GetMinter(ctx)
			assert.Equal(t, tc.expectedInflation, minter.CurrentInflationRate)
		})
	}
}
