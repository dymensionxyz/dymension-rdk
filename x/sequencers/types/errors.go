package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sequencers module sentinel errors
var (
	ErrMultipleDymintSequencers = sdkerrors.Register(ModuleName, 1, "multiple dymint sequencers not supported")
	ErrSequencerNotFound        = sdkerrors.Register(ModuleName, 2, "sequencer address not found")
	ErrHistoricalInfoNotFound   = sdkerrors.Register(ModuleName, 3, "historical info not found")
)
