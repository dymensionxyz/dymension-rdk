#!/bin/bash
BASEDIR=$(dirname "$0")

if [ "$EVM_ENABLED" ]; then
  echo "EVM build enabled"
fi

. "$BASEDIR"/shared.sh
. "$BASEDIR"/set_genesis_config.sh


# ---------------------------- initial parameters ---------------------------- #
# Assuming 1,000,000RAP tokens
TOKEN_AMOUNT=${TOKEN_AMOUNT:-100000000000000000000000000urax}
#half is staked
STAKING_AMOUNT=${STAKING_AMOUNT:-500000000000urax}
SEQUENCER_AMOUNT=${SEQUENCER_AMOUNT:-10000000000udym}

CONFIG_DIRECTORY="$CHAIN_DIR/config"
GENESIS_FILE="$CONFIG_DIRECTORY/genesis.json"
TENDERMINT_CONFIG_FILE="$CONFIG_DIRECTORY/config.toml"
CLIENT_CONFIG_FILE="$CONFIG_DIRECTORY/client.toml"
APP_CONFIG_FILE="$CONFIG_DIRECTORY/app.toml"

# --------------------------------- run init --------------------------------- #
if ! command -v "$EXECUTABLE" >/dev/null; then
  echo "$EXECUTABLE does not exist"
  exit 1
fi

# TODO: run this check only if settlement is set to dymension
if ! command -v "$SETTLEMENT_EXECUTABLE" >/dev/null; then
  echo "$SETTLEMENT_EXECUTABLE does not exist"
  exit 1
fi

# Verify that a genesis file doesn't exists for the dymension chain
if [ -f "$GENESIS_FILE" ]; then
  printf "\n======================================================================================================\n"
  echo "A genesis file already exists [$GENESIS_FILE]. building the chain will delete all previous chain data. continue? (y/n)"
  printf "\n======================================================================================================\n"
  read -r answer
  if [ "$answer" != "${answer#[Yy]}" ]; then
    rm -rf "$CHAIN_DIR"
  else
    exit 1
  fi
fi


$EXECUTABLE dymint unsafe-reset-all  --home "$CHAIN_DIR"
$EXECUTABLE init "$MONIKER" --chain-id "$CHAIN_ID" --home "$CHAIN_DIR"

if [ -n "$LOG_FILE_PATH" ]; then
  mkdir -p "$(dirname "$LOG_FILE_PATH")" # create parent directories if they don't exist
  touch "$LOG_FILE_PATH" # create the file
  echo "Log file created at $LOG_FILE_PATH"
else
  echo "LOG_FILE_PATH is not set. using stdout"
fi


# ------------------------------- client config ------------------------------ #
$EXECUTABLE config keyring-backend test
$EXECUTABLE config chain-id "$CHAIN_ID"

# -------------------------------- app config -------------------------------- #
sed -i'' -e "s/^minimum-gas-prices *= .*/minimum-gas-prices = \"0$DENOM\"/" "$APP_CONFIG_FILE"
sed -i'' -e '/\[api\]/,+3 s/enable *= .*/enable = true/' "$APP_CONFIG_FILE"
sed -i'' -e "/\[api\]/,+9 s/address *= .*/address = \"tcp:\/\/$API_ADDRESS\"/" "$APP_CONFIG_FILE"
sed -i'' -e "/\[grpc\]/,+6 s/address *= .*/address = \"$GRPC_LADDRESS\"/" "$APP_CONFIG_FILE"
sed -i'' -e "/\[grpc-web\]/,+7 s/address *= .*/address = \"$GRPC_WEB_LADDRESS\"/" "$APP_CONFIG_FILE"
sed -i'' -e "/\[rpc\]/,+3 s/laddr *= .*/laddr = \"tcp:\/\/$RPC_LADDRESS\"/" "$TENDERMINT_CONFIG_FILE"
sed -i'' -e "/\[p2p\]/,+3 s/laddr *= .*/laddr = \"tcp:\/\/$P2P_LADDRESS\"/" "$TENDERMINT_CONFIG_FILE"
sed -i'' -e "s/^persistent_peers *= .*/persistent_peers = \"$ROLLAPP_PEERS\"/" "$TENDERMINT_CONFIG_FILE"

if [ -n "$UNSAFE_CORS" ]; then
  echo "Setting CORS"
  sed -ie 's/enabled-unsafe-cors.*$/enabled-unsafe-cors = true/' "$APP_CONFIG_FILE"
  sed -ie 's/enable-unsafe-cors.*$/enabled-unsafe-cors = true/' "$APP_CONFIG_FILE"
  sed -ie 's/cors_allowed_origins.*$/cors_allowed_origins = ["*"]/' "$TENDERMINT_CONFIG_FILE"
fi

# ------------------------------ genesis config ------------------------------ #
set_distribution_params
set_gov_params
set_minting_params

set_denom "$DENOM"

if [ "$EVM_ENABLED" ]; then
  set_EVM_params
fi

# --------------------- adding keys and genesis accounts --------------------- #

#local genesis account
$EXECUTABLE keys add "$KEY_NAME_ROLLAPP" --keyring-backend test --home "$CHAIN_DIR"
$EXECUTABLE add-genesis-account "$KEY_NAME_ROLLAPP" "$TOKEN_AMOUNT" --keyring-backend test --home "$CHAIN_DIR"

#If using settlement layer, make sure the sequencer account is funded
if [ "$SETTLEMENT_LAYER" = "dymension" ]; then
    #add account for sequencer on the hub
    $SETTLEMENT_EXECUTABLE keys add "$KEY_NAME_DYM" --keyring-backend test --keyring-dir $KEYRING_PATH
    SEQ_ACCOUNT_ON_HUB="$($SETTLEMENT_EXECUTABLE keys show -a $KEY_NAME_DYM --keyring-dir $KEYRING_PATH --keyring-backend test)"
    echo "Current balance of sequencer account on hub[$SEQ_ACCOUNT_ON_HUB]: "
    $SETTLEMENT_EXECUTABLE q bank balances "$SEQ_ACCOUNT_ON_HUB" --node "$SETTLEMENT_RPC"

    echo "Make sure the sequencer account [$SEQ_ACCOUNT_ON_HUB] is funded"
    echo "From within the hub node: \"$SETTLEMENT_EXECUTABLE tx bank send $KEY_NAME_GENESIS $SEQ_ACCOUNT_ON_HUB $SEQUENCER_AMOUNT --keyring-backend test\""
    echo "Press to continue..." 
    read -r answer
    fi


if [ "$ROLLAPP_PEERS" != "" ]; then
  printf "\n======================================================================================================"
  ROLLAPP_PEERS
  echo "ROLLAPP_PEERS defined, assuming full node for existing network"
  echo "To join existing chain, copy the genesis file to $GENESIS_FILE"
  exit 0
fi

echo "Do you want to include staker on genesis? (Y/n) "
read -r answer
if [ ! "$answer" != "${answer#[Nn]}" ] ;then
  $EXECUTABLE gentx "$KEY_NAME_ROLLAPP" "$STAKING_AMOUNT" --chain-id "$CHAIN_ID" --keyring-backend test --home "$CHAIN_DIR"
  $EXECUTABLE collect-gentxs --home "$CHAIN_DIR"
fi

$EXECUTABLE validate-genesis


echo "Do you want to register sequencer on genesis? (Y/n) "
read -r answer
if [ ! "$answer" != "${answer#[Nn]}" ] ;then
  set_sequencers
fi

