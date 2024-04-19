package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgHubGenesisEvent{}, "hub-genesis/HubGenesisEvent", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgHubGenesisEvent{})
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

var (
	amino     = codec.NewLegacyAmino()
	moduleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
