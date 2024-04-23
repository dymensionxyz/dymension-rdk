package types

import (
	errorsmod "cosmossdk.io/errors"
)

// errors
var (
	ErrDuplicate      = errorsmod.Register(ModuleName, 301, "duplicate")
	ErrInvalidSigners = errorsmod.Register(ModuleName, 302, "signers for vesting tx should be 1")
	ErrNoPermission   = errorsmod.Register(ModuleName, 303, "signer not in the permissions list to create vesting account")
)
