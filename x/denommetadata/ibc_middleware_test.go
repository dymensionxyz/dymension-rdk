package denommetadata_test

import (
	"encoding/json"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"
	"github.com/stretchr/testify/require"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

func TestIBCMiddleware_OnRecvPacket(t *testing.T) {
	tests := []struct {
		name         string
		keeper       *mockKeeper
		packet       channeltypes.Packet
		wantAck      exported.Acknowledgement
		wantSentData []byte
		wantCreated  bool
	}{
		{
			name:   "valid packet data with packet metadata",
			keeper: &mockKeeper{},
			packet: channeltypes.Packet{
				Data: mustMarshalJSON(validFungibleTokenWithPacketMetadata),
			},
			wantAck:      emptyResult,
			wantSentData: mustMarshalJSON(validFungibleToken),
			wantCreated:  true,
		}, {
			name:   "invalid packet data",
			keeper: &mockKeeper{},
			packet: channeltypes.Packet{
				Data: []byte(``),
			},
			wantAck:      channeltypes.NewErrorAcknowledgement(errortypes.ErrInvalidType),
			wantSentData: nil,
			wantCreated:  false,
		}, {
			name:   "no memo",
			keeper: &mockKeeper{},
			packet: channeltypes.Packet{
				Data: []byte(`{"memo": ""}`),
			},
			wantAck:      emptyResult,
			wantSentData: []byte(`{"memo": ""}`),
			wantCreated:  false,
		}, {
			name:   "custom memo",
			keeper: &mockKeeper{},
			packet: channeltypes.Packet{
				Data: []byte(`{"memo": "thanks for the sweater, grandma!"}`),
			},
			wantAck:      emptyResult,
			wantSentData: []byte(`{"memo": "thanks for the sweater, grandma!"}`),
			wantCreated:  false,
		}, {
			name:   "memo has empty packet metadata",
			keeper: &mockKeeper{},
			packet: channeltypes.Packet{
				Data: []byte(`{"memo": "{\"packet_metadata\":\"\"}"}`),
			},
			wantAck:      emptyResult,
			wantSentData: []byte(`{"memo": "{\"packet_metadata\":\"\"}"}`),
			wantCreated:  false,
		}, {
			name:   "memo has empty denom metadata",
			keeper: &mockKeeper{},
			packet: channeltypes.Packet{
				Data: []byte(`{"memo": "{\"packet_metadata\":{\"denom_metadata\":null}}"}`),
			},
			wantAck:      emptyResult,
			wantSentData: []byte(`{"memo": "{\"packet_metadata\":{\"denom_metadata\":null}}"}`),
			wantCreated:  false,
		}, {
			name:   "denom metadata already exists in keeper",
			keeper: &mockKeeper{hasDenomMetaData: true},
			packet: channeltypes.Packet{
				Data: mustMarshalJSON(validFungibleTokenWithPacketMetadata),
			},
			wantAck:      emptyResult,
			wantSentData: mustMarshalJSON(validFungibleToken),
			wantCreated:  false,
		}, {
			name:   "keeper error",
			keeper: &mockKeeper{err: fmt.Errorf("whatever")},
			packet: channeltypes.Packet{
				Data: mustMarshalJSON(validFungibleTokenWithPacketMetadata),
			},
			wantAck:      channeltypes.NewErrorAcknowledgement(fmt.Errorf("whatever")),
			wantSentData: nil,
			wantCreated:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &mockIBCModule{}
			im := denommetadata.NewIBCMiddleware(tt.keeper, app)
			got := im.OnRecvPacket(sdk.Context{}, tt.packet, sdk.AccAddress{})
			require.Equal(t, tt.wantAck, got)
			require.Equal(t, string(tt.wantSentData), string(app.sentData))
			require.Equal(t, tt.wantCreated, tt.keeper.created)
		})
	}
}

var (
	emptyResult = channeltypes.Acknowledgement{}

	validFungibleToken = transfertypes.FungibleTokenPacketData{
		Denom:    "adym",
		Amount:   "100",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     validUserMemo,
	}
	validFungibleTokenWithPacketMetadata = transfertypes.FungibleTokenPacketData{
		Denom:    "adym",
		Amount:   "100",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     string(mustMarshalJSON(validMemoData)),
	}
	validMemoData = &types.MemoData{
		UserMemo: validUserMemo,
		PacketMetadata: &types.PacketMetadata{
			DenomMetadata: &validDenomMetadata,
		},
	}
	validUserMemo      = "user memo"
	validDenomMetadata = banktypes.Metadata{
		Description: "Denom of the Hub",
		Base:        "adym",
		Display:     "DYM",
		Name:        "DYM",
		Symbol:      "adym",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "adym",
				Exponent: 0,
			}, {
				Denom:    "DYM",
				Exponent: 18,
			},
		},
	}
)

func mustMarshalJSON(v any) []byte {
	bz, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bz
}

type mockIBCModule struct {
	porttypes.IBCModule
	sentData []byte
}

func (m *mockIBCModule) OnRecvPacket(_ sdk.Context, p channeltypes.Packet, _ sdk.AccAddress) exported.Acknowledgement {
	m.sentData = p.Data
	return emptyResult
}

type mockKeeper struct {
	hasDenomMetaData, created bool
	err                       error
}

func (m *mockKeeper) CreateDenomMetadata(sdk.Context, ...types.DenomMetadata) error {
	m.created = m.err == nil
	return m.err
}

func (m mockKeeper) HasDenomMetaData(sdk.Context, string) bool { return m.hasDenomMetaData }
