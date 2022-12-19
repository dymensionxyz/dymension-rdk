BASEDIR=$(dirname "$0")
echo "$BASEDIR"
source "$BASEDIR"/shared.sh

#Register Sequencer
DESCRIPTION="{\"Moniker\":\"$MONIKER_NAME\",\"Identity\":\"\",\"Website\":\"\",\"SecurityContact\":\"\",\"Details\":\"\"}";

#Use default keys of the settlement node
SEQ_ADDRESS="$("$SETTLEMENT_EXECUTABLE" keys show "$KEY_NAME_DYM" -a --keyring-backend test)"
SEQ_PUB_KEY="$("$SETTLEMENT_EXECUTABLE" keys show "$KEY_NAME_DYM" -p --keyring-backend test)"


##TODO: check if KEY_NAME_DYM exist. if not create and fund


$SETTLEMENT_EXECUTABLE tx sequencer create-sequencer "$SEQ_ADDRESS" "$SEQ_PUB_KEY" "$ROLLAPP_ID" "$DESCRIPTION" \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test


echo "check expected result with:"
echo "$SETTLEMENT_EXECUTABLE q sequencer show-sequencer $SEQ_ADDRESS"