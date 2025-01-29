package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/staking interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// custom register

	legacy.RegisterAminoMsg(cdc, &MsgCreateValidatorERC20{}, "dymension/MsgCreateValidatorERC20")
	legacy.RegisterAminoMsg(cdc, &MsgDelegateERC20{}, "dymension/MsgDelegateERC20")
}

// RegisterInterfaces registers the x/staking interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateValidatorERC20{},
		&MsgDelegateERC20{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// var (
// 	amino     = codec.NewLegacyAmino()
// 	ModuleCdc = codec.NewAminoCodec(amino)
// )

// // func init() {
// // 	RegisterLegacyAminoCodec(amino)
// // 	cryptocodec.RegisterCrypto(amino)
// // 	sdk.RegisterLegacyAminoCodec(amino)

// // 	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
// // 	// used to properly serialize MsgGrant and MsgExec instances
// // 	RegisterLegacyAminoCodec(authzcodec.Amino)
// // }
