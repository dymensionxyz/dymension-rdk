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

	HubKeyPrefix = "Hub/value/"
)

func HubKey(
	hubId string,
) []byte {
	var key []byte

	hubIdBytes := []byte(hubId)
	key = append(key, hubIdBytes...)
	key = append(key, []byte("/")...)

	return key
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}
