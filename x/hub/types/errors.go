package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/hub module sentinel errors
var (
	ErrGenesisEventNotTriggered = errorsmod.Register(ModuleName, 1000, "genesis event not triggered yet")
	ErrMismatchedChannelID      = errorsmod.Register(ModuleName, 1001, "mismatched channel id")
)
