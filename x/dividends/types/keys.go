package types

var (
	// ModuleName defines the module name.
	ModuleName = "dividends"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for slashing.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName

	// KeyLastGaugeID defines key for setting last gauge ID.
	KeyLastGaugeID = []byte{0x01}

	// KeyPrefixGauges defines prefix key for storing reference key for all gauges.
	KeyPrefixGauges = []byte{0x02}
)
