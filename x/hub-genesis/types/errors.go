package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrFailedGetClientState         = errorsmod.Register(ModuleName, 1001, "failed to get client state")
	ErrChainIDMismatch              = errorsmod.Register(ModuleName, 1002, "chain ID not matches with the channel")
	ErrInvalidGenesisTokens         = errorsmod.Register(ModuleName, 1003, "invalid genesis token")
	ErrGenesisEventAlreadyTriggered = errorsmod.Register(ModuleName, 1004, "genesis event already triggered")
	ErrGenesisInsufficentBalance    = errorsmod.Register(ModuleName, 1005, "insufficent balance in module account to lock genesis tokens")
	ErrLockingGenesisTokens         = errorsmod.Register(ModuleName, 1006, "failed to lock genesis tokens")
)
