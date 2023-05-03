#!/bin/bash
BASEDIR=$(dirname "$0")
. "$BASEDIR"/../shared.sh


#TODO: make common function
SEQ_ACCOUNT_ON_HUB=$(getSeqAddrOnHub)
echo "Current balance of sequencer account on hub[$SEQ_ACCOUNT_ON_HUB]: "
$SETTLEMENT_EXECUTABLE q bank balances "$SEQ_ACCOUNT_ON_HUB" --node "$SETTLEMENT_RPC"

echo "Transfer funds if needed and continue..."
read -r answer

#Register Sequencer
DESCRIPTION="{\"Moniker\":\"$MONIKER\",\"Identity\":\"\",\"Website\":\"\",\"SecurityContact\":\"\",\"Details\":\"\"}";
SEQ_PUB_KEY="$($EXECUTABLE dymint show-sequencer --home $ROLLAPP_CHAIN_DIR)"

$SETTLEMENT_EXECUTABLE tx sequencer create-sequencer "$SEQ_PUB_KEY" "$ROLLAPP_ID" "$DESCRIPTION" \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test \
  --broadcast-mode block \
  --node "$SETTLEMENT_RPC"