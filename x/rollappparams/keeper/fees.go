package keeper

import (
	"fmt"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	"github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"
)

// CheckTxFeeWithMinGasPrices implements the default fee logic, where the minimum price per
// unit of gas is fixed and set by each validator, and the tx priority is computed from the gas price.
func (k Keeper) CheckTxFeeWithMinGasPrices() ante.TxFeeChecker {
	return func(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
		feeTx, ok := tx.(sdk.FeeTx)
		if !ok {
			return nil, 0, errorsmod.Wrapf(errortypes.ErrTxDecode, "Tx must be a FeeTx")
		}

		feeCoins := feeTx.GetFee()
		gas := feeTx.GetGas()

		// Ensure that the provided fees meets a minimum threshold for the validator,
		// if this is a CheckTx. This is only for local mempool purposes, and thus
		// is only ran on CheckTx.
		if ctx.IsCheckTx() {
			if err := k.CheckFeeCoinsAgainstMinGasPrices(ctx, feeCoins, gas); err != nil {
				return nil, 0, err
			}
		}

		priority := getTxPriority(feeCoins, int64(gas))
		return feeCoins, priority, nil
	}
}

// CheckFeeCoinsAgainstMinGasPrices checks if the provided fee coins are greater than or equal to the
// required fees, that are based on the minimum gas prices and the gas. If not, it will return an error.
func (k Keeper) CheckFeeCoinsAgainstMinGasPrices(ctx sdk.Context, feeCoins sdk.Coins, gas uint64) error {
	var (
		validatorMinGasPrices = ctx.MinGasPrices()
		globalMinGasPrices    = k.GetParams(ctx).MinGasPrices
		minGasPrices          sdk.DecCoins

		validatorEmpty = validatorMinGasPrices.IsZero()
		globalEmpty    = globalMinGasPrices.IsZero()
	)

	// Possible cases:
	//
	//  | Global    | Validator | Result                       |
	//  |-----------|-----------|------------------------------|
	//  | empty     | empty     | all txs are accepted         |
	//  | empty     | non-empty | validator values             |
	//  | non-empty | empty     | global values                |
	//  | non-empty | non-empty | intersect(global, validator) |

	switch {
	case globalEmpty && validatorEmpty:
		return nil

	case globalEmpty:
		minGasPrices = validatorMinGasPrices // self assignment for clarity

	case validatorEmpty:
		minGasPrices = globalMinGasPrices

	default:
		minGasPrices = types.IntersectMinGasPrices(validatorMinGasPrices, globalMinGasPrices)

		if minGasPrices.IsZero() {
			return fmt.Errorf("validator has not specified any gas denoms indicated by governance: validator min gas prices: %s, governance min gas prices: %s", validatorMinGasPrices, globalMinGasPrices)
		}
	}

	requiredFees := make(sdk.Coins, len(minGasPrices))

	// Determine the required fees by multiplying each required minimum gas
	// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
	glDec := sdk.NewDec(int64(gas))
	for i, gp := range minGasPrices {
		fee := gp.Amount.Mul(glDec)
		requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
	}

	if !feeCoins.IsAnyGTE(requiredFees) {
		return errorsmod.Wrapf(errortypes.ErrInsufficientFee, "got: %s required: %s", feeCoins, requiredFees)
	}

	return nil
}

// getTxPriority returns a naive tx priority based on the amount of the smallest denomination of the gas price
// provided in a transaction.
// NOTE: This implementation should be used with a great consideration as it opens potential attack vectors
// where txs with multiple coins could not be prioritized as expected.
func getTxPriority(fees sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fees {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}

	return priority
}
