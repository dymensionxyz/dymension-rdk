#100DYM
AMOUNT="100000000udym"

# ---------------------- create sequencer account on hub --------------------- #
create_new_account() {
  echo "creating account"
  echo "$SETTLEMENT_EXECUTABLE keys add $KEY_NAME_DYM --keyring-backend test"
  $SETTLEMENT_EXECUTABLE keys add $KEY_NAME_DYM --keyring-backend test
}

# ------------------------------- fund account ------------------------------- #
fund_account() {
  echo "funding"
  NEW_ADDRESS="$("$SETTLEMENT_EXECUTABLE" keys show "$KEY_NAME_DYM" -a --keyring-backend test)"
  echo "$SETTLEMENT_EXECUTABLE tx bank send $KEY_NAME_GENESIS $NEW_ADDRESS $AMOUNT"
  $SETTLEMENT_EXECUTABLE tx bank send "$KEY_NAME_GENESIS" "$NEW_ADDRESS" "$AMOUNT" \
    --chain-id "$SETTLEMENT_CHAIN_ID" \
    --keyring-backend test \
    --broadcast-mode block
}


create_and_fund_seq_account() {
  create_new_account

  echo "waiting for keyring to be updated"
  #TODO: validate keyring is updated
  sleep 5
  fund_account
}