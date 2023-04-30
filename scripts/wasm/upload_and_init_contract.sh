#!/bin/sh

ROLLAPP_CHAIN_ID=${ROLLAPP_CHAIN_ID:-rollapp}
MONIKER=${MONIKER:-rollapp-validator}
KEY_NAME=${KEY_NAME:-validator}
NAMESPACE_ID=${NAMESPACE_ID:-000000000000FFFF}
TOKEN_AMOUNT=${TOKEN_AMOUNT:-1000000000000000000000urap}
STAKING_AMOUNT=${STAKING_AMOUNT:-600000000000000000000urap}

# Setting up the correct parameters
export TXFLAG="--chain-id $ROLLAPP_CHAIN_ID --gas-prices 0.25urap --gas auto --gas-adjustment 1.3 -y --output json -b block"

# Storing the binary on chain
RES=$(rollappd tx wasm store cw20_base.wasm --from alice $TXFLAG)

# Getting the code id for the stored binary
CODE_ID=$(echo $RES | jq -r '.logs[0].events[1].attributes[0].value')

# Printing the code id
echo $CODE_ID

# Querying the list of contracts instantiated with the code id above
rollappd query wasm list-contract-by-code $CODE_ID --output json

# This should return an empty list, for now.