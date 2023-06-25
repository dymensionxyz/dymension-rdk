# ---------------------------------------------------------------------------- #
#                               SETTLEMENT CONFIG                              #
# ---------------------------------------------------------------------------- #
SETTLEMENT_EXECUTABLE=${SETTLEMENT_EXECUTABLE:-dymd}
KEY_NAME_GENESIS=${KEY_NAME_GENESIS:-"local-user"}
SETTLEMENT_CHAIN_ID=${SETTLEMENT_CHAIN_ID:-"dymension_100-1"}
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
# ROLLAPP_CHAIN_DIR="$HOME/.rollapp_evm"
# ROLLAPP_CHAIN_ID="rollappevm_100-1"

# ---------------------------------- GLOBAL ---------------------------------- #
KEY_NAME_ROLLAPP=${KEY_NAME_ROLLAPP:-"rol-user"}
DENOM=${DENOM:-"urax"}

EXECUTABLE=${EXECUTABLE:-rollappd}
ROLLAPP_CHAIN_DIR=${ROLLAPP_CHAIN_DIR:-$HOME/.rollapp}
ROLLAPP_CHAIN_ID=${ROLLAPP_CHAIN_ID:-rollapp}
MONIKER=${MONIKER:-$ROLLAPP_CHAIN_ID-sequencer}

RPC_LADDRESS=${RPC_LADDRESS:-"0.0.0.0:26657"}
P2P_LADDRESS=${P2P_LADDRESS:-"0.0.0.0:26656"}
GRPC_LADDRESS=${GRPC_LADDRESS:-"0.0.0.0:8080"}
GRPC_WEB_LADDRESS=${GRPC_WEB_LADDRESS:-"0.0.0.0:8081"}
API_ADDRESS=${API_ADDRESS:-"0.0.0.0:1417"}
UNSAFE_CORS=${UNSAFE_CORS:-"true"}

LOG_LEVEL=${LOG_LEVEL:-"info"}
# LOG_FILE_PATH=${LOG_FILE_PATH:-"$ROLLAPP_CHAIN_DIR/log/rollapp.log"}
MAX_LOG_SIZE=${MAX_LOG_SIZE:-"2000"}
MODULE_LOG_LEVEL_OVERRIDE=${MODULE_LOG_LEVEL_OVERRIDE:-""}

# ------------------------------- dymint config ------------------------------ #
KEY_NAME_DYM=${KEY_NAME_DYM:-"sequencer"}

#Its the keyring that will be used by dymint sequencer
KEYRING_PATH=${KEYRING_PATH:-$ROLLAPP_CHAIN_DIR/sequencer_keys}

# ---------------------------------------------------------------------------- #
#                                  IBC CONFIG                                  #
# ---------------------------------------------------------------------------- #
RELAYER_KEY_FOR_ROLLAP=${RELAYER_KEY_FOR_ROLLAP:-"relayer-rollapp-key"}
RELAYER_KEY_FOR_HUB=${RELAYER_KEY_FOR_HUB:-"relayer-hub-key"}
RELAYER_PATH=${RELAYER_PATH:-"hub-rollapp"}
ROLLAPP_RPC_FOR_RELAYER=${ROLLAPP_RPC_FOR_RELAYER:-"tcp://127.0.0.1:26657"}
SETTLEMENT_RPC_FOR_RELAYER=${SETTLEMENT_RPC_FOR_RELAYER:-$SETTLEMENT_RPC}
ROLLAPP_CHANNEL_NAME=${ROLLAPP_CHANNEL_NAME:-"channel-0"}
HUB_CHANNEL_NAME=${HUB_CHANNEL_NAME:-"channel-0"}
