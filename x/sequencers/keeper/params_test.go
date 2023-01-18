package keeper_test

import (
	"testing"

	"github.com/dymensionxyz/rollapp/testutil/utils"
	"github.com/dymensionxyz/rollapp/x/sequencers/testutils"
	"github.com/dymensionxyz/rollapp/x/sequencers/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestGetParams(t *testing.T) {
	app := utils.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	k := testutils.NewTestSequencer(ctx)

	//TODO: change params and validate
	params := types.DefaultParams()
	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
