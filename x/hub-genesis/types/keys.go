package types

const (
	ModuleName = "hubgenesis"

	// StoreKey is the default store key for hub-genesis module.
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the hub-genesis module
	QuerierRoute = StoreKey
)

var (
	StateKey       = []byte{0x01}
	GenesisInfoKey = []byte{0x02}
)
