BASEDIR=$(dirname "$0")
echo "$BASEDIR"
source "$BASEDIR"/shared.sh


#Register rollapp 
$SETTLEMENT_EXECUTABLE tx rollapp create-rollapp "$ROLLAPP_ID" stamp1 "genesis-path/1" 3 100 '{"Addresses":[]}' \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test \
  --broadcast-mode block