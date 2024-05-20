package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// DONTCOVER

var (
	ErrorUnknownProposalType       = sdkerrors.Register(ModuleName, 10000, "unknown proposal type")
	ErrorInvalidrequest            = sdkerrors.Register(ModuleName, 10001, "invalid request")
	ErrorMaxLimitReachedByProvider = sdkerrors.Register(ModuleName, 10002, "provider reached maximum limit to create gas tanks")
	ErrorFeeConsumptionFailure     = sdkerrors.Register(ModuleName, 10003, "fee cannot be deducted from gas tank")
)
