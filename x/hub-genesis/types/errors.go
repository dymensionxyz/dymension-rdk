package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrFailedGetClientState         = errorsmod.Register(ModuleName, 1001, "failed to get client state")
	ErrChainIDMismatch              = errorsmod.Register(ModuleName, 1002, "chain ID not matches with the channel")
	ErrInvalidGenesisTokens         = errorsmod.Register(ModuleName, 1003, "invalid genesis token") // TODO: use it(?) where was it originally used?
	ErrGenesisEventAlreadyTriggered = errorsmod.Register(ModuleName, 1004, "genesis event already triggered")
	ErrLockingGenesisTokens         = errorsmod.Register(ModuleName, 1005, "failed to lock genesis tokens")
)
