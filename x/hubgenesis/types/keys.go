package types

const (
	// module name.
	ModuleName = "hubgenesis"

	// StoreKey is the default store key for mint.
	StoreKey = ModuleName

	// RouterKey is the message route for hub genesis.
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the minting querier.
	QueryParameters = "parameters"
)

var StateKey = []byte{0x01}
