package types

import (
	"strconv"

	"cosmossdk.io/collections"
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

	// Stub variable to store the operator address from the InitChain request
	InitChainStubAddr = "initchainstubaddr"
)

var (
	// Keys for store prefixes

	SequencersKey           = []byte{0x21} // prefix for each key to a sequencer
	SequencersByConsAddrKey = []byte{0x22} // prefix for each key to a sequencer index, by pubkey
	ValidatorUpdateKey      = []byte{0x23} // prefix for each key to a validator updates list
	SequencersRewardAddrKey = []byte{0x24} // prefix for each key to a validator reward addr
	WhitelistedRelayersKey  = []byte{0x25} // prefix for each key to a validator whitelisted relayer list

	HistoricalInfoKey = []byte{0x50} // prefix for the historical info
)

// GetSequencerKey creates the key for the sequencer with address
func GetSequencerKey(operatorAddr sdk.ValAddress) []byte {
	return append(SequencersKey, address.MustLengthPrefix(operatorAddr)...)
}

// GetSequencerByConsAddrKey creates the key for the sequencer with pubkey
func GetSequencerByConsAddrKey(addr sdk.ConsAddress) []byte {
	return append(SequencersByConsAddrKey, address.MustLengthPrefix(addr)...)
}

func GetSequencerRewardAddrKey(operator sdk.ValAddress) []byte {
	return append(SequencersRewardAddrKey, address.MustLengthPrefix(operator)...)
}

// GetHistoricalInfoKey returns a key prefix for indexing HistoricalInfo objects.
func GetHistoricalInfoKey(height int64) []byte {
	return append(HistoricalInfoKey, []byte(strconv.FormatInt(height, 10))...)
}

func WhitelistedRelayersPrefix() collections.Prefix {
	return collections.NewPrefix(WhitelistedRelayersKey)
}
