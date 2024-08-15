package keeper

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
}

func Test_checkSigAccNumber(t *testing.T) {
	type args struct {
		ctx        sdk.Context
		acc        uint64
		keyAndSig  *types.KeyAndSig
		payloadApp codec.ProtoMarshaler
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkSigAccNumber(tt.args.ctx, tt.args.acc, tt.args.keyAndSig, tt.args.payloadApp)
			if !tt.wantErr(t, err, fmt.Sprintf("checkSigAccNumber(%v, %v, %v, %v)", tt.args.ctx, tt.args.acc, tt.args.keyAndSig, tt.args.payloadApp)) {
				return
			}
			assert.Equalf(t, tt.want, got, "checkSigAccNumber(%v, %v, %v, %v)", tt.args.ctx, tt.args.acc, tt.args.keyAndSig, tt.args.payloadApp)
		})
	}
}
