BASEDIR=$(dirname "$0")
source "$BASEDIR"/shared.sh
source "$BASEDIR"/ibc_transfer_commands.sh

IBC_PORT=transfer
IBC_VERSION=ics20-1

HUB_GENESIS_ADDR=$(dymd keys show $KEY_NAME_GENESIS -a --keyring-backend test)
ROLLAPP_GENESIS_ADDR=$(rollappd keys show $KEY_NAME_ROLLAPP -a --keyring-backend test --home $CHAIN_DIR)
TOKENS_AMOUNT=5555
URAP_AMOUNT="$TOKENS_AMOUNT"urap
DYM_AMOUNT="$TOKENS_AMOUNT"dym

usage ()
{
	echo "Usage: $(basename $0) <argument>"
	echo ""
	echo "argument :" 
    echo "-q: query balances of $KEY_NAME_GENESIS on hub and $KEY_NAME_ROLLAPP on rollapp"
    echo "rol2hub: ibc-transfer of $URAP_AMOUNT to $KEY_NAME_GENESIS from $KEY_NAME_ROLLAPP"
    echo "hub_back: transfer back the tokens from the hub to the rollapp"
	echo "hub2rol: ibc-transfer of $DYM_AMOUNT to $KEY_NAME_ROLLAPP from $KEY_NAME_GENESIS"
    echo "hub_back: transfer back the tokens from the hub to the rollapp"
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