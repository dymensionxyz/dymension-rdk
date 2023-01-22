BASEDIR=$(dirname "$0")
source "$BASEDIR"/shared.sh

IBC_PORT=transfer
IBC_VERSION=ics20-1


# ---------------------------------------------------------------------------- #
#                                     utils                                    #
# ---------------------------------------------------------------------------- #
fund_hub_account() {
    FROM=${1:-$KEY_NAME_GENESIS}
    TO=${2:-$(rly keys show $SETTLEMENT_CHAIN_ID)}
    AMOUNT=${3:-100000000udym}
    echo "funding $TO with $AMOUNT from $FROM"
    
    $SETTLEMENT_EXECUTABLE tx bank send "$FROM" "$TO" "$AMOUNT" \
        --chain-id "$SETTLEMENT_CHAIN_ID" \
        --keyring-backend test \
        --broadcast-mode block
}

fund_rollapp_account() {
    FROM=${1:-$KEY_NAME_ROLLAPP}
    TO=${2:-$(rly keys show $CHAIN_ID)}
    AMOUNT=${3:-100000000urap}
    echo "funding $TO with $AMOUNT from $FROM"
    
    $EXECUTABLE tx bank send "$FROM" "$TO" "$AMOUNT" \
        --chain-id "$CHAIN_ID" \
        --keyring-backend test \
        --home $CHAIN_DIR \
        --broadcast-mode block
}


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
SETTLEMENT_CONFIG="{\"node_address\": \"http://$SETTLEMENT_RPC\", \"rollapp_id\": \"$ROLLAPP_ID\", \"dym_account_name\": \"$KEY_NAME_DYM\", \"keyring_home_dir\": \"$KEYRING_PATH\", \"keyring_backend\":\"test\"}"
rly config init --settlement-config "$SETTLEMENT_CONFIG"

echo '# ------------------------- adding chains to rly config ------------------------- #'
tmp=$(mktemp)
ROLLAPP_IBC_CONF_FILE="$BASEDIR/ibc/rollapp.json"
jq --arg key "$RELAYER_KEY_FOR_ROLLAP" '.value.key = $key' $ROLLAPP_IBC_CONF_FILE > "$tmp" && mv "$tmp" $ROLLAPP_IBC_CONF_FILE
jq --arg chain "$CHAIN_ID" '.value."chain-id" = $chain' $ROLLAPP_IBC_CONF_FILE > "$tmp" && mv "$tmp" $ROLLAPP_IBC_CONF_FILE
jq --arg rpc "tcp://$RPC_PORT" '.value."rpc-addr" = $rpc' $ROLLAPP_IBC_CONF_FILE > "$tmp" && mv "$tmp" $ROLLAPP_IBC_CONF_FILE

HUB_IBC_CONF_FILE="$BASEDIR/ibc/hub.json"
jq --arg key "$RELAYER_KEY_FOR_HUB" '.value.key = $key' $HUB_IBC_CONF_FILE > "$tmp" && mv "$tmp" $HUB_IBC_CONF_FILE
jq --arg chain "$SETTLEMENT_CHAIN_ID" '.value."chain-id" = $chain' $HUB_IBC_CONF_FILE > "$tmp" && mv "$tmp" $HUB_IBC_CONF_FILE
jq --arg rpc "tcp://$SETTLEMENT_RPC" '.value."rpc-addr" = $rpc' $HUB_IBC_CONF_FILE > "$tmp" && mv "$tmp" $HUB_IBC_CONF_FILE

rly chains add --file "$BASEDIR/ibc/rollapp.json" "$CHAIN_ID"
rly chains add --file "$BASEDIR/ibc/hub.json" "$SETTLEMENT_CHAIN_ID"

echo '# -------------------------------- creating keys ------------------------------- #'
rly keys add "$CHAIN_ID" "$RELAYER_KEY_FOR_ROLLAP"
rly keys add "$SETTLEMENT_CHAIN_ID" "$RELAYER_KEY_FOR_HUB"

echo '# ------------------------------- fund rly account on hub ------------------------------ #'
#TODO: check if accounts exist and funded before funding again
fund_hub_account
echo '# ------------------------------- fund rly on rollapp ------------------------------ #'
fund_rollapp_account

#TODO: validate by code the accounts are actually funded
dymd q bank balances "$(rly keys show $SETTLEMENT_CHAIN_ID)"
rollappd q bank balances "$(rly keys show $CHAIN_ID)" --home $CHAIN_DIR

echo '# -------------------------------- creating IBC link ------------------------------- #'
rly paths new "$CHAIN_ID" "$SETTLEMENT_CHAIN_ID" "$RELAYER_PATH" --src-port "$IBC_PORT" --dst-port "$IBC_PORT" --version "$IBC_VERSION"
rly transact link "$RELAYER_PATH" --src-port "$IBC_PORT" --dst-port "$IBC_PORT" --version "$IBC_VERSION"