package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/rollapp/x/sequencers/keeper"
	"github.com/dymensionxyz/rollapp/x/sequencers/testutils"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	k, ctx := testutils.NewTestSequencer(t)

	q := keeper.Querier{Keeper: *k}

	wctx := sdk.WrapSDKContext(ctx)

	response, err := q.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: types.DefaultParams()}, response)
}
