BASEDIR=$(dirname "$0")
source "$BASEDIR"/shared.sh

IBC_PORT=transfer
IBC_VERSION=ics20-1


# --------------------------------- rly init --------------------------------- #
RLY_PATH="$HOME/.relayer"
RLY_CONFIG_FILE="$RLY_PATH/config/config.yaml"

if [ -f "$RLY_CONFIG_FILE" ]; then
  printf "======================================================================================================\n"
  echo "A rly config file already exists. Overwrite? (y/N)"
  printf "======================================================================================================\n"
  read -r answer
  if [[ "$answer" == "Y" || "$answer" == "y" ]]; then
    rm -rf "$RLY_PATH"
  fi
fi

echo '# -------------------------- initializing rly config ------------------------- #'
rly config init --settlement-config "$SETTLEMENT_CONFIG"

echo '# ------------------------- adding chains to rly config ------------------------- #'
tmp=$(mktemp)
ROLLAPP_IBC_CONF_FILE="$BASEDIR/ibc/rollapp.json"
jq --arg key "$RELAYER_KEY_FOR_ROLLAP" '.value.key = $key' $ROLLAPP_IBC_CONF_FILE > "$tmp" && mv "$tmp" $ROLLAPP_IBC_CONF_FILE
jq --arg chain "$CHAIN_ID" '.value."chain-id" = $chain' $ROLLAPP_IBC_CONF_FILE > "$tmp" && mv "$tmp" $ROLLAPP_IBC_CONF_FILE
jq --arg rpc "tcp://$ROLLAPP_RPC_FOR_RELAYER" '.value."rpc-addr" = $rpc' $ROLLAPP_IBC_CONF_FILE > "$tmp" && mv "$tmp" $ROLLAPP_IBC_CONF_FILE

HUB_IBC_CONF_FILE="$BASEDIR/ibc/hub.json"
jq --arg key "$RELAYER_KEY_FOR_HUB" '.value.key = $key' $HUB_IBC_CONF_FILE > "$tmp" && mv "$tmp" $HUB_IBC_CONF_FILE
jq --arg chain "$SETTLEMENT_CHAIN_ID" '.value."chain-id" = $chain' $HUB_IBC_CONF_FILE > "$tmp" && mv "$tmp" $HUB_IBC_CONF_FILE
jq --arg rpc "tcp://$SETTLEMENT_RPC_FOR_RELAYER" '.value."rpc-addr" = $rpc' $HUB_IBC_CONF_FILE > "$tmp" && mv "$tmp" $HUB_IBC_CONF_FILE

rly chains add --file "$BASEDIR/ibc/rollapp.json" "$CHAIN_ID"
rly chains add --file "$BASEDIR/ibc/hub.json" "$SETTLEMENT_CHAIN_ID"

echo '# -------------------------------- creating keys ------------------------------- #'
rly keys add "$CHAIN_ID" "$RELAYER_KEY_FOR_ROLLAP"
rly keys add "$SETTLEMENT_CHAIN_ID" "$RELAYER_KEY_FOR_HUB"

echo "# ------------------------------- balance of rly account on hub [$(rly keys show "$SETTLEMENT_CHAIN_ID")]------------------------------ #"
$SETTLEMENT_EXECUTABLE q bank balances "$(rly keys show "$SETTLEMENT_CHAIN_ID")" --node tcp://"$SETTLEMENT_RPC_FOR_RELAYER"

echo "# ------------------------------- balance of rly account on rollapp [$(rly keys show "$CHAIN_ID")] ------------------------------ #"
$EXECUTABLE q bank balances "$(rly keys show "$CHAIN_ID")" --node tcp://"$ROLLAPP_RPC_FOR_RELAYER"

read -r -p "waiting to fund accounts. Press to continue..."

echo '# -------------------------------- creating IBC link ------------------------------- #'
rly paths new "$CHAIN_ID" "$SETTLEMENT_CHAIN_ID" "$RELAYER_PATH" --src-port "$IBC_PORT" --dst-port "$IBC_PORT" --version "$IBC_VERSION"
rly transact link "$RELAYER_PATH" --src-port "$IBC_PORT" --dst-port "$IBC_PORT" --version "$IBC_VERSION"