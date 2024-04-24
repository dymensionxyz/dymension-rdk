package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
)

// RegisterLegacyAminoCodec registers the necessary x/staking interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateGovernor{}, "governors/MsgCreateGovernor")
	legacy.RegisterAminoMsg(cdc, &MsgEditGovernor{}, "governors/MsgEditGovernor")
	legacy.RegisterAminoMsg(cdc, &MsgDelegate{}, "governors/MsgDelegate")
	legacy.RegisterAminoMsg(cdc, &MsgUndelegate{}, "governors/MsgUndelegate")
	legacy.RegisterAminoMsg(cdc, &MsgBeginRedelegate{}, "governors/MsgBeginRedelegate")
	legacy.RegisterAminoMsg(cdc, &MsgCancelUnbondingDelegation{}, "governors/MsgCancelUnbondingDelegation")

	cdc.RegisterInterface((*isStakeAuthorization_Governors)(nil), nil)
	cdc.RegisterConcrete(&StakeAuthorization_AllowList{}, "governors/StakeAuthorization/AllowList", nil)
	cdc.RegisterConcrete(&StakeAuthorization_DenyList{}, "governors/StakeAuthorization/DenyList", nil)
	cdc.RegisterConcrete(&StakeAuthorization{}, "governors/StakeAuthorization", nil)
}

// RegisterInterfaces registers the x/staking interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateGovernor{},
		&MsgEditGovernor{},
		&MsgDelegate{},
		&MsgUndelegate{},
		&MsgBeginRedelegate{},
		&MsgCancelUnbondingDelegation{},
	)
	registry.RegisterImplementations(
		(*authz.Authorization)(nil),
		&StakeAuthorization{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
}
