package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/rollapp/testutil/utils"
	"github.com/dymensionxyz/rollapp/x/sequencers/keeper"
	"github.com/dymensionxyz/rollapp/x/sequencers/testutils"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestParamsQuery(t *testing.T) {
	app := utils.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	k := keeper.Querier{Keeper: *testutils.NewTestSequencer(ctx)}

	wctx := sdk.WrapSDKContext(ctx)

	response, err := k.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: types.DefaultParams()}, response)
}
