package keeper_test

import (
	"testing"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/testutils"
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
				MaxSequencers:     5,
				HistoricalEntries: 999,
			},
		},
	}
	for _, tC := range testCases {
		k, ctx := testutils.NewTestSequencerKeeper(t)
		k.SetParams(ctx, tC.params)
		t.Run(tC.desc, func(t *testing.T) {
			require.EqualValues(t, tC.params, k.GetParams(ctx))
		})
	}
}
