package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(authtypes.AccountI) (stop bool))
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI // only used for simulation

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	GetSupply(ctx sdk.Context, denom string) sdk.Coin

	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}

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
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) DelegationI

	// MaxGovernors returns the maximum amount of bonded governors
	MaxGovernors(sdk.Context) uint32

	/* ----------------------------- removed methods ---------------------------- */
	/*
		// slash the validator and delegators of the validator, specifying offence height, offence power, and slash fraction
		Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec) math.Int
		Jail(sdk.Context, sdk.ConsAddress)   // jail a validator
		Unjail(sdk.Context, sdk.ConsAddress) // unjail a validator
		ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) ValidatorI // get a particular validator by consensus address
		// iterate through the consensus validator set of the last block by operator address, execute func for each validator
		IterateLastValidators(sdk.Context,
			func(index int64, validator ValidatorI) (stop bool))
	*/

}

// DelegationSet expected properties for the set of all delegations for a particular (noalias)
type DelegationSet interface {
	// iterate through all delegations from one delegator by governor-AccAddress,
	//   execute func for each governor
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation DelegationI) (stop bool))
}

// Event Hooks
// These can be utilized to communicate between a staking keeper and another
// keeper which must take particular actions when governors/delegators change
// state. The second keeper must implement this interface, which then the
// staking keeper can call.

// StakingHooks event hooks for staking governor object (noalias)
type StakingHooks interface {
	AfterGovernorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error   // Must be called when a governor is created
	BeforeGovernorModified(ctx sdk.Context, valAddr sdk.ValAddress) error // Must be called when a governor's state changes
	AfterGovernorRemoved(ctx sdk.Context, valAddr sdk.ValAddress) error   // Must be called when a governor is deleted

	AfterGovernorBonded(ctx sdk.Context, valAddr sdk.ValAddress) error         // Must be called when a governor is bonded
	AfterGovernorBeginUnbonding(ctx sdk.Context, valAddr sdk.ValAddress) error // Must be called when a governor begins unbonding

	BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error        // Must be called when a delegation is created
	BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error // Must be called when a delegation's shares are modified
	BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error        // Must be called when a delegation is removed
	AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error
}
