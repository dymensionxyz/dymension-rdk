BASEDIR=$(dirname "$0")
echo "$BASEDIR"
source "$BASEDIR"/shared.sh

$EXECUTABLE tx sequencers create-sequencer \
  --pubkey $($EXECUTABLE dymint show-sequencer) \
  --broadcast-mode block \
  --moniker $MONIKER \
  --chain-id $CHAIN_ID \
  --from $(rollappd keys show -a ${KEY_NAME_ROLLAPP} --keyring-backend test) \
  --keyring-backend test