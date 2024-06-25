package denommetadata_test

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"
	"github.com/stretchr/testify/require"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
	hubtypes "github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

func TestIBCMiddleware_OnRecvPacket(t *testing.T) {
	tests := []struct {
		name           string
		bankKeeper     *mockBankKeeper
		transferKeeper *mockTransferKeeper
		hubKeeper      *mockHubKeeper
		hooks          *mockERC20Hook

		memoData         *memoData
		wantAck          exported.Acknowledgement
		wantSentMemoData *memoData
		wantCreated      bool
	}{
		{
			name:             "valid packet data with packet metadata",
			bankKeeper:       &mockBankKeeper{},
			transferKeeper:   &mockTransferKeeper{},
			hooks:            &mockERC20Hook{},
			memoData:         validMemoData,
			wantAck:          emptyResult,
			wantSentMemoData: nil,
			wantCreated:      true,
		}, {
			name:             "valid packet data with packet metadata and user memo",
			bankKeeper:       &mockBankKeeper{},
			transferKeeper:   &mockTransferKeeper{},
			hooks:            &mockERC20Hook{},
			memoData:         validMemoDataWithUserMemo,
			wantAck:          emptyResult,
			wantSentMemoData: validUserMemo,
			wantCreated:      true,
		}, {
			name:             "no memo",
			bankKeeper:       &mockBankKeeper{},
			transferKeeper:   &mockTransferKeeper{},
			hooks:            &mockERC20Hook{},
			memoData:         nil,
			wantAck:          emptyResult,
			wantSentMemoData: nil,
			wantCreated:      false,
		}, {
			name:             "custom memo",
			bankKeeper:       &mockBankKeeper{},
			transferKeeper:   &mockTransferKeeper{},
			hooks:            &mockERC20Hook{},
			memoData:         validUserMemo,
			wantAck:          emptyResult,
			wantSentMemoData: validUserMemo,
			wantCreated:      false,
		}, {
			name:             "memo has empty denom metadata",
			bankKeeper:       &mockBankKeeper{},
			transferKeeper:   &mockTransferKeeper{},
			hooks:            &mockERC20Hook{},
			memoData:         invalidMemoDataNoDenomMetadata,
			wantAck:          emptyResult,
			wantSentMemoData: nil,
			wantCreated:      false,
		}, {
			name:             "denom metadata already exists in keeper",
			bankKeeper:       &mockBankKeeper{hasDenomMetaData: true},
			transferKeeper:   &mockTransferKeeper{},
			hooks:            &mockERC20Hook{},
			memoData:         validMemoData,
			wantAck:          emptyResult,
			wantSentMemoData: nil,
			wantCreated:      false,
		}, {
			name:             "failed to create erc20 contract",
			bankKeeper:       &mockBankKeeper{},
			transferKeeper:   &mockTransferKeeper{},
			hooks:            &mockERC20Hook{err: fmt.Errorf("failed to create erc20 contract")},
			memoData:         validMemoData,
			wantAck:          channeltypes.NewErrorAcknowledgement(fmt.Errorf("failed to create erc20 contract")),
			wantSentMemoData: nil,
			wantCreated:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &mockIBCModule{}
			im := denommetadata.NewIBCModule(
				app,
				tt.bankKeeper,
				tt.transferKeeper,
				tt.hubKeeper,
				types.NewMultiDenommetadataHooks(tt.hooks),
			)
			var memo string
			if tt.memoData != nil {
				memo = mustMarshalJSON(tt.memoData)
			}
			packetData := packetDataWithMemo(memo)
			packet := channeltypes.Packet{Data: packetData, SourcePort: "transfer", SourceChannel: "channel-0"}
			got := im.OnRecvPacket(sdk.Context{}, packet, sdk.AccAddress{})
			require.Equal(t, tt.wantAck, got)
			if !tt.wantAck.Success() {
				return
			}
			var wantMemo string
			if tt.wantSentMemoData != nil {
				wantMemo = mustMarshalJSON(tt.wantSentMemoData)
			}
			require.Equal(t, string(packetDataWithMemo(wantMemo)), string(app.sentData))
			require.Equal(t, tt.wantCreated, tt.bankKeeper.created)
		})
	}
}

