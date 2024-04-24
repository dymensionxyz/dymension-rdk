package types

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegationI delegation bond for a delegated proof of stake system
type DelegationI interface {
	GetDelegatorAddr() sdk.AccAddress // delegator sdk.AccAddress for the bond
	GetGovernorAddr() sdk.ValAddress  // validator operator address
	GetShares() sdk.Dec               // amount of validator's shares held in this delegation
}

// GovernorI expected validator functions
type GovernorI interface {
	GetMoniker() string                                      // moniker of the validator
	GetStatus() BondStatus                                   // status of the validator
	IsBonded() bool                                          // check if has a bonded status
	IsUnbonded() bool                                        // check if has status unbonded
	IsUnbonding() bool                                       // check if has status unbonding
	GetOperator() sdk.ValAddress                             // operator address to receive/return validators coins
	GetTokens() math.Int                                     // validation tokens
	GetBondedTokens() math.Int                               // validator bonded tokens
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
