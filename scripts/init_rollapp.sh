BASEDIR=$(dirname "$0")
source "$BASEDIR"/shared.sh

# ---------------------------- initial parameters ---------------------------- #
TOKEN_AMOUNT=${TOKEN_AMOUNT:-1000000000000000000000urap}
STAKING_AMOUNT=${STAKING_AMOUNT:-500000000000000000000urap}

CONFIG_DIRECTORY="$CHAIN_DIR/config"
GENESIS_FILE="$CONFIG_DIRECTORY/genesis.json"
TENDERMINT_CONFIG_FILE="$CONFIG_DIRECTORY/config.toml"
CLIENT_CONFIG_FILE="$CONFIG_DIRECTORY/client.toml"
APP_CONFIG_FILE="$CONFIG_DIRECTORY/app.toml"

# --------------------------------- run init --------------------------------- #
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

# ------------------------------- client config ------------------------------ #
sed -i'' -e "s/^chain-id *= .*/chain-id = \"$CHAIN_ID\"/" "$CLIENT_CONFIG_FILE"
sed -i'' -e "s/^node *= .*/node = \"tcp:\/\/$RPC_PORT\"/" "$CLIENT_CONFIG_FILE"

# -------------------------------- app config -------------------------------- #
sed -i'' -e 's/^minimum-gas-prices *= .*/minimum-gas-prices = "0urap"/' "$APP_CONFIG_FILE"

# ------------------------------ genesis config ------------------------------ #
sed -i'' -e 's/bond_denom": ".*"/bond_denom": "urap"/' "$GENESIS_FILE"
sed -i'' -e 's/mint_denom": ".*"/mint_denom": "urap"/' "$GENESIS_FILE"
#TODO: set genesis params (rewards distribution, infaltion, staking denom)


$EXECUTABLE keys add "$KEY_NAME_ROLLAPP" --keyring-backend test --home "$CHAIN_DIR"
$EXECUTABLE add-genesis-account "$KEY_NAME_ROLLAPP" "$TOKEN_AMOUNT" --keyring-backend test --home "$CHAIN_DIR"

printf "\n\n======================================================================================================\n"
echo "Press any key to continue generating genesis validator..."
printf "======================================================================================================\n\n"
read -p ""

$EXECUTABLE gentx "$KEY_NAME_ROLLAPP" "$STAKING_AMOUNT" --chain-id "$CHAIN_ID" --keyring-backend test --home "$CHAIN_DIR"
$EXECUTABLE collect-gentxs --home "$CHAIN_DIR"