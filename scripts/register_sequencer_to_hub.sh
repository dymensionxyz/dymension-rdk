BASEDIR=$(dirname "$0")
source "$BASEDIR"/shared.sh

#Register Sequencer
DESCRIPTION="{\"Moniker\":\"$MONIKER\",\"Identity\":\"\",\"Website\":\"\",\"SecurityContact\":\"\",\"Details\":\"\"}";
SEQ_PUB_KEY="$($EXECUTABLE dymint show-sequencer --home $CHAIN_DIR)"

$SETTLEMENT_EXECUTABLE tx sequencer create-sequencer "$SEQ_PUB_KEY" "$ROLLAPP_ID" "$DESCRIPTION" \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test \
  --keyring-dir "$KEYRING_PATH" \
  --broadcast-mode block \
  --node tcp://"$SETTLEMENT_RPC"