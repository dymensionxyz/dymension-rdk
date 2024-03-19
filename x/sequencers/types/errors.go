package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/sequencers module sentinel errors
var (
	ErrNoSequencerOnGenesis   = errorsmod.Register(ModuleName, 1, "sequencer on genesis required")
	ErrSequencerNotFound      = errorsmod.Register(ModuleName, 2, "sequencer address not found")
	ErrHistoricalInfoNotFound = errorsmod.Register(ModuleName, 3, "historical info not found")

	ErrEmptyValidatorAddr              = errorsmod.Register(ModuleName, 100, "empty validator address")
	ErrNoValidatorFound                = errorsmod.Register(ModuleName, 101, "validator does not exist")
	ErrValidatorOwnerExists            = errorsmod.Register(ModuleName, 102, "validator already exist for this operator address; must use new validator operator address")
	ErrValidatorPubKeyExists           = errorsmod.Register(ModuleName, 103, "validator already exist for this pubkey; must use new validator pubkey")
	ErrValidatorPubKeyTypeNotSupported = errorsmod.Register(ModuleName, 104, "validator pubkey type is not supported")
	ErrEmptyDelegatorAddr              = errorsmod.Register(ModuleName, 105, "empty delegator address")
	ErrEmptyValidatorPubKey            = errorsmod.Register(ModuleName, 106, "empty validator public key")
	ErrSequencerNotRegistered          = errorsmod.Register(ModuleName, 107, "sequencer not registered on the hub")
)
