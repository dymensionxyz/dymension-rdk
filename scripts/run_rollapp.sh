BASEDIR=$(dirname "$0")
echo "$BASEDIR"
source "$BASEDIR"/shared.sh


#TODO: print run configuration

SETTLEMENT_CONFIG="{\"node_address\": \"http://$SETTLEMENT_RPC\", \"rollapp_id\": \"$ROLLAPP_ID\", \"dym_account_name\": \"$KEY_NAME_DYM\", \"keyring_home_dir\": \"$KEYRING_PATH\", \"keyring_backend\":\"test\"}"
SETTLEMENT_CONFIG_MOCK="{\"root_dir\": \""$HOME"/.rollapp\", \"db_path\": \"data\"}" \


#TODO: make settlement mock a parameter
$EXECUTABLE start --dymint.aggregator \
  --dymint.da_layer mock \
  --dymint.settlement_layer dymension \
  --dymint.settlement_config "$SETTLEMENT_CONFIG" \
  --dymint.block_batch_size 20 \
  --dymint.namespace_id "$NAMESPACE_ID" \
  --dymint.block_time 0.5s \
  --rpc.laddr "tcp://0.0.0.0:26667" \
  --p2p.laddr "tcp://0.0.0.0:26666" \
  --grpc.address "0.0.0.0:9080" \
  --grpc-web.address "0.0.0.0:9081"

