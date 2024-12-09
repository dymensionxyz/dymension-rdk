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

func (m *MsgSendTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Relayer)}
}

func (m *MsgSendTransfer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Relayer)
	if err != nil {
		return errorsmod.Wrap(errors.Join(gerrc.ErrInvalidArgument, err), "get relayer addr from bech32")
	}
	return nil
}
