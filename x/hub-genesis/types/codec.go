package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
}

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

var (
	amino     = codec.NewLegacyAmino()
	moduleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
