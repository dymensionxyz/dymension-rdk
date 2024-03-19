package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrUnknownHubID                 = sdkerrors.Register(ModuleName, 1000, "unknown hub id")
	ErrInvalidGenesisChannelId      = sdkerrors.Register(ModuleName, 1001, "invalid genesis channel id")
	ErrInvalidGenesisChainId        = sdkerrors.Register(ModuleName, 1002, "invalid genesis chain id")
	ErrGenesisEventAlreadyTriggered = sdkerrors.Register(ModuleName, 1003, "genesis event already triggered")
	ErrGenesisNoCoinsOnModuleAcc    = sdkerrors.Register(ModuleName, 1004, "no coins on module account")
	ErrLockingGenesisTokens         = sdkerrors.Register(ModuleName, 1005, "failed to lock genesis tokens")
)
