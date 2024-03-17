package types

const (
	// module name
	ModuleName = "denommetadata"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName

	// PermissionedAddressesKey is the key for the permissioned addresses
	PermissionedAddressesKey = "PermissionedAddresses"
)

// KVStore key prefixes
var (
	ParamsKey = []byte{0x00} // Prefix for params key
)
