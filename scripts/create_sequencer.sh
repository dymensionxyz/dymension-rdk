#!/bin/bash

BASEDIR=$(dirname "$0")
. "$BASEDIR"/shared.sh

$EXECUTABLE tx sequencers create-sequencer \
  --pubkey $($EXECUTABLE dymint show-sequencer --home $CHAIN_DIR) \
  --broadcast-mode block \
  --moniker $MONIKER \
  --chain-id $CHAIN_ID \
  --from $($EXECUTABLE keys show -a $KEY_NAME_ROLLAPP --keyring-backend test --home $CHAIN_DIR) \
  --keyring-backend test \
  --home $CHAIN_DIR