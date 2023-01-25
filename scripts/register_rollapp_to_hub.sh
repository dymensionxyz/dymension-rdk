BASEDIR=$(dirname "$0")
source "$BASEDIR"/shared.sh


#Register rollapp 
$SETTLEMENT_EXECUTABLE tx rollapp create-rollapp "$ROLLAPP_ID" stamp1 "genesis-path/1" 3 1 '{"Addresses":[]}' \
  --from "$KEY_NAME_GENESIS" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test \
  --broadcast-mode block \
  --node tcp://"$SETTLEMENT_RPC"