package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// cdc.RegisterConcrete(&MsgCreateSequencer{}, "sequencer/CreateSequencer", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// registry.RegisterImplementations((*sdk.Msg)(nil),
	// 	&MsgCreateSequencer{},
	// )

	// msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
