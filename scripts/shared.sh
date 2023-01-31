# ---------------------------------------------------------------------------- #
#                               SETTLEMENT CONFIG                              #
# ---------------------------------------------------------------------------- #
KEY_NAME_GENESIS=${KEY_NAME_GENESIS:-"local-user"}
SETTLEMENT_EXECUTABLE=${SETTLEMENT_EXECUTABLE:-dymd}
SETTLEMENT_CHAIN_ID=${SETTLEMENT_CHAIN_ID:-"local-testnet"}
SETTLEMENT_RPC=${SETTLEMENT_RPC:-"127.0.0.1:36657"}


# ---------------------------------------------------------------------------- #
#                                ROLLAPP CONFIG                                #
# ---------------------------------------------------------------------------- #
KEY_NAME_ROLLAPP=${KEY_NAME_ROLLAPP:-"rol-user"}

EXECUTABLE=${EXECUTABLE:-rollappd}
CHAIN_DIR=${CHAIN_DIR:-$HOME/.rollapp}
CHAIN_ID=${CHAIN_ID:-rollapp}
ROLLAPP_ID=${ROLLAPP_ID:-$CHAIN_ID}
MONIKER=${MONIKER:-$CHAIN_ID-sequencer}
NAMESPACE_ID=${NAMESPACE_ID:-"000000000000ffff"}

RPC_LADDRESS=${RPC_LADDRESS:-"0.0.0.0:26667"}
P2P_LADDRESS=${P2P_LADDRESS:-"0.0.0.0:26668"}
GRPC_LADDRESS=${GRPC_LADDRESS:-"0.0.0.0:9080"}
GRPC_WEB_LADDRESS=${GRPC_WEB_LADDRESS:-"0.0.0.0:9081"}
API_ADDRESS=${API_ADDRESS:-"0.0.0.0:1417"}
LOG_LEVEL=${LOG_LEVEL:-"info"}
P2P_SEEDS=${P2P_SEEDS:-""}
UNSAFE_CORS=${UNSAFE_CORS:-""}
ROLLAPP_PEERS=${ROLLAPP_PEERS:-""}


# ------------------------------- dymint config ------------------------------ #
#TODO: rename to sequencer key name
#TODO: make most params based on chain ID
KEY_NAME_DYM=${KEY_NAME_DYM:-"local-sequencer"}
KEYRING_PATH=${KEYRING_PATH:-"$HOME/.rollapp"}

AGGREGATOR=${AGGREGATOR:-"true"}
BATCH_SIZE=${BATCH_SIZE:-"60"}
BLOCK_TIME=${BLOCK_TIME:-"0.2s"}

# DA CONFIG
DA_LAYER=${DA_LAYER:-"celestia"}
DA_LC_ENDPOINT=${DA_LC_ENDPOINT:-"127.0.0.1:26659"}
DA_NAMESPACE_ID=${DA_NAMESPACE_ID:-"[0,0,0,0,0,0,255,255]"}
DA_LAYER_CONFIG=${DA_LAYER_CONFIG:-"{\"base_url\": \"http:\/\/$DA_LC_ENDPOINT\", \"timeout\": 60000000000, \"fee\":9000, \"gas_limit\": 20000000, \"namespace_id\":$DA_NAMESPACE_ID}"}

# Settlement config
SETTLEMENT_LAYER=${SETTLEMENT_LAYER:-"dymension"}
SETTLEMENT_CONFIG=${SETTLEMENT_CONFIG:-"{\"node_address\": \"http://$SETTLEMENT_RPC\", \"rollapp_id\": \"$ROLLAPP_ID\", \"dym_account_name\": \"$KEY_NAME_DYM\", \"keyring_home_dir\": \"$KEYRING_PATH\", \"keyring_backend\":\"test\"}"}
SETTLEMENT_CONFIG_MOCK=${SETTLEMENT_CONFIG_MOCK:-"{\"root_dir\": \"$HOME/.rollapp\", \"db_path\": \"data\", \"proposer_pub_key\":\"$HOME/.rollapp/config/priv_validator_key.json\"}"}


# ---------------------------------------------------------------------------- #
#                                  IBC CONFIG                                  #
# ---------------------------------------------------------------------------- #
RELAYER_KEY_FOR_ROLLAP=${RELAYER_KEY_FOR_ROLLAP:-"relayer-rollapp-key"}
RELAYER_KEY_FOR_HUB=${RELAYER_KEY_FOR_HUB:-"relayer-hub-key"}
RELAYER_PATH=${RELAYER_PATH:-"hub-rollapp"}
ROLLAPP_RPC_FOR_RELAYER=${ROLLAPP_RPC_FOR_RELAYER:-$RPC_LADDRESS}
SETTLEMENT_RPC_FOR_RELAYER=${SETTLEMENT_RPC_FOR_RELAYER:-$SETTLEMENT_RPC}
ROLLAPP_CHANNEL_NAME=${ROLLAPP_CHANNEL_NAME:-"channel-0"}
HUB_CHANNEL_NAME=${HUB_CHANNEL_NAME:-"channel-0"}
