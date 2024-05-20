package types

// Event types for the gasless module.
const (
	EventTypeCreateGasTank       = "create_gas_tank"
	EventTypeAuthorizeActors     = "authorize_actors"
	EventTypeUpdateGasTankStatus = "update_gas_tank_status"
	EventTypeUpdateGasTankConfig = "update_gas_tank_config"
	EventTypeBlockConsumer       = "block_consumer"
	EventTypeUnblockConsumer     = "unblock_consumer"
	EventTypeFeeConsumption      = "fee_consumption"

	AttributeKeyCreator                = "creator"
	AttributeKeyProvider               = "provider"
	AttributeKeyActor                  = "actor"
	AttributeKeyConsumer               = "consumer"
	AttributeKeyGasTankID              = "gas_tank_id"
	AttributeKeyFeeDenom               = "fee_denom"
	AttributeKeyAuthorizedActors       = "authorized_actors"
	AttributeKeyGasTankStatus          = "gas_tank_status"
	AttributeKeyMaxFeeUsagePerTx       = "max_fee_usage_per_tx"
	AttributeKeyMaxTxsCountPerConsumer = "max_txs_count_per_consumer"
	AttributeKeyMaxFeeUsagePerConsumer = "max_fee_usage_per_consumer"
	AttributeKeyTxsAllowed             = "txs_allowed"
	AttributeKeyContractsAllowed       = "contracts_allowed"
	AttributeKeyFeeConsumptionMessage  = "message"
	AttributeKeyFeeSource              = "fee_source"
	AttributeKeyFailedGasTankIDs       = "failed_gas_tank_ids"
	AttributeKeyFailedGasTankErrors    = "failed_gas_tank_errors"
	AttributeKeySucceededGtid          = "succeeded_gas_tank_id"
)
