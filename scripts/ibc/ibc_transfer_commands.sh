
# ------------------------ moving urap to hub and back ----------------------- #
from_rollapp_2_hub() {
    AMOUNT=$URAP_AMOUNT
    echo "sending $AMOUNT to $HUB_GENESIS_ADDR from $KEY_NAME_ROLLAPP"

    $EXECUTABLE tx ibc-transfer transfer "$IBC_PORT" "$ROLLAPP_CHANNEL_NAME" "$HUB_GENESIS_ADDR" "$AMOUNT" \
    --from $KEY_NAME_ROLLAPP \
    --chain-id "$ROLLAPP_CHAIN_ID" \
    --broadcast-mode block \
    --keyring-backend test \
    --home $ROLLAPP_CHAIN_DIR
}


from_hub_back_2_rollapp() {
    URAP_HASHED_DENOM=$(dymd q ibc-transfer denom-hash $IBC_PORT/$HUB_CHANNEL_NAME/urap | cut -d ' ' -f 2)
    AMOUNT="$TOKENS_AMOUNT"ibc/"$URAP_HASHED_DENOM"
    echo "sending $AMOUNT to $ROLLAPP_GENESIS_ADDR from $KEY_NAME_GENESIS"

    $SETTLEMENT_EXECUTABLE tx ibc-transfer transfer "$IBC_PORT" "$HUB_CHANNEL_NAME" "$ROLLAPP_GENESIS_ADDR" "$AMOUNT" \
    --from $KEY_NAME_GENESIS \
    --chain-id "$SETTLEMENT_CHAIN_ID" --broadcast-mode block --keyring-backend test
}



# ---------------------- moving DYM to rollapp and back ---------------------- #
from_hub_2_rollapp() {
    AMOUNT=$DYM_AMOUNT
    echo "sending $AMOUNT to $ROLLAPP_GENESIS_ADDR from $KEY_NAME_GENESIS"

    $SETTLEMENT_EXECUTABLE tx ibc-transfer transfer "$IBC_PORT" "$HUB_CHANNEL_NAME" "$ROLLAPP_GENESIS_ADDR" "$AMOUNT" \
    --from $KEY_NAME_GENESIS \
    --chain-id "$SETTLEMENT_CHAIN_ID" \
    --broadcast-mode block \
    --packet-timeout-timestamp 100000000000000000 \
    --packet-timeout-height 0-0 \
    --keyring-backend test
}

from_rollapp_back_2_hub() {
    DYM_HASHED_DENOM=$(rollappd q ibc-transfer denom-hash $IBC_PORT/$ROLLAPP_CHANNEL_NAME/udym | cut -d ' ' -f 2)
    AMOUNT="$TOKENS_AMOUNT"ibc/"$DYM_HASHED_DENOM"
    echo "sending $AMOUNT to $HUB_GENESIS_ADDR from $KEY_NAME_ROLLAPP"

    $EXECUTABLE tx ibc-transfer transfer "$IBC_PORT" "$ROLLAPP_CHANNEL_NAME" "$HUB_GENESIS_ADDR" "$AMOUNT" \
        --from $KEY_NAME_ROLLAPP \
        --chain-id "$ROLLAPP_CHAIN_ID" \
        --broadcast-mode block \
        --keyring-backend test \
        --home $ROLLAPP_CHAIN_DIR
}


query_test_accounts() {
    echo '# ------------------------------------ .. ------------------------------------ #'
    echo "Rollapp:"
    $EXECUTABLE q bank balances $ROLLAPP_GENESIS_ADDR --node "$ROLLAPP_RPC_FOR_RELAYER"

    echo '# ------------------------------------ .. ------------------------------------ #'
    echo "Hub:"
    dymd q bank balances $HUB_GENESIS_ADDR --node "$SETTLEMENT_RPC_FOR_RELAYER"

    # echo '# ------------------------------------ .. ------------------------------------ #'
    # echo "Rollapp1:"
    # COUNTERPARTY_ROLLAPP_GENESIS_ADDR=$(rollappd keys show $KEY_NAME_ROLLAPP -a --keyring-backend test)
    # rollappd q bank balances $COUNTERPARTY_ROLLAPP_GENESIS_ADDR
}


