package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrDuplicate          = errorsmod.Register(ModuleName, 200, "duplicate")
	ErrBlank              = errorsmod.Register(ModuleName, 201, "address cannot be blank")
	ErrNoPermission       = errorsmod.Register(ModuleName, 202, "signer not in the permissions list to create or update denom metadata")
	ErrDenomAlreadyExists = errorsmod.Register(ModuleName, 203, "denom metadata is already registered")
	ErrDenomDoesNotExist  = errorsmod.Register(ModuleName, 204, "unable to find denom metadata registered")
)
