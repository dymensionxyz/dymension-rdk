package keeper

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeepers "github.com/dymensionxyz/dymension-rdk/testutil/keepers"
	"github.com/dymensionxyz/dymension-rdk/testutil/utils"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUpdateHappyPath(t *testing.T) {
	app := utils.Setup(t, false)
	k, ctx := testkeepers.NewTestSequencerKeeperFromApp(app)

	msgServer := msgServer{*k}

	wctx := sdk.WrapSDKContext(ctx)

	_, err := msgServer.CreateSequencer(wctx, &types.MsgCreateSequencer{})
	require.NoError(t, err)

	_, err = msgServer.UpdateSequencer(wctx, &types.MsgUpdateSequencer{})
	require.NoError(t, err)
}

func Test_msgServer_CreateSequencer(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		goCtx context.Context
		msg   *types.MsgCreateSequencer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.MsgCreateSequencerResponse
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := msgServer{
				Keeper: tt.fields.Keeper,
			}
			got, err := m.CreateSequencer(tt.args.goCtx, tt.args.msg)
			if !tt.wantErr(t, err, fmt.Sprintf("CreateSequencer(%v, %v)", tt.args.goCtx, tt.args.msg)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CreateSequencer(%v, %v)", tt.args.goCtx, tt.args.msg)
		})
	}
}
