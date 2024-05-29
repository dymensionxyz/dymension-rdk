package denommetadata_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"
	"github.com/stretchr/testify/require"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata"
	denommetadatamoduletypes "github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

func TestIBCMiddleware_OnRecvPacket(t *testing.T) {
	tests := []struct {
		name        string
		keeper      *mockKeeper
		packet      types.Packet
		want        exported.Acknowledgement
		wantCreated bool
	}{
		{
			name:   "valid packet data",
			keeper: &mockKeeper{},
			packet: types.Packet{
				Data: validDenomMetadata,
			},
			want:        emptyResult,
			wantCreated: true,
		}, {
			name:   "invalid packet data",
			keeper: &mockKeeper{},
			packet: types.Packet{
				Data: []byte(``),
			},
			want:        types.NewErrorAcknowledgement(errortypes.ErrInvalidType),
			wantCreated: false,
		}, {
			name:   "no memo",
			keeper: &mockKeeper{},
			packet: types.Packet{
				Data: []byte(`{"memo": ""}`),
			},
			want:        emptyResult,
			wantCreated: false,
		}, {
			name:   "bad memo",
			keeper: &mockKeeper{},
			packet: types.Packet{
				Data: []byte(`{"memo": "bad"}`),
			},
			want:        emptyResult,
			wantCreated: false,
		}, {
			name:   "memo has invalid metadata",
			keeper: &mockKeeper{},
			packet: types.Packet{
				Data: []byte(`{"memo": "{\"denom_metadata\":\"invalid\"}"}`),
			},
			want:        types.NewErrorAcknowledgement(fmt.Errorf("whatever")),
			wantCreated: false,
		}, {
			name:   "denom metadata already exists in keeper",
			keeper: &mockKeeper{hasDenomMetaData: true},
			packet: types.Packet{
				Data: validDenomMetadata,
			},
			want:        emptyResult,
			wantCreated: false,
		}, {
			name:   "keeper error",
			keeper: &mockKeeper{err: fmt.Errorf("whatever")},
			packet: types.Packet{
				Data: validDenomMetadata,
			},
			want:        types.NewErrorAcknowledgement(fmt.Errorf("whatever")),
			wantCreated: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := denommetadata.NewIBCMiddleware(tt.keeper, mockIBCModule{})
			got := im.OnRecvPacket(sdk.Context{}, tt.packet, sdk.AccAddress{})
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantCreated, tt.keeper.created)
		})
	}
}

var (
	emptyResult        = types.Acknowledgement{}
	validDenomMetadata = []byte(`{"memo": "{\"denom_metadata\":{\"name\":\"name\",\"symbol\":\"symbol\",\"base\":\"base\",\"display\":\"display\",\"denom_units\":[{\"denom\":\"base\",\"exponent\":0},{\"denom\":\"display\",\"exponent\":18}]}}"}`)
)

type mockIBCModule struct {
	porttypes.IBCModule
}

func (m mockIBCModule) OnRecvPacket(sdk.Context, types.Packet, sdk.AccAddress) exported.Acknowledgement {
	return emptyResult
}

type mockKeeper struct {
	hasDenomMetaData, created bool
	err                       error
}

func (m *mockKeeper) CreateDenomMetadata(sdk.Context, ...denommetadatamoduletypes.DenomMetadata) error {
	m.created = m.err == nil
	return m.err
}

func (m mockKeeper) HasDenomMetaData(sdk.Context, string) bool { return m.hasDenomMetaData }
