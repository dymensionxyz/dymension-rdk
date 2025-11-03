package types

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

// ConvertAmount converts an amount from one token to another based on their decimal precisions
// For example, converting from 6 decimals to 18 decimals: amount * 10^(18-6) = amount * 10^12
// Converting from 18 decimals to 6 decimals: amount / 10^(18-6) = amount / 10^12
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
