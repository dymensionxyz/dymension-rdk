package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sequencers module sentinel errors
var (
	ErrMultipleDymintSequencers = sdkerrors.Register(ModuleName, 1, "multiple dymint sequencers not supported")
	ErrNoSequencerOnInitChain   = sdkerrors.Register(ModuleName, 2, "no sequencer defined on InitChain")
	ErrFailedInitChain          = sdkerrors.Register(ModuleName, 3, "failed to initialize sequencer on InitChain")
	ErrSequencerNotFound        = sdkerrors.Register(ModuleName, 4, "sequencer address not found")
	ErrHistoricalInfoNotFound   = sdkerrors.Register(ModuleName, 5, "historical info not found")
)
