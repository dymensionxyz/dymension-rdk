BASEDIR=$(dirname "$0")
. "$BASEDIR"/../shared.sh
. "$BASEDIR"/ibc_transfer_commands.sh

IBC_PORT=transfer
IBC_VERSION=ics20-1

#TODO:fix this. The hub addr shouldn't depend on the local keyring
HUB_GENESIS_ADDR=$(dymd keys show $KEY_NAME_GENESIS -a --keyring-backend test)
ROLLAPP_GENESIS_ADDR=$($EXECUTABLE keys show $KEY_NAME_ROLLAPP -a --keyring-backend test --home $ROLLAPP_CHAIN_DIR)
TOKENS_AMOUNT=500000000
URAP_AMOUNT="$TOKENS_AMOUNT""$DENOM"
DYM_AMOUNT="$TOKENS_AMOUNT"udym

usage ()
{
	echo "Usage: $(basename $0) <argument>"
	echo ""
	echo "argument :" 
    echo "-q: query balances of $KEY_NAME_GENESIS on hub and $KEY_NAME_ROLLAPP on rollapp"
    echo "rol2hub: ibc-transfer of $URAP_AMOUNT to $KEY_NAME_GENESIS from $KEY_NAME_ROLLAPP"
    echo "hub_back: transfer back the tokens from the hub to the rollapp"
	echo "hub2rol: ibc-transfer of $DYM_AMOUNT to $KEY_NAME_ROLLAPP from $KEY_NAME_GENESIS"
    echo "rol_back: transfer back the tokens from the rollapp to the hub"
	echo ""

	exit 1
}


if [ $# -eq 0 ]; then
	echo "No arguments supplied"
    usage
	exit 1
fi

case $1 in
	"-q")
		query_test_accounts
		exit 0
		;;
	rol2hub)
		from_rollapp_2_hub
		exit 0
		;;
    hub_back)
		from_hub_back_2_rollapp
		exit 0
		;;
	hub2rol)
		from_hub_2_rollapp
		exit 0
		;;
	rol_back)
		from_rollapp_back_2_hub
		exit 0
		;;
    *)
        echo "Unknown"
		usage
		exit 1
esac