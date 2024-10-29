package types_test

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func TestMustNewWhitelistedRelayers(t *testing.T) {
	// generate random addresses
	addr1 := utils.AccAddress()
	addr2 := utils.AccAddress()
	addr3 := utils.AccAddress()
	addr4 := utils.AccAddress()
	addr5 := utils.AccAddress()

	// create a list of relayers
	relayers := []string{
		addr1.String(),
		addr2.String(),
		addr3.String(),
		addr4.String(),
		addr5.String(),
	}

	// form a slice of whitelisted relayers sdk.AccAddress
	expected := types.MustNewWhitelistedRelayers(relayers)

	// check that the list matches the initial set of addrs
	require.ElementsMatch(t, relayers, expected.Relayers)

	// deterministic test: validate that the return value is independent
	// of the order of inputted relater addresses
	for i := 0; i < 100; i++ {
		// shuffle the slice
		rand.Shuffle(len(relayers), func(i, j int) {
			relayers[i], relayers[j] = relayers[j], relayers[i]
		})

		actual := types.MustNewWhitelistedRelayers(relayers)
		require.Equal(t, expected, actual)
	}
}
