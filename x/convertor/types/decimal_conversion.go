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
