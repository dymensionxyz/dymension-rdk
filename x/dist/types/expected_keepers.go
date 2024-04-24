package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper expected staking keeper (noalias)
type SequencerKeeper interface {
	GetSequencerByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (stakingtypes.ValidatorI, bool)
}

// StakingKeeper expected staking keeper (noalias)
type StakingKeeper interface {
	GetLastTotalPower(ctx sdk.Context) math.Int

	// iterate through governors by operator address, execute func for each governor
	IterateValidators(sdk.Context,
		func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	// iterate through bonded governors by operator address, execute func for each governor
	IterateBondedValidatorsByPower(sdk.Context,
		func(index int64, validator stakingtypes.ValidatorI) (stop bool))

	Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI // get a particular governor by operator address

	// Delegation allows for getting a particular delegation for a given governor
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI

	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))

	GetAllSDKDelegations(ctx sdk.Context) []stakingtypes.Delegation

	//unused
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI
}
