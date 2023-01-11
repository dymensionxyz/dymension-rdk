BASEDIR=$(dirname "$0")
echo "$BASEDIR"
source "$BASEDIR"/shared.sh

#Register Sequencer
DESCRIPTION="{\"Moniker\":\"$MONIKER_NAME\",\"Identity\":\"\",\"Website\":\"\",\"SecurityContact\":\"\",\"Details\":\"\"}";

#Use default keys of the settlement node
SEQ_PUB_KEY="$(${EXECUTABLE} dymint show-sequencer)"


##TODO: check if KEY_NAME_DYM exist. if not create and fund


$SETTLEMENT_EXECUTABLE tx sequencer create-sequencer "$SEQ_PUB_KEY" "$ROLLAPP_ID" "$DESCRIPTION" \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test \
  --broadcast-mode block