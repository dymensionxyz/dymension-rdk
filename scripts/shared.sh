# ---------------------------------------------------------------------------- #
#                               SETTLEMENT CONFIG                              #
# ---------------------------------------------------------------------------- #
SETTLEMENT_EXECUTABLE=${SETTLEMENT_EXECUTABLE:-dymd}
KEY_NAME_GENESIS=${KEY_NAME_GENESIS:-"local-user"}
SETTLEMENT_CHAIN_ID=${SETTLEMENT_CHAIN_ID:-"local-testnet"}
SETTLEMENT_RPC=${SETTLEMENT_RPC:-"tcp://127.0.0.1:36657"}

# ---------------------------------------------------------------------------- #
#                                ROLLAPP CONFIG                                #
# ---------------------------------------------------------------------------- #
# ----------------------------- UNCOMMENT FOR EVM ---------------------------- #
# EVM_ENABLED=true

# # Assuming 100,000,000REVM tokens
# # evm uses 10^18 decimal precision for arevm
# DENOM=arevm
# result=$(echo "100 * 10^6 * 10^18" | bc)
# staking_result=$(echo "$result / 2000" | bc)
# TOKEN_AMOUNT="$result""$DENOM"
# STAKING_AMOUNT="$staking_result""$DENOM"

# EXECUTABLE="rollapp_evm"
# CHAIN_DIR="$HOME/.rollapp_evm"
# CHAIN_ID="rollappevm_100_1"

# ---------------------------------- GLOBAL ---------------------------------- #
KEY_NAME_ROLLAPP=${KEY_NAME_ROLLAPP:-"rol-user"}
DENOM=${DENOM:-"urax"}

EXECUTABLE=${EXECUTABLE:-rollappd}
CHAIN_DIR=${CHAIN_DIR:-$HOME/.rollapp}
CHAIN_ID=${CHAIN_ID:-rollapp}
ROLLAPP_ID=${ROLLAPP_ID:-$CHAIN_ID}
MONIKER=${MONIKER:-$CHAIN_ID-sequencer}

RPC_LADDRESS=${RPC_LADDRESS:-"0.0.0.0:26657"}
P2P_LADDRESS=${P2P_LADDRESS:-"0.0.0.0:26656"}
GRPC_LADDRESS=${GRPC_LADDRESS:-"0.0.0.0:8080"}
GRPC_WEB_LADDRESS=${GRPC_WEB_LADDRESS:-"0.0.0.0:8081"}
API_ADDRESS=${API_ADDRESS:-"0.0.0.0:1317"}
UNSAFE_CORS=${UNSAFE_CORS:-""}

LOG_LEVEL=${LOG_LEVEL:-"info"}
# LOG_FILE_PATH=${LOG_FILE_PATH:-"$CHAIN_DIR/log/rollapp.log"}
MAX_LOG_SIZE=${MAX_LOG_SIZE:-"2000"}
MODULE_LOG_LEVEL_OVERRIDE=${MODULE_LOG_LEVEL_OVERRIDE:-""}

P2P_SEEDS=${P2P_SEEDS:-""}
ROLLAPP_PEERS=${ROLLAPP_PEERS:-""}


# ------------------------------- dymint config ------------------------------ #
KEY_NAME_DYM=${KEY_NAME_DYM:-"$MONIKER"}
KEYRING_PATH=${KEYRING_PATH:-"$HOME/.dymension"}

AGGREGATOR=${AGGREGATOR:-"true"}
BATCH_SIZE=${BATCH_SIZE:-"60"}
BLOCK_TIME=${BLOCK_TIME:-"0.2s"}
DYMINT_FEES=${DYMINT_FEES:-"2000000udym"}

# DA CONFIG
DA_LAYER=${DA_LAYER:-"mock"}
DA_LC_ENDPOINT=${DA_LC_ENDPOINT:-"127.0.0.1:26659"}
NAMESPACE_ID=${NAMESPACE_ID:-"000000000000ffff"}
DA_NAMESPACE_ID=${DA_NAMESPACE_ID:-"[0,0,0,0,0,0,255,255]"}
DA_LAYER_CONFIG=${DA_LAYER_CONFIG:-"{\"base_url\": \"http:\/\/$DA_LC_ENDPOINT\", \"timeout\": 60000000000, \"fee\":20000, \"gas_limit\": 20000000, \"namespace_id\":$DA_NAMESPACE_ID}"}

# Settlement config
SETTLEMENT_LAYER=${SETTLEMENT_LAYER:-"mock"}
SETTLEMENT_CONFIG=${SETTLEMENT_CONFIG:-"{\"node_address\": \"$SETTLEMENT_RPC\", \"rollapp_id\": \"$ROLLAPP_ID\", \"dym_account_name\": \"$KEY_NAME_DYM\", \"keyring_home_dir\": \"$KEYRING_PATH\", \"keyring_backend\":\"test\", \"gas_fees\": \"$DYMINT_FEES\"}"}
SETTLEMENT_CONFIG_MOCK=${SETTLEMENT_CONFIG_MOCK:-"{\"root_dir\": \"$CHAIN_DIR\", \"db_path\": \"data\", \"proposer_pub_key\":\"$CHAIN_DIR/config/priv_validator_key.json\"}"}


# ---------------------------------------------------------------------------- #
#                                  IBC CONFIG                                  #
# ---------------------------------------------------------------------------- #
RELAYER_FEES=${RELAYER_FEES:-$DYMINT_FEES}
RELAYER_KEY_FOR_ROLLAP=${RELAYER_KEY_FOR_ROLLAP:-"relayer-rollapp-key"}
RELAYER_KEY_FOR_HUB=${RELAYER_KEY_FOR_HUB:-"relayer-hub-key"}
RELAYER_PATH=${RELAYER_PATH:-"hub-rollapp"}
ROLLAPP_RPC_FOR_RELAYER=${ROLLAPP_RPC_FOR_RELAYER:-"tcp://$RPC_LADDRESS"}
SETTLEMENT_RPC_FOR_RELAYER=${SETTLEMENT_RPC_FOR_RELAYER:-$SETTLEMENT_RPC}
ROLLAPP_CHANNEL_NAME=${ROLLAPP_CHANNEL_NAME:-"channel-0"}
HUB_CHANNEL_NAME=${HUB_CHANNEL_NAME:-"channel-0"}

RELAYER_SETTLEMENT_CONFIG=${RELAYER_SETTLEMENT_CONFIG:-"{\"node_address\": \"$SETTLEMENT_RPC\", \"rollapp_id\": \"$ROLLAPP_ID\", \"dym_account_name\": \"$KEY_NAME_DYM\", \"keyring_home_dir\": \"$KEYRING_PATH\", \"keyring_backend\":\"test\", \"gas_fees\": \"$RELAYER_FEES\"}"}

