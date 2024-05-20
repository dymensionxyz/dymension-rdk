package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name.
	ModuleName = "gasless"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for slashing.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName
)

var (
	TxGTIDsKeyPrefix     = []byte{0xa0}
	LastGasTankIDKey     = []byte{0xa1}
	GasTankKeyPrefix     = []byte{0xa2}
	GasConsumerKeyPrefix = []byte{0xa3}
)

func GetTxGTIDsKey(txPathOrContractAddress string) []byte {
	return append(TxGTIDsKeyPrefix, LengthPrefixString(txPathOrContractAddress)...)
}

func GetAllTxGTIDsKey() []byte {
	return TxGTIDsKeyPrefix
}

func GetLastGasTankIDKey() []byte {
	return LastGasTankIDKey
}

func GetGasTankKey(gasTankID uint64) []byte {
	return append(GasTankKeyPrefix, sdk.Uint64ToBigEndian(gasTankID)...)
}

func GetAllGasTanksKey() []byte {
	return GasTankKeyPrefix
}

func GetGasConsumerKey(consumer sdk.AccAddress) []byte {
	return append(GasConsumerKeyPrefix, address.MustLengthPrefix(consumer)...)
}

func GetAllGasConsumersKey() []byte {
	return GasConsumerKeyPrefix
}

// LengthPrefixString returns length-prefixed bytes representation
// of a string.
func LengthPrefixString(s string) []byte {
	bz := []byte(s)
	bzLen := len(bz)
	if bzLen == 0 {
		return bz
	}
	return append([]byte{byte(bzLen)}, bz...)
}
