#!/bin/bash
BASEDIR=$(dirname "$0")
. "$BASEDIR"/../shared.sh

MAX_SEQUENCERS=5

#Register rollapp 
$SETTLEMENT_EXECUTABLE tx rollapp create-rollapp "$ROLLAPP_CHAIN_ID" stamp1 "genesis-path/1" 3 "$MAX_SEQUENCERS" '{"Addresses":[]}' \
  --from "$KEY_NAME_GENESIS" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test \
  --broadcast-mode block \
  --node "$SETTLEMENT_RPC"
