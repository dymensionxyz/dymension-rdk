package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrDANotSupported = errorsmod.Register(ModuleName, 1, "da type not supported")
)
