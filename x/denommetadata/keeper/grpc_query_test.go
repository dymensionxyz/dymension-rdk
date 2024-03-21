package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/testutils"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	k, ctx := testutils.NewTestDenommetadataKeeper(t)

	q := keeper.Querier{Keeper: *k}

	wctx := sdk.WrapSDKContext(ctx)

	response, err := q.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: types.DefaultParams()}, response)
}
