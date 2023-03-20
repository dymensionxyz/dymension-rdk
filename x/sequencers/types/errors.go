package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sequencers module sentinel errors
var (
	ErrNoSequencerOnGenesis   = sdkerrors.Register(ModuleName, 1, "sequencer on genesis required")
	ErrSequencerNotFound      = sdkerrors.Register(ModuleName, 2, "sequencer address not found")
	ErrHistoricalInfoNotFound = sdkerrors.Register(ModuleName, 3, "historical info not found")

	ErrEmptyValidatorAddr              = sdkerrors.Register(ModuleName, 100, "empty validator address")
	ErrNoValidatorFound                = sdkerrors.Register(ModuleName, 101, "validator does not exist")
	ErrValidatorOwnerExists            = sdkerrors.Register(ModuleName, 102, "validator already exist for this operator address; must use new validator operator address")
	ErrValidatorPubKeyExists           = sdkerrors.Register(ModuleName, 103, "validator already exist for this pubkey; must use new validator pubkey")
	ErrValidatorPubKeyTypeNotSupported = sdkerrors.Register(ModuleName, 104, "validator pubkey type is not supported")
	ErrEmptyDelegatorAddr              = sdkerrors.Register(ModuleName, 105, "empty delegator address")
	ErrEmptyValidatorPubKey            = sdkerrors.Register(ModuleName, 106, "empty validator public key")
	ErrSequencerNotRegistered          = sdkerrors.Register(ModuleName, 107, "sequencer not registered on the hub")
)
