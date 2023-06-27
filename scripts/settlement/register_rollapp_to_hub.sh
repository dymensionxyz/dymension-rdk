#!/bin/bash
BASEDIR=$(dirname "$0")
. "$BASEDIR"/../shared.sh

MAX_SEQUENCERS=5

#Register rollapp 
$SETTLEMENT_EXECUTABLE tx rollapp create-rollapp "$ROLLAPP_CHAIN_ID" stamp1 "genesis-path/1" 3 "$MAX_SEQUENCERS" '{"Addresses":[]}' \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test \
  --keyring-dir "$KEYRING_PATH" \
  --broadcast-mode block \
  --node "$SETTLEMENT_RPC"
