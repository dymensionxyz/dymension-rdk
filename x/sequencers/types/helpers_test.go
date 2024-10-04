package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/stretchr/testify/require"
)

func TestKeyConversions(t *testing.T) {
	addr := "dym1e2xpcrtxgl9cmug2w4fck3k2mcuktqa0gkme0p"

	hrp, dymAddrBytes, err := bech32.DecodeAndConvert(addr)
	require.NoError(t, err)
	require.Equal(t, "dym", hrp)

	newAddr, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32AccountAddrPrefix(), dymAddrBytes)
	require.NoError(t, err)

	t.Log(newAddr)
}
