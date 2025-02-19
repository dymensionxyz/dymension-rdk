package utils

import "cosmossdk.io/math"

var (
	// ADYM represents 1 ADYM.
	// It is not real DYM, just a number representation for convenience.
	ADYM = math.NewInt(1)

	// DYM represents 1 DYM. Equals to 10^18 ADYM.
	// It is not real DYM, just a number representation for convenience.
	DYM = math.NewIntWithDecimal(1, 18)
)
