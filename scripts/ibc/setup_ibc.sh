#!/bin/bash

BASEDIR=$(dirname "$0")
. "$BASEDIR"/../shared.sh

IBC_PORT=transfer
IBC_VERSION=ics20-1


# --------------------------------- rly init --------------------------------- #
RLY_PATH="$HOME/.relayer"
RLY_CONFIG_FILE="$RLY_PATH/config/config.yaml"
ROLLAPP_IBC_CONF_FILE="$BASEDIR/rollapp.json"
HUB_IBC_CONF_FILE="$BASEDIR/hub.json"

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
rly config init

echo '# ------------------------- adding chains to rly config ------------------------- #'
tmp=$(mktemp)

jq --arg key "$RELAYER_KEY_FOR_ROLLAP" '.value.key = $key' $ROLLAPP_IBC_CONF_FILE > "$tmp" && mv "$tmp" $ROLLAPP_IBC_CONF_FILE
jq --arg chain "$CHAIN_ID" '.value."chain-id" = $chain' $ROLLAPP_IBC_CONF_FILE > "$tmp" && mv "$tmp" $ROLLAPP_IBC_CONF_FILE
jq --arg rpc "$ROLLAPP_RPC_FOR_RELAYER" '.value."rpc-addr" = $rpc' $ROLLAPP_IBC_CONF_FILE > "$tmp" && mv "$tmp" $ROLLAPP_IBC_CONF_FILE
jq --arg denom "0.0$DENOM" '.value."gas-prices" = $denom' $ROLLAPP_IBC_CONF_FILE > "$tmp" && mv "$tmp" $ROLLAPP_IBC_CONF_FILE
if [ "$EVM_ENABLED" ]; then
  jq '.value."account-prefix" = "ethm"' $ROLLAPP_IBC_CONF_FILE > "$tmp" && mv "$tmp" $ROLLAPP_IBC_CONF_FILE
fi

jq --arg key "$RELAYER_KEY_FOR_HUB" '.value.key = $key' $HUB_IBC_CONF_FILE > "$tmp" && mv "$tmp" $HUB_IBC_CONF_FILE
jq --arg chain "$SETTLEMENT_CHAIN_ID" '.value."chain-id" = $chain' $HUB_IBC_CONF_FILE > "$tmp" && mv "$tmp" $HUB_IBC_CONF_FILE
jq --arg rpc "$SETTLEMENT_RPC_FOR_RELAYER" '.value."rpc-addr" = $rpc' $HUB_IBC_CONF_FILE > "$tmp" && mv "$tmp" $HUB_IBC_CONF_FILE

rly chains add --file "$ROLLAPP_IBC_CONF_FILE" "$CHAIN_ID"
rly chains add --file "$HUB_IBC_CONF_FILE" "$SETTLEMENT_CHAIN_ID"

echo '# -------------------------------- creating keys ------------------------------- #'
rly keys add "$CHAIN_ID" "$RELAYER_KEY_FOR_ROLLAP"
rly keys add "$SETTLEMENT_CHAIN_ID" "$RELAYER_KEY_FOR_HUB"

RLY_HUB_ADDR=$(rly keys show "$SETTLEMENT_CHAIN_ID")
RLY_ROLLAPP_ADDR=$(rly keys show "$CHAIN_ID")

echo "# ------------------------------- balance of rly account on hub [$RLY_HUB_ADDR]------------------------------ #"
$SETTLEMENT_EXECUTABLE q bank balances "$(rly keys show "$SETTLEMENT_CHAIN_ID")" --node "$SETTLEMENT_RPC_FOR_RELAYER"
echo "From within the hub node: \"$SETTLEMENT_EXECUTABLE tx bank send $KEY_NAME_GENESIS $RLY_HUB_ADDR 100000000udym --keyring-backend test --broadcast-mode block\""

echo "# ------------------------------- balance of rly account on rollapp [$RLY_ROLLAPP_ADDR] ------------------------------ #"
$EXECUTABLE q bank balances "$(rly keys show "$CHAIN_ID")" --node "$ROLLAPP_RPC_FOR_RELAYER"
echo "From within the rollapp node: \"$EXECUTABLE tx bank send $KEY_NAME_ROLLAPP $RLY_ROLLAPP_ADDR 100000000$DENOM --keyring-backend test --broadcast-mode block\""

echo "waiting to fund accounts. Press to continue..."
read -r answer

echo '# -------------------------------- creating IBC link ------------------------------- #'
rly chains set-settlement "$SETTLEMENT_CHAIN_ID"
rly paths new "$CHAIN_ID" "$SETTLEMENT_CHAIN_ID" "$RELAYER_PATH" --src-port "$IBC_PORT" --dst-port "$IBC_PORT" --version "$IBC_VERSION"
rly transact link -t300s "$RELAYER_PATH" --src-port "$IBC_PORT" --dst-port "$IBC_PORT" --version "$IBC_VERSION"