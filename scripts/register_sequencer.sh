BASEDIR=$(dirname "$0")
echo "$BASEDIR"
source "$BASEDIR"/shared.sh

#Register Sequencer
DESCRIPTION="{\"Moniker\":\"$MONIKER_NAME\",\"Identity\":\"\",\"Website\":\"\",\"SecurityContact\":\"\",\"Details\":\"\"}";
CREATOR_ADDRESS="$("$SETTLEMENT_EXECUTABLE" keys show "$KEY_NAME_DYM" -a --keyring-backend test)"
CREATOR_PUB_KEY="$("$SETTLEMENT_EXECUTABLE" keys show "$KEY_NAME_DYM" -p --keyring-backend test)"


$SETTLEMENT_EXECUTABLE tx sequencer create-sequencer "$CREATOR_ADDRESS" "$CREATOR_PUB_KEY" "$ROLLAPP_ID" "$DESCRIPTION" \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test


