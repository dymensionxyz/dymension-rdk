source shared.sh

dymd tx rollapp create-rollapp "$ROLLAPP_ID" stamp1 "genesis-path/1" 3 100 '{"Addresses":[]}' \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test




#Sequencer register
DESCRIPTION="{\"Moniker\":\"$MONIKER_NAME\",\"Identity\":\"\",\"Website\":\"\",\"SecurityContact\":\"\",\"Details\":\"\"}";
CREATOR_ADDRESS="$(dymd keys show "$KEY_NAME_DYM" -a --keyring-backend test)"
CREATOR_PUB_KEY="$(dymd keys show "$KEY_NAME_DYM" -p --keyring-backend test)"


echo "$CREATOR_ADDRESS" 
echo "$CREATOR_PUB_KEY"


dymd tx sequencer create-sequencer "$CREATOR_ADDRESS" "$CREATOR_PUB_KEY" "$ROLLAPP_ID" "$DESCRIPTION" \
  --from "$KEY_NAME_DYM" \
  --chain-id "$SETTLEMENT_CHAIN_ID" \
  --keyring-backend test


