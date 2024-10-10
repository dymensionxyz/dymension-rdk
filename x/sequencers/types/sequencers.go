package types

import (
	"fmt"
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/utils/addressutils"
)

func MustNewWhitelistedRelayers(relayers []string) WhitelistedRelayers {
	convertedRelayers := make([]string, 0, len(relayers))
	for _, r := range relayers {
		relayer, err := addressutils.Bech32ToAddr[sdk.AccAddress](r)
		if err != nil {
			panic(fmt.Errorf("convert bech32 to relayer address: %s: %w", r, err))
		}
		convertedRelayers = append(convertedRelayers, relayer.String())
	}
	slices.Sort(convertedRelayers)
	return WhitelistedRelayers{Relayers: convertedRelayers}
}

func (wr WhitelistedRelayers) Validate() error {
	relayers := make(map[string]struct{}, len(wr.Relayers))
	for _, r := range wr.Relayers {
		if _, ok := relayers[r]; ok {
			return fmt.Errorf("duplicated relayer: %s", r)
		}
		relayers[r] = struct{}{}

		relayer, err := addressutils.Bech32ToAddr[sdk.AccAddress](r)
		if err != nil {
			return fmt.Errorf("convert bech32 to relayer address: %s: %w", r, err)
		}
		err = sdk.VerifyAddressFormat(relayer)
		if err != nil {
			return fmt.Errorf("invalid relayer address: %s: %w", r, err)
		}
	}
	return nil
}
