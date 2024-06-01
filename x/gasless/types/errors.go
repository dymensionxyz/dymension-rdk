package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// DONTCOVER

var (
	ErrorFeeConsumptionFailure = sdkerrors.Register(ModuleName, 10001, "fee cannot be deducted from gas tank")
)
