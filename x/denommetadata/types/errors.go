package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrDuplicate = errorsmod.Register(ModuleName, 1, "duplicate")
	ErrBlank     = errorsmod.Register(ModuleName, 2, "address cannot be blank")
)
