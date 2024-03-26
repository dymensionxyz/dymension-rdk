package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/sequencers module sentinel errors
var (
	ErrMultipleDymintSequencers = errorsmod.Register(ModuleName, 1, "multiple dymint sequencers not supported")
	ErrNoSequencerOnInitChain   = errorsmod.Register(ModuleName, 2, "no sequencer defined on InitChain")
	ErrFailedInitChain          = errorsmod.Register(ModuleName, 4, "failed to initialize sequencer on InitChain")
	ErrFailedInitGenesis        = errorsmod.Register(ModuleName, 5, "failed to initialize sequencer on InitGenesis")
	ErrSequencerNotFound        = errorsmod.Register(ModuleName, 6, "sequencer address not found")
	ErrHistoricalInfoNotFound   = errorsmod.Register(ModuleName, 7, "historical info not found")
)