func TestIBCRecvMiddleware_OnAcknowledgementPacket(t *testing.T) {
	type fields struct {
		IBCModule      porttypes.IBCModule
		bankKeeper     *mockBankKeeper
		transferKeeper *mockTransferKeeper
		hubKeeper      *mockHubKeeper
		hooks          *mockERC20Hook
	}
	type args struct {
		ctx             sdk.Context
		packet          channeltypes.Packet
		acknowledgement []byte
		relayer         sdk.AccAddress
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid packet data with packet metadata",
			fields: fields{
				IBCModule:      &mockIBCModule{},
				bankKeeper:     &mockBankKeeper{},
				transferKeeper: &mockTransferKeeper{},
				hooks:          &mockERC20Hook{},
				hubKeeper: &mockHubKeeper{
					hub: hubtypes.Hub{},
				},
			},
			args: args{
				ctx:             sdk.Context{},
				packet:          channeltypes.Packet{Data: packetDataWithMemo(mustMarshalJSON(validMemoData)), SourcePort: "transfer", SourceChannel: "channel-0"},
				acknowledgement: emptyResult.Acknowledgement(),
				relayer:         sdk.AccAddress{},
			},
			wantErr: false,
		}, {
			name: "valid packet data with packet metadata and user memo",
			fields: fields{
				IBCModule:      &mockIBCModule{},
				bankKeeper:     &mockBankKeeper{},
				transferKeeper: &mockTransferKeeper{},
				hooks:          &mockERC20Hook{},
				hubKeeper:      &mockHubKeeper{},
			},
			args: args{
				ctx:             sdk.Context{},
				packet:          channeltypes.Packet{Data: packetDataWithMemo(mustMarshalJSON(validMemoDataWithUserMemo)), SourcePort: "transfer", SourceChannel: "channel-0"},
				acknowledgement: emptyResult.Acknowledgement(),
				relayer:         sdk.AccAddress{},
			},
			wantErr: false,
		}, {
			name: "no memo",
			fields: fields{
				IBCModule:      &mockIBCModule{},
				bankKeeper:     &mockBankKeeper{},
				transferKeeper: &mockTransferKeeper{},
				hooks:          &mockERC20Hook{},
				hubKeeper:      &mockHubKeeper{},
			},
			args: args{
				ctx:             sdk.Context{},
				packet:          channeltypes.Packet{Data: packetDataWithMemo(""), SourcePort: "transfer", SourceChannel: "channel-0"},
				acknowledgement: emptyResult.Acknowledgement(),
				relayer:         sdk.AccAddress{},
			},
			wantErr: false,
		}, {
			name: "custom memo",
			fields: fields{
				IBCModule:      &mockIBCModule{},
				bankKeeper:     &mockBankKeeper{},
				transferKeeper: &mockTransferKeeper{},
				hooks:          &mockERC20Hook{},
				hubKeeper:      &mockHubKeeper{},
			},
			args: args{
				ctx:             sdk.Context{},
				packet:          channeltypes.Packet{Data: packetDataWithMemo(mustMarshalJSON(validUserMemo)), SourcePort: "transfer", SourceChannel: "channel-0"},
				acknowledgement: emptyResult.Acknowledgement(),
				relayer:         sdk.AccAddress{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := denommetadata.NewIBCModule(
				tt.fields.IBCModule,
				tt.fields.bankKeeper,
				tt.fields.transferKeeper,
				tt.fields.hubKeeper,
				types.NewMultiDenommetadataHooks(tt.fields.hooks),
			)
			err := im.OnAcknowledgementPacket(tt.args.ctx, tt.args.packet, tt.args.acknowledgement, tt.args.relayer)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

var (
	emptyResult   = channeltypes.Acknowledgement{}
	validUserMemo = &memoData{
		User: &validUserData,
	}
	validMemoDataWithUserMemo = &memoData{
		MemoData: validMemoData.MemoData,
		User:     &validUserData,
	}
	validUserData = userData{Data: "data"}
	validMemoData = &memoData{
		MemoData: types.MemoData{
			DenomMetadata: &validDenomMetadata,
		},
	}
	invalidMemoDataNoDenomMetadata = &memoData{
		MemoData: types.MemoData{},
	}
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

type memoData struct {
	types.MemoData
	User *userData `json:"user,omitempty"`
}

type userData struct {
	Data string `json:"data"`
}

func packetDataWithMemo(memo string) []byte {
	byt, _ := types.ModuleCdc.MarshalJSON(&transfertypes.FungibleTokenPacketData{
		Denom:    "adym",
		Amount:   "100",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     memo,
	})
	return byt
}

func mustMarshalJSON(v any) string {
	bz, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bz)
}

type mockIBCModule struct {
	porttypes.IBCModule
	sentData []byte
}

func (m *mockIBCModule) OnRecvPacket(_ sdk.Context, p channeltypes.Packet, _ sdk.AccAddress) exported.Acknowledgement {
	m.sentData = p.Data
	return emptyResult
}

func (m *mockIBCModule) OnAcknowledgementPacket(sdk.Context, channeltypes.Packet, []byte, sdk.AccAddress) error {
	return nil
}

type mockBankKeeper struct {
	hasDenomMetaData, created bool
}

func (m *mockBankKeeper) SetDenomMetaData(sdk.Context, banktypes.Metadata) {
	m.created = true
}

func (m *mockBankKeeper) GetDenomMetaData(sdk.Context, string) (banktypes.Metadata, bool) {
	return banktypes.Metadata{}, m.hasDenomMetaData
}

type mockTransferKeeper struct {
	hasDT   bool
	created bool
}

func (m *mockTransferKeeper) HasDenomTrace(sdk.Context, tmbytes.HexBytes) bool {
	return m.hasDT
}

func (m *mockTransferKeeper) SetDenomTrace(sdk.Context, transfertypes.DenomTrace) {
	m.created = true
}

func (m *mockTransferKeeper) OnRecvPacket(sdk.Context, channeltypes.Packet, sdk.AccAddress) exported.Acknowledgement {
	return emptyResult
}

type mockHubKeeper struct {
	hub hubtypes.Hub
}

func (m *mockHubKeeper) SetState(ctx sdk.Context, state hubtypes.State) {
	m.hub = state.Hub
}

func (m *mockHubKeeper) GetState(ctx sdk.Context) hubtypes.State {
	return hubtypes.State{Hub: m.hub}
}

type mockERC20Hook struct {
	createCalled bool
	err          error
	sync.Mutex
}

func (m *mockERC20Hook) AfterDenomMetadataCreation(sdk.Context, banktypes.Metadata) error {
	m.Lock()
	defer m.Unlock()
	m.createCalled = m.err == nil
	return m.err
}

func (m *mockERC20Hook) AfterDenomMetadataUpdate(sdk.Context, banktypes.Metadata) error { return nil }
