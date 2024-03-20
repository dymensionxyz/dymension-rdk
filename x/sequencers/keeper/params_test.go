package keeper_test

import (
	"testing"

	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	testCases := []struct {
		desc   string
		params types.Params
	}{
		{
			desc:   "Default params",
			params: types.DefaultParams(),
		},
		{
			desc: "non-default params",
			params: types.Params{
				UnbondingTime:     100,
				HistoricalEntries: 999,
			},
		},
	}
	for _, tC := range testCases {
		app := utils.Setup(t, false)
		k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

		k.SetParams(ctx, tC.params)
		t.Run(tC.desc, func(t *testing.T) {
			require.EqualValues(t, tC.params, k.GetParams(ctx))
		})
	}
}
