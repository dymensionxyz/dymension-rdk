package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrSequencerNotFound           = errorsmod.Register(ModuleName, 6, "sequencer address not found")
	ErrHistoricalInfoNotFound      = errorsmod.Register(ModuleName, 7, "historical info not found")
	ErrWhitelistedRelayersNotFound = errorsmod.Register(ModuleName, 8, "whitelisted relayers not found")
)
