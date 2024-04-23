package types

import (
	errorsmod "cosmossdk.io/errors"
)

// errors
var (
	ErrDuplicate          = errorsmod.Register(ModuleName, 301, "duplicate")
	ErrNoPermission       = errorsmod.Register(ModuleName, 302, "signer not in the permissions list to create vesting account")
)