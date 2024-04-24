package types

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GovernorSet expected properties for the set of all governors (noalias)
type GovernorSet interface {
	// iterate through governors by operator address, execute func for each governor
	IterateGovernors(sdk.Context,
		func(index int64, governor GovernorI) (stop bool))

	// iterate through bonded governors by operator address, execute func for each governor
	IterateBondedGovernorsByPower(sdk.Context,
		func(index int64, governor GovernorI) (stop bool))

	Governor(sdk.Context, sdk.ValAddress) GovernorI // get a particular governor by operator address
	TotalBondedTokens(sdk.Context) math.Int         // total bonded tokens within the governor set
	StakingTokenSupply(sdk.Context) math.Int        // total staking token supply

	// Delegation allows for getting a particular delegation for a given governor
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI

	// MaxGovernors returns the maximum amount of bonded governors
	MaxGovernors(sdk.Context) uint32
}

type StakingComptability interface {
	// iterate through validators by operator address, execute func for each validator
	IterateValidators(sdk.Context,
		func(index int64, validator stakingtypes.ValidatorI) (stop bool))

	Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI            // get a particular validator by operator address
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI // get a particular validator by consensus address

	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI

	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))

	GetAllSDKDelegations(ctx sdk.Context) []stakingtypes.Delegation
}

// DelegationSet expected properties for the set of all delegations for a particular (noalias)
type DelegationSet interface {
	// iterate through all delegations from one delegator by governor-AccAddress,
	//   execute func for each governor
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))
}

// GovernorI expected validator functions
type GovernorI interface {
	GetMoniker() string          // moniker of the validator
	GetStatus() BondStatus       // status of the validator
	IsBonded() bool              // check if has a bonded status
	IsUnbonded() bool            // check if has status unbonded
	IsUnbonding() bool           // check if has status unbonding
	GetOperator() sdk.ValAddress // operator address to receive/return validators coins
	GetTokens() math.Int         // validation tokens
	GetBondedTokens() math.Int   // validator bonded tokens
	// todo: remove
	GetConsensusPower(math.Int) int64                        // validation power in tendermint
	GetCommission() sdk.Dec                                  // validator commission rate
	GetMinSelfDelegation() math.Int                          // validator minimum self delegation
	GetDelegatorShares() sdk.Dec                             // total outstanding delegator shares
	TokensFromShares(sdk.Dec) sdk.Dec                        // token worth of provided delegator shares
	TokensFromSharesTruncated(sdk.Dec) sdk.Dec               // token worth of provided delegator shares, truncated
	TokensFromSharesRoundUp(sdk.Dec) sdk.Dec                 // token worth of provided delegator shares, rounded up
	SharesFromTokens(amt math.Int) (sdk.Dec, error)          // shares worth of delegator's bond
	SharesFromTokensTruncated(amt math.Int) (sdk.Dec, error) // truncated shares worth of delegator's bond
}
