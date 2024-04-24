package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/dymensionxyz/dymension-rdk/x/governors/types"
)

// StakingKeeper expected staking keeper (noalias)
type SequencerKeeper interface {
	GetSequencerByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (stakingtypes.Governor, bool)
}

// StakingKeeper expected staking keeper (noalias)
type StakingKeeper interface {
	GetLastTotalPower(ctx sdk.Context) math.Int

	// iterate through governors by operator address, execute func for each governor
	IterateGovernors(sdk.Context,
		func(index int64, governor stakingtypes.GovernorI) (stop bool))
	// iterate through bonded governors by operator address, execute func for each governor
	IterateBondedGovernorsByPower(sdk.Context,
		func(index int64, governor stakingtypes.GovernorI) (stop bool))

	Governor(sdk.Context, sdk.ValAddress) stakingtypes.GovernorI // get a particular governor by operator address

	// Delegation allows for getting a particular delegation for a given governor
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI

	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))

	GetAllSDKDelegations(ctx sdk.Context) []stakingtypes.Delegation
}
