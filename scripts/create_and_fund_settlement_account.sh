BASEDIR=$(dirname "$0")
echo "$BASEDIR"
source "$BASEDIR"/shared.sh

AMOUNT="100000000dym"


#TODO: allow to get accounts from cmdline arguments
echo "genesis account: $KEY_NAME_GENESIS"
echo "new account: $KEY_NAME_DYM"


#add new user to the keyring
$SETTLEMENT_EXECUTABLE keys add $KEY_NAME_DYM --keyring-backend test

#send funds
GENESIS_ADDRESS="$("$SETTLEMENT_EXECUTABLE" keys show "$KEY_NAME_GENESIS" -a --keyring-backend test)"
NEW_ADDRESS="$("$SETTLEMENT_EXECUTABLE" keys show "$KEY_NAME_DYM" -a --keyring-backend test)"



echo "$SETTLEMENT_EXECUTABLE tx bank send $GENESIS_ADDRESS $NEW_ADDRESS $AMOUNT"
$SETTLEMENT_EXECUTABLE tx bank send "$GENESIS_ADDRESS" "$NEW_ADDRESS" "$AMOUNT" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test



echo "check expected result with:"
echo "$SETTLEMENT_EXECUTABLE q bank balances $NEW_ADDRESS --chain-id $SETTLEMENT_CHAIN_ID"