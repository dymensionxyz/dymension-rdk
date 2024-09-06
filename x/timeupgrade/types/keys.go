package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "timeupgrade"

	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	// Keys for store prefixes

	SequencersKey = []byte{0x21} // prefix for each key to a sequencer
)

// GetSequencerKey creates the key for the sequencer with address
func GetSequencerKey(operatorAddr sdk.ValAddress) []byte {
	return append(SequencersKey, address.MustLengthPrefix(operatorAddr)...)
}
