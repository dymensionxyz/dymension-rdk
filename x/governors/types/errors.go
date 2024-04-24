package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/staking module sentinel errors
//
// TODO: Many of these errors are redundant. They should be removed and replaced
// by sdkerrors.ErrInvalidRequest.
//
// REF: https://github.com/cosmos/cosmos-sdk/issues/5450
var (
	ErrEmptyGovernorAddr              = sdkerrors.Register(ModuleName, 2, "empty Governor address")
	ErrNoGovernorFound                = sdkerrors.Register(ModuleName, 3, "Governor does not exist")
	ErrGovernorOwnerExists            = sdkerrors.Register(ModuleName, 4, "Governor already exist for this operator address; must use new Governor operator address")
	ErrGovernorPubKeyExists           = sdkerrors.Register(ModuleName, 5, "Governor already exist for this pubkey; must use new Governor pubkey")
	ErrGovernorPubKeyTypeNotSupported = sdkerrors.Register(ModuleName, 6, "Governor pubkey type is not supported")
	ErrGovernorJailed                 = sdkerrors.Register(ModuleName, 7, "Governor for this address is currently jailed")
	ErrBadRemoveGovernor              = sdkerrors.Register(ModuleName, 8, "failed to remove Governor")
	ErrCommissionNegative             = sdkerrors.Register(ModuleName, 9, "commission must be positive")
	ErrCommissionHuge                 = sdkerrors.Register(ModuleName, 10, "commission cannot be more than 100%")
	ErrCommissionGTMaxRate            = sdkerrors.Register(ModuleName, 11, "commission cannot be more than the max rate")
	ErrCommissionUpdateTime           = sdkerrors.Register(ModuleName, 12, "commission cannot be changed more than once in 24h")
	ErrCommissionChangeRateNegative   = sdkerrors.Register(ModuleName, 13, "commission change rate must be positive")
	ErrCommissionChangeRateGTMaxRate  = sdkerrors.Register(ModuleName, 14, "commission change rate cannot be more than the max rate")
	ErrCommissionGTMaxChangeRate      = sdkerrors.Register(ModuleName, 15, "commission cannot be changed more than max change rate")
	ErrSelfDelegationBelowMinimum     = sdkerrors.Register(ModuleName, 16, "Governor's self delegation must be greater than their minimum self delegation")
	ErrMinSelfDelegationDecreased     = sdkerrors.Register(ModuleName, 17, "minimum self delegation cannot be decrease")
	ErrEmptyDelegatorAddr             = sdkerrors.Register(ModuleName, 18, "empty delegator address")
	ErrNoDelegation                   = sdkerrors.Register(ModuleName, 19, "no delegation for (address, Governor) tuple")
	ErrBadDelegatorAddr               = sdkerrors.Register(ModuleName, 20, "delegator does not exist with address")
	ErrNoDelegatorForAddress          = sdkerrors.Register(ModuleName, 21, "delegator does not contain delegation")
	ErrInsufficientShares             = sdkerrors.Register(ModuleName, 22, "insufficient delegation shares")
	ErrDelegationGovernorEmpty        = sdkerrors.Register(ModuleName, 23, "cannot delegate to an empty Governor")
	ErrNotEnoughDelegationShares      = sdkerrors.Register(ModuleName, 24, "not enough delegation shares")
	ErrNotMature                      = sdkerrors.Register(ModuleName, 25, "entry not mature")
	ErrNoUnbondingDelegation          = sdkerrors.Register(ModuleName, 26, "no unbonding delegation found")
	ErrMaxUnbondingDelegationEntries  = sdkerrors.Register(ModuleName, 27, "too many unbonding delegation entries for (delegator, Governor) tuple")
	ErrNoRedelegation                 = sdkerrors.Register(ModuleName, 28, "no redelegation found")
	ErrSelfRedelegation               = sdkerrors.Register(ModuleName, 29, "cannot redelegate to the same Governor")
	ErrTinyRedelegationAmount         = sdkerrors.Register(ModuleName, 30, "too few tokens to redelegate (truncates to zero tokens)")
	ErrBadRedelegationDst             = sdkerrors.Register(ModuleName, 31, "redelegation destination Governor not found")
	ErrTransitiveRedelegation         = sdkerrors.Register(ModuleName, 32, "redelegation to this Governor already in progress; first redelegation to this Governor must complete before next redelegation")
	ErrMaxRedelegationEntries         = sdkerrors.Register(ModuleName, 33, "too many redelegation entries for (delegator, src-Governor, dst-Governor) tuple")
	ErrDelegatorShareExRateInvalid    = sdkerrors.Register(ModuleName, 34, "cannot delegate to Governors with invalid (zero) ex-rate")
	ErrBothShareMsgsGiven             = sdkerrors.Register(ModuleName, 35, "both shares amount and shares percent provided")
	ErrNeitherShareMsgsGiven          = sdkerrors.Register(ModuleName, 36, "neither shares amount nor shares percent provided")
	ErrCommissionLTMinRate            = sdkerrors.Register(ModuleName, 40, "commission cannot be less than min rate")
)
