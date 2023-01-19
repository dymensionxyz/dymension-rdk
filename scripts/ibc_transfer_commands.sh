
# ------------------------ moving urap to hub and back ----------------------- #
from_rollapp_2_hub() {
    AMOUNT=$URAP_AMOUNT
    echo "sending $AMOUNT to $HUB_GENESIS_ADDR from $KEY_NAME_ROLLAPP"

    $EXECUTABLE tx ibc-transfer transfer "$IBC_PORT" "$ROLLAPP_CHANNEL_NAME" "$HUB_GENESIS_ADDR" "$AMOUNT" \
    --from $KEY_NAME_ROLLAPP \
    --chain-id "$CHAIN_ID" \
    --broadcast-mode block \
    --keyring-backend test \
    --home $CHAIN_DIR
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
    --keyring-backend test
}

from_rollapp_back_2_hub() {
    DYM_HASHED_DENOM=$(rollappd q ibc-transfer denom-hash $IBC_PORT/$ROLLAPP_CHANNEL_NAME/dym | cut -d ' ' -f 2)
    AMOUNT="$TOKENS_AMOUNT"ibc/"$DYM_HASHED_DENOM"
    echo "sending $AMOUNT to $HUB_GENESIS_ADDR from $KEY_NAME_ROLLAPP"

    $EXECUTABLE tx ibc-transfer transfer "$IBC_PORT" "$ROLLAPP_CHANNEL_NAME" "$HUB_GENESIS_ADDR" "$AMOUNT" \
        --from $KEY_NAME_ROLLAPP \
        --chain-id "$CHAIN_ID" \
        --broadcast-mode block \
        --keyring-backend test \
        --home $CHAIN_DIR
}


#TODO: this multi party test should be moved to testing/infra repo
# from_rollappX_2_rollappY_through_hub() {
#     # TODO: take address from command
#     MEMO='{"forward":{"receiver":"rol1g0f0cth5acca6agtlshv25kpxf33j3kdkkdzxs","port":"transfer","channel":"channel-0"}}'

#     rollappd tx ibc-transfer transfer "$IBC_PORT" "$ROLLAPP_CHANNEL_NAME" "$HUB_GENESIS_ADDR" 95urap \
#     --from $KEY_NAME_ROLLAPP \
#     --chain-id "$CHAIN_ID" \
#      --broadcast-mode block \
#      --keyring-backend test \
#       --home $CHAIN_DIR \
#       --memo $MEMO
# }


query_test_accounts() {
    echo '# ------------------------------------ .. ------------------------------------ #'
    echo "Rollapp:"
    rollappd q bank balances $ROLLAPP_GENESIS_ADDR --home $CHAIN_DIR

    echo '# ------------------------------------ .. ------------------------------------ #'
    echo "Hub:"
    dymd q bank balances $HUB_GENESIS_ADDR

    # echo '# ------------------------------------ .. ------------------------------------ #'
    # echo "Rollapp1:"
    # COUNTERPARTY_ROLLAPP_GENESIS_ADDR=$(rollappd keys show $KEY_NAME_ROLLAPP -a --keyring-backend test)
    # rollappd q bank balances $COUNTERPARTY_ROLLAPP_GENESIS_ADDR
}


