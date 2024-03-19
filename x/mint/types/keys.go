package types

const (
	// module name.
	ModuleName = "mint"

	// StoreKey is the default store key for mint.
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey
)

// MinterKey is the key to use for the keeper store.
var MinterKey = []byte{0x00}
