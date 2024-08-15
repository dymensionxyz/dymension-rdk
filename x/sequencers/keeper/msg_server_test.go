package keeper

import (
	"context"
	"fmt"
	"testing"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/assert"
)

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
