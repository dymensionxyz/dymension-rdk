package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/tendermint/tendermint/crypto"
)

type AddressType int32

const (
	// the 32 bytes length address type of ADR 028.
	AddressType32Bytes AddressType = 0
	// the default 20 bytes length address type.
	AddressType20Bytes AddressType = 1
)

// DeriveAddress derives an address with the given address length type, module name, and
// address derivation name
func DeriveAddress(addressType AddressType, moduleName, name string) sdk.AccAddress {
	switch addressType {
	case AddressType32Bytes:
		return address.Module(moduleName, []byte(name))
	case AddressType20Bytes:
		return sdk.AccAddress(crypto.AddressHash([]byte(moduleName + name)))
	default:
		return sdk.AccAddress{}
	}
}
