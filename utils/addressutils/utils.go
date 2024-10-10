package addressutils

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/tendermint/tendermint/crypto"
)

type AddressType int32

const (
	// AddressType32Bytes is the 32 bytes length address type of ADR 028.
	AddressType32Bytes AddressType = 0
	// AddressType20Bytes is the default 20 bytes length address type.
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

// Bech32ToAddr casts an arbitrary-prefixed bech32 string to either sdk.AccAddress or sdk.ValAddress.
func Bech32ToAddr[T sdk.AccAddress | sdk.ValAddress](addr string) (T, error) {
	_, bytes, err := bech32.DecodeAndConvert(addr)
	if err != nil {
		return nil, fmt.Errorf("decoding bech32 addr: %w", err)
	}
	return bytes, nil
}
