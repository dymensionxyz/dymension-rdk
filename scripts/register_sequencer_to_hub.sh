BASEDIR=$(dirname "$0")
source "$BASEDIR"/shared.sh

#Register Sequencer
DESCRIPTION="{\"Moniker\":\"$MONIKER_NAME\",\"Identity\":\"\",\"Website\":\"\",\"SecurityContact\":\"\",\"Details\":\"\"}";
SEQ_PUB_KEY="$($EXECUTABLE dymint show-sequencer --home $CHAIN_DIR)"
SEQ_ACCOUNT_ON_HUB="$($SETTLEMENT_EXECUTABLE keys show -a $KEY_NAME_DYM --home $CHAIN_DIR --keyring-dir $KEYRING_PATH --keyring-backend test)"


#TODO: this should check the address provided, not through keyring!
echo "Current balance of sequencer account on hub: "
$SETTLEMENT_EXECUTABLE q bank balances "$SEQ_ACCOUNT_ON_HUB" --node tcp://"$SETTLEMENT_RPC"


read -r -p "Transfer funds if needed and continue..."

$SETTLEMENT_EXECUTABLE tx sequencer create-sequencer "$SEQ_PUB_KEY" "$ROLLAPP_ID" "$DESCRIPTION" \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test \
  --broadcast-mode block \
  --node tcp://"$SETTLEMENT_RPC"