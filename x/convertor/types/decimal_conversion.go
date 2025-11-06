package types

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

// ConvertAmount converts an amount from one token to another based on their decimal precisions.
//
// Scaling up (e.g., 6 decimals → 18 decimals):
//   - amount * 10^(18-6) = amount * 10^12
//   - Example: 1000000 (1 USDC) → 1000000000000000000 (1 USDC in 18 decimals)
//   - No precision loss occurs when scaling up
//
// Scaling down (e.g., 18 decimals → 6 decimals):
//   - amount / 10^(18-6) = amount / 10^12
//   - Example: 1000000000000000001 → 1000000 (1 wei is permanently lost)
//   - ⚠️  PRECISION LOSS: Any fractional amount smaller than the target precision is truncated
//   - The lost precision cannot be recovered and represents a permanent loss of value
//   - Users should be aware that sending amounts with precision beyond the target decimals
//     will result in the excess precision being discarded
//
// Best Practice:
//   - When sending from rollapp (18 decimals) to hub (6 decimals), ensure amounts are
//     multiples of 10^12 to avoid precision loss
//   - Example: Send 1000000000000000000 (safe) instead of 1000000000000000001 (loses 1 wei)
func ConvertAmount(amount math.Int, fromDecimals, toDecimals uint32) (math.Int, error) {
	if fromDecimals == toDecimals {
		return amount, nil
	}

	if amount.IsNegative() {
		return math.Int{}, errorsmod.Wrap(gerrc.ErrInvalidArgument, "amount cannot be negative")
	}

	amountBig := amount.BigInt()

	if fromDecimals < toDecimals {
		// Scaling up: multiply by 10^(toDecimals - fromDecimals)
		decimalDiff := toDecimals - fromDecimals
		multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimalDiff)), nil)
		result := new(big.Int).Mul(amountBig, multiplier)
		return sdk.NewIntFromBigInt(result), nil
	}

	// Scaling down: divide by 10^(fromDecimals - toDecimals)
	decimalDiff := fromDecimals - toDecimals
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimalDiff)), nil)
	result := new(big.Int).Div(amountBig, divisor)

	return sdk.NewIntFromBigInt(result), nil
}

// ClearPrecisionLoss removes precision that would be lost during decimal conversion by
// performing a round-trip conversion (down then up) and returning the result.
//
// This function is useful for calculating "dust" - the amount that would be lost when
// converting from higher to lower decimal precision. The difference between the original
// amount and the result of this function is the dust that must be burned.
//
// Example with 18 -> 6 decimals:
//   - Input:  1000000000000000001 (1 token + 1 wei with 18 decimals)
//   - Output: 1000000000000000000 (1 token with 18 decimals, dust cleared)
//   - Dust:   1000000000000000001 - 1000000000000000000 = 1 wei
//
// The function only clears precision when scaling down (fromDecimals > toDecimals).
// When scaling up or when decimals are equal, the amount is returned unchanged since
// no precision is lost.
//
// Parameters:
//   - amount: The original amount in fromDecimals precision
//   - fromDecimals: The decimal precision of the original amount
//   - toDecimals: The target decimal precision (typically lower than fromDecimals)
//
// Returns:
//   - The amount with precision loss cleared (in fromDecimals precision)
//   - Error if the conversion fails
func ClearPrecisionLoss(amount math.Int, fromDecimals, toDecimals uint32) (math.Int, error) {
	if fromDecimals <= toDecimals {
		// No precision loss when scaling up or when decimals are equal
		return amount, nil
	}

	// Scale down to lower decimals (this truncates any excess precision)
	convertedAmt, err := ConvertAmount(amount, fromDecimals, toDecimals)
	if err != nil {
		return math.Int{}, err
	}

	// Scale back up to original decimals (precision loss is now cleared)
	return ConvertAmount(convertedAmt, toDecimals, fromDecimals)
}
