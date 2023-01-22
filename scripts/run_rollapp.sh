BASEDIR=$(dirname "$0")
source "$BASEDIR"/shared.sh


# echo "SETTLEMENT_CONFIG: $SETTLEMENT_CONFIG"

SETTLEMENT_CONFIG="{\"node_address\": \"http://$SETTLEMENT_RPC\", \"rollapp_id\": \"$ROLLAPP_ID\", \"dym_account_name\": \"$KEY_NAME_DYM\", \"keyring_home_dir\": \"$KEYRING_PATH\", \"keyring_backend\":\"test\"}"

#TODO: make running a mock through a parameter
$EXECUTABLE start --dymint.aggregator \
  --dymint.da_layer mock \
  --dymint.settlement_layer dymension \
  --dymint.settlement_config "$SETTLEMENT_CONFIG" \
  --dymint.block_batch_size 60 \
  --dymint.namespace_id "$NAMESPACE_ID" \
  --dymint.block_time 0.5s \
  --home $CHAIN_DIR

