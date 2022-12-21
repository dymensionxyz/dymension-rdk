package types

import (
	"strconv"
)

const (
	// ModuleName defines the module name
	ModuleName = "sequencers"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_sequencers"
)

// Iota generated keys/key prefixes (as a byte), supports 256 possible values
const (
	// PortByteKey defines the byte key to store the port ID in store
	PortByteKey byte = iota

	// LastDistributionTransmissionByteKey defines the byte key to store the last distribution transmission
	LastDistributionTransmissionByteKey

	// UnbondingTimeKeyString is the byte key for storing the unbonding period
	UnbondingTimeByteKey

	// ProviderClientKeyString is the byte key for storing the clientID of the provider client
	ProviderClientByteKey

	// ProviderChannelKeyString is the byte key for storing the channelID of the CCV channel
	ProviderChannelByteKey

	// PendingChangesKeyString is the byte key that will store any pending validator set changes
	// received over CCV channel but not yet flushed over ABCI
	PendingChangesByteKey

	// HistoricalInfoKey is the byte prefix that will store the historical info for a given height
	HistoricalInfoBytePrefix

	// PacketMaturityTimePrefix is the byte prefix that will store maturity time for each received VSC packet
	PacketMaturityTimeBytePrefix

	// HeightValsetUpdateIDPrefix is the byte prefix that will store the mapping from block height to valset update ID
	HeightValsetUpdateIDBytePrefix

	// OutstandingDowntimePrefix is the byte prefix that will store the validators outstanding downtime by consensus address
	OutstandingDowntimeBytePrefix

	// PendingSlashRequestsPrefix is the byte prefix that will store a list of slash request that must be sent
	// to the provider chain once the CCV channel is established
	PendingSlashRequestsBytePrefix

	// CrossChainValidatorPrefix is the byte prefix that will store cross-chain validators by consensus address
	CrossChainValidatorBytePrefix
)

var (
	// Keys for store prefixes
	// Last* values are constant during a block.
	LastValidatorPowerKey = []byte{0x11} // prefix for each key to a validator index, for bonded validators
	LastTotalPowerKey     = []byte{0x12} // prefix for the total power

	ValidatorsKey             = []byte{0x21} // prefix for each key to a validator
	ValidatorsByConsAddrKey   = []byte{0x22} // prefix for each key to a validator index, by pubkey
	ValidatorsByPowerIndexKey = []byte{0x23} // prefix for each key to a validator index, sorted by power

	DelegationKey                    = []byte{0x31} // key for a delegation
	UnbondingDelegationKey           = []byte{0x32} // key for an unbonding-delegation
	UnbondingDelegationByValIndexKey = []byte{0x33} // prefix for each key for an unbonding-delegation, by validator operator
	RedelegationKey                  = []byte{0x34} // key for a redelegation
	RedelegationByValSrcIndexKey     = []byte{0x35} // prefix for each key for an redelegation, by source validator operator
	RedelegationByValDstIndexKey     = []byte{0x36} // prefix for each key for an redelegation, by destination validator operator

	UnbondingQueueKey    = []byte{0x41} // prefix for the timestamps in unbonding queue
	RedelegationQueueKey = []byte{0x42} // prefix for the timestamps in redelegations queue
	ValidatorQueueKey    = []byte{0x43} // prefix for the timestamps in validator queue

	HistoricalInfoKey = []byte{0x50} // prefix for the historical info
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetHistoricalInfoKey returns a key prefix for indexing HistoricalInfo objects.
func GetHistoricalInfoKey(height int64) []byte {
	return append(HistoricalInfoKey, []byte(strconv.FormatInt(height, 10))...)
}
