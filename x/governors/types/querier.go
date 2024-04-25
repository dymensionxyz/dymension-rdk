package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the staking Querier
const (
	QueryGovernors                     = "governors"
	QueryGovernor                      = "governor"
	QueryDelegatorDelegations          = "delegatorDelegations"
	QueryDelegatorUnbondingDelegations = "delegatorUnbondingDelegations"
	QueryRedelegations                 = "redelegations"
	QueryGovernorDelegations           = "governorDelegations"
	QueryGovernorRedelegations         = "governorRedelegations"
	QueryGovernorUnbondingDelegations  = "governorUnbondingDelegations"
	QueryDelegation                    = "delegation"
	QueryUnbondingDelegation           = "unbondingDelegation"
	QueryDelegatorGovernors            = "delegatorGovernors"
	QueryDelegatorGovernor             = "delegatorGovernor"
	QueryPool                          = "pool"
	QueryParameters                    = "parameters"
	QueryHistoricalInfo                = "historicalInfo"
)

// defines the params for the following queries:
// - 'custom/staking/delegatorDelegations'
// - 'custom/staking/delegatorUnbondingDelegations'
// - 'custom/staking/delegatorGovernors'
type QueryDelegatorParams struct {
	DelegatorAddr sdk.AccAddress
}

func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/governor'
// - 'custom/staking/governorDelegations'
// - 'custom/staking/governorUnbondingDelegations'
type QueryGovernorParams struct {
	GovernorAddr sdk.ValAddress
	Page, Limit  int
}

func NewQueryGovernorParams(governorAddr sdk.ValAddress, page, limit int) QueryGovernorParams {
	return QueryGovernorParams{
		GovernorAddr: governorAddr,
		Page:         page,
		Limit:        limit,
	}
}

// defines the params for the following queries:
// - 'custom/staking/redelegation'
type QueryRedelegationParams struct {
	DelegatorAddr   sdk.AccAddress
	SrcGovernorAddr sdk.ValAddress
	DstGovernorAddr sdk.ValAddress
}

func NewQueryRedelegationParams(delegatorAddr sdk.AccAddress, srcGovernorAddr, dstGovernorAddr sdk.ValAddress) QueryRedelegationParams {
	return QueryRedelegationParams{
		DelegatorAddr:   delegatorAddr,
		SrcGovernorAddr: srcGovernorAddr,
		DstGovernorAddr: dstGovernorAddr,
	}
}

// QueryGovernorsParams defines the params for the following queries:
// - 'custom/staking/governors'
type QueryGovernorsParams struct {
	Page, Limit int
	Status      string
}

func NewQueryGovernorsParams(page, limit int, status string) QueryGovernorsParams {
	return QueryGovernorsParams{page, limit, status}
}
