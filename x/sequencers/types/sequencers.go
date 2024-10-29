package types

import (
	"fmt"
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension-rdk/utils/addressutils"
)

const maxWhitelistedRelayers = 10

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

func ValidateWhitelistedRelayers(wlr []string) error {
	if len(wlr) > maxWhitelistedRelayers {
		return fmt.Errorf("maximum allowed relayers is %d", maxWhitelistedRelayers)
	}
	relayers := make(map[string]struct{}, len(wlr))
	for _, r := range wlr {
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
