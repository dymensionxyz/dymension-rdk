BASEDIR=$(dirname "$0")
source "$BASEDIR"/shared.sh
source "$BASEDIR"/create_and_fund_seq_account.sh

#Register Sequencer
DESCRIPTION="{\"Moniker\":\"$MONIKER_NAME\",\"Identity\":\"\",\"Website\":\"\",\"SecurityContact\":\"\",\"Details\":\"\"}";
SEQ_PUB_KEY="$($EXECUTABLE dymint show-sequencer --home $CHAIN_DIR)"

key_name_dym_exists="$SETTLEMENT_EXECUTABLE keys show $KEY_NAME_DYM --keyring-backend test"
if $key_name_dym_exists > /dev/null; then
  echo "$KEY_NAME_DYM EXIST!"
else
  echo "$KEY_NAME_DYM not found - creating and funding"
  create_and_fund_seq_account
fi

$SETTLEMENT_EXECUTABLE tx sequencer create-sequencer "$SEQ_PUB_KEY" "$ROLLAPP_ID" "$DESCRIPTION" \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test \
  --broadcast-mode block