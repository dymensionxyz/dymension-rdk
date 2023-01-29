#!/bin/bash
BASEDIR=$(dirname "$0")
. "$BASEDIR"/shared.sh

# If the settlement layer is set to dymension, use $SETTLEMENT_CONFIG otherwise use $SETTLEMENT_CONFIG_MOCK
if [ ! "$SETTLEMENT_LAYER" = "dymension" ]; then
  echo "using mock settlement layer"
  SETTLEMENT_CONFIG="$SETTLEMENT_CONFIG_MOCK"
fi

if [ "$DA_LAYER" = "mock" ]; then
  echo "using mock DA layer"
  DA_LAYER_CONFIG="30s"
fi

# If aggregator is set to true pass the aggregator flag
if [ "$AGGREGATOR" = "true" ]; then
  AGGREGATOR_FLAG="--dymint.aggregator"
else
  AGGREGATOR_FLAG=""
fi

$EXECUTABLE start $AGGREGATOR_FLAG \
  --dymint.da_layer "$DA_LAYER" \
  --dymint.da_config "$DA_LAYER_CONFIG" \
  --dymint.settlement_layer "$SETTLEMENT_LAYER" \
  --dymint.settlement_config "$SETTLEMENT_CONFIG" \
  --dymint.block_batch_size "$BATCH_SIZE" \
  --dymint.namespace_id "$NAMESPACE_ID" \
  --dymint.block_time "$BLOCK_TIME" \
  --p2p.seeds "$ROLLAPP_SEEDS" \
  --home "$CHAIN_DIR" \
  --log_level "$LOG_LEVEL"

