package types

// staking module event types
const (
	EventTypeCompleteUnbonding         = "complete_unbonding"
	EventTypeCompleteRedelegation      = "complete_redelegation"
	EventTypeCreateGovernor            = "create_Governor"
	EventTypeEditGovernor              = "edit_Governor"
	EventTypeDelegate                  = "delegate"
	EventTypeUnbond                    = "unbond"
	EventTypeCancelUnbondingDelegation = "cancel_unbonding_delegation"
	EventTypeRedelegate                = "redelegate"

	AttributeKeyGovernor          = "Governor"
	AttributeKeyCommissionRate    = "commission_rate"
	AttributeKeyMinSelfDelegation = "min_self_delegation"
	AttributeKeySrcGovernor       = "source_Governor"
	AttributeKeyDstGovernor       = "destination_Governor"
	AttributeKeyDelegator         = "delegator"
	AttributeKeyCreationHeight    = "creation_height"
	AttributeKeyCompletionTime    = "completion_time"
	AttributeKeyNewShares         = "new_shares"
	AttributeValueCategory        = ModuleName
)
