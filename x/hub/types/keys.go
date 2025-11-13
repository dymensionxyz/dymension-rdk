package types

const (
	// ModuleName defines the module name.
	ModuleName = "hubs"
	// StoreKey defines the primary store key.
	StoreKey = ModuleName
	// RouterKey is the message route for hub genesis.
	RouterKey = ModuleName
	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey
)

const (
	RegisteredHubDenomsKeyPrefix   = "registeredHubDenoms/value/"
	DecimalConversionPairKeyPrefix = "decimalConversionPair/value/"
)
