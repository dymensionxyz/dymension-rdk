package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func AddrBytesToAddrString(addr []byte) (string, error) {
	newAddr, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), addr)
	if err != nil {
		return "", err
	}
	_, err = sdk.AccAddressFromBech32(newAddr)
	if err != nil {
		return "", err
	}
	return newAddr, nil
}
