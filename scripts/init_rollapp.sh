BASEDIR=$(dirname "$0")
echo "$BASEDIR"
source "$BASEDIR"/shared.sh

# ---------------------------- initial parameters ---------------------------- #
TOKEN_AMOUNT=${TOKEN_AMOUNT:-1000000000000000000000urap}
STAKING_AMOUNT=${STAKING_AMOUNT:-500000000000000000000urap}
GENESIS_FILE="$CHAIN_DIR"/config/genesis.json


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


$EXECUTABLE tendermint unsafe-reset-all
$EXECUTABLE init "$MONIKER" --chain-id "$CHAIN_ID"


sed -i'' -e 's/^minimum-gas-prices *= .*/minimum-gas-prices = "0urap"/' "$CHAIN_DIR"/config/app.toml
sed -i'' -e 's/bond_denom": ".*"/bond_denom": "urap"/' "$CHAIN_DIR"/config/genesis.json
sed -i'' -e 's/mint_denom": ".*"/mint_denom": "urap"/' "$CHAIN_DIR"/config/genesis.json


#TODO: set rewards precentegas correctly

$EXECUTABLE keys add "$KEY_NAME_ROLLAPP" --keyring-backend test
$EXECUTABLE add-genesis-account "$KEY_NAME_ROLLAPP" "$TOKEN_AMOUNT" --keyring-backend test

read -p "Press any key to continue generating genesis validator..."


$EXECUTABLE gentx "$KEY_NAME_ROLLAPP" "$STAKING_AMOUNT" --chain-id "$CHAIN_ID" --keyring-backend test
$EXECUTABLE collect-gentxs