package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

var (
	_ sdk.Msg = (*MsgSendTransfer)(nil)
)

const (
	TypeMsgSendTransfer = "send_transfer"
)

var (
	_ sdk.Msg = &MsgSendTransfer{}
	//_ legacytx.LegacyMsg = &MsgSendTransfer{}
)

//func RegisterCodec(cdc *codec.LegacyAmino) {
//	cdc.RegisterConcrete(&MsgSetCanonicalClient{}, "lightclient/SetCanonicalClient", nil)
//}

//func (msg *MsgSetCanonicalClient) Route() string {
//	return ModuleName
//}
//
//func (msg *MsgSetCanonicalClient) Type() string {
//	return MsgSendTransfer
//}

func (m *MsgSendTransfer) GetSigners() []sdk.AccAddress {
	a, _ := sdk.AccAddressFromBech32(m.Relayer)
	return []sdk.AccAddress{a}
}

func (m *MsgSendTransfer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Relayer)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get relayer addr from bech32")
	}
	if m.ChannelId == "" {
		return errorsmod.Wrap(gerrc.ErrInvalidArgument, "channel id is empty")
	}
	return nil
}
