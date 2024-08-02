package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/sequencers module sentinel errors
var (
	ErrDANotSupported = errorsmod.Register(ModuleName, 1, "da type not supported")
)
