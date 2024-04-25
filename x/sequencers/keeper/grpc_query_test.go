package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	q := keeper.Querier{Keeper: *k}

	wctx := sdk.WrapSDKContext(ctx)

	response, err := q.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: types.DefaultParams()}, response)
}

func TestPermissionsQuery(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)
	q := keeper.Querier{Keeper: *k}

	wctx := sdk.WrapSDKContext(ctx)

	accAddr := utils.AccAddress()

	request := &types.QueryPermissionsRequest{
		Address: sdk.MustBech32ifyAddressBytes(sdk.Bech32PrefixAccAddr, accAddr),
	}

	response, err := q.Permissions(wctx, request)
	require.NoError(t, err)
	require.Equal(t, &types.QueryPermissionsResponse{Permissions: ""}, response)

	k.GrantPermissions(ctx, accAddr, types.NewPermissionsList([]string{"test1", "test2"}))
	response, err = q.Permissions(wctx, request)
	require.NoError(t, err)
	require.Equal(t, &types.QueryPermissionsResponse{Permissions: "test1\ntest2"}, response)
}
