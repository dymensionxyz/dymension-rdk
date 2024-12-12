package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
)

func TestCheckFeeCoinsAgainstMinGasPrices(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestRollappParamsKeeperFromApp(app)

	tests := []struct {
		name                  string
		feeCoins              sdk.Coins
		gas                   uint64
		validatorMinGasPrices sdk.DecCoins
		globalMinGasPrices    sdk.DecCoins
		expectErr             bool
	}{
		{
			name:                  "valid fee, global is empty, validator is empty",
			feeCoins:              sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(1000))),
			gas:                   1_000_000,
			validatorMinGasPrices: sdk.DecCoins{},
			globalMinGasPrices:    sdk.DecCoins{},
			expectErr:             false,
		},
		{
			name:                  "valid fee, global is empty, validator is not empty",
			feeCoins:              sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(1000))),
			gas:                   1_000_000,
			validatorMinGasPrices: sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 3))}, // 0.001 adym
			globalMinGasPrices:    sdk.DecCoins{},
			expectErr:             false,
		},
		{
			name:                  "insufficient fee, global is empty, validator is not empty",
			feeCoins:              sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(500))),
			gas:                   1_000_000,
			validatorMinGasPrices: sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 3))}, // 0.001 adym
			globalMinGasPrices:    sdk.DecCoins{},
			expectErr:             true,
		},
		{
			name:                  "valid fee, global is not empty, validator is empty",
			feeCoins:              sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(1000))),
			gas:                   1_000_000,
			validatorMinGasPrices: sdk.DecCoins{},
			globalMinGasPrices:    sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 3))}, // 0.001 adym
			expectErr:             false,
		},
		{
			name:                  "insufficient fee, global is not empty, validator is empty",
			feeCoins:              sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(500))),
			gas:                   1_000_000,
			validatorMinGasPrices: sdk.DecCoins{},
			globalMinGasPrices:    sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 3))}, // 0.001 adym
			expectErr:             true,
		},
		{
			name:                  "valid fee, global is not empty, validator is not empty, no intersection",
			feeCoins:              sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(1000))),
			gas:                   1_000_000,
			validatorMinGasPrices: sdk.DecCoins{sdk.NewDecCoinFromDec("uatom", sdk.NewDecWithPrec(1, 3))}, // 0.001 uatom
			globalMinGasPrices:    sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 3))},  // 0.001 adym
			expectErr:             true,
		},
		{
			name:                  "insufficient fee, global is greater than validator",
			feeCoins:              sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(1000))),
			gas:                   1_000_000,
			validatorMinGasPrices: sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 3))}, // 0.001 uatom
			globalMinGasPrices:    sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 2))}, // 0.01 adym
			expectErr:             true,
		},
		{
			name:                  "valid fee, global is greater than validator",
			feeCoins:              sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(10000))),
			gas:                   1_000_000,
			validatorMinGasPrices: sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 3))}, // 0.001 uatom
			globalMinGasPrices:    sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 2))}, // 0.01 adym
			expectErr:             false,
		},
		{
			name:                  "insufficient fee, validator is greater than global",
			feeCoins:              sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(1000))),
			gas:                   1_000_000,
			validatorMinGasPrices: sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 2))}, // 0.01 uatom
			globalMinGasPrices:    sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 3))}, // 0.001 adym
			expectErr:             true,
		},
		{
			name:                  "valid fee, validator is greater than global",
			feeCoins:              sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(10000))),
			gas:                   1_000_000,
			validatorMinGasPrices: sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 2))}, // 0.01 uatom
			globalMinGasPrices:    sdk.DecCoins{sdk.NewDecCoinFromDec("adym", sdk.NewDecWithPrec(1, 3))}, // 0.001 adym
			expectErr:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx = ctx.WithMinGasPrices(tt.validatorMinGasPrices)
			err := k.SetMinGasPrices(ctx, tt.globalMinGasPrices)
			require.NoError(t, err)

			err = k.CheckFeeCoinsAgainstMinGasPrices(ctx, tt.feeCoins, tt.gas)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
