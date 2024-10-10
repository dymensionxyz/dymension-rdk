package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func Bech32ToAddr[T sdk.AccAddress | sdk.ValAddress](addr string) (T, error) {
	_, bytes, err := bech32.DecodeAndConvert(addr)
	if err != nil {
		return nil, fmt.Errorf("decoding bech32 addr: %w", err)
	}
	return bytes, nil
}
