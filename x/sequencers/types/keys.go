package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "sequencers"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the module's message route key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

var (
	// Keys for store prefixes
	ValidatorsKey           = []byte{0x21} // prefix for each key to a validator
	ValidatorsByConsAddrKey = []byte{0x22} // prefix for each key to a validator index, by pubkey

	DymintSeqKey = []byte{0x23} // prefix for each key sequencer registerd on dymint, by pubkey

	HistoricalInfoKey = []byte{0x50} // prefix for the historical info
)

// GetValidatorKey creates the key for the validator with address
// VALUE: staking/Validator
func GetValidatorKey(operatorAddr sdk.ValAddress) []byte {
	return append(ValidatorsKey, address.MustLengthPrefix(operatorAddr)...)
}

// GetValidatorByConsAddrKey creates the key for the validator with pubkey
// VALUE: validator operator address ([]byte)
func GetValidatorByConsAddrKey(addr sdk.ConsAddress) []byte {
	return append(ValidatorsByConsAddrKey, address.MustLengthPrefix(addr)...)
}

// DymintSeqKey creates the key for the validator with address
// VALUE: staking/Validator
func GetDymintSeqKey(addr sdk.ConsAddress) []byte {
	return append(DymintSeqKey, address.MustLengthPrefix(addr)...)
}

// GetHistoricalInfoKey returns a key prefix for indexing HistoricalInfo objects.
func GetHistoricalInfoKey(height int64) []byte {
	return append(HistoricalInfoKey, []byte(strconv.FormatInt(height, 10))...)
}
