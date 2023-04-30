#!/bin/bash
BASEDIR=$(dirname "$0")
. "$BASEDIR"/../shared.sh


MAX_SEQUENCERS=5

#TODO: make common function
SEQ_ACCOUNT_ON_HUB="$($SETTLEMENT_EXECUTABLE keys show -a $KEY_NAME_DYM --keyring-dir $KEYRING_PATH --keyring-backend test)"
echo "Current balance of sequencer account on hub[$SEQ_ACCOUNT_ON_HUB]: "
$SETTLEMENT_EXECUTABLE q bank balances "$SEQ_ACCOUNT_ON_HUB" --node "$SETTLEMENT_RPC"

echo "Transfer funds if needed and continue..."
read -r answer
#Register rollapp 
$SETTLEMENT_EXECUTABLE tx rollapp create-rollapp "$ROLLAPP_ID" stamp1 "genesis-path/1" 3 "$MAX_SEQUENCERS" '{"Addresses":[]}' \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test \
  --keyring-dir "$KEYRING_PATH" \
  --broadcast-mode block \
  --node "$SETTLEMENT_RPC"
