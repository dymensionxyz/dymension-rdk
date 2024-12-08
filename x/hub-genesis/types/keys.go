package types

import "cosmossdk.io/collections"

const (
	ModuleName = "hubgenesis"

	// StoreKey is the default store key for hub-genesis module.
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the hub-genesis module
	QuerierRoute = StoreKey
)

var (
	StateKey           = []byte{0x01}
	GenesisInfoKey     = []byte{0x02}
	OngoingChannelsKey = []byte{0x03}
)

func OngoingChannelsPrefix() collections.Prefix {
	return collections.NewPrefix(OngoingChannelsKey)
}
