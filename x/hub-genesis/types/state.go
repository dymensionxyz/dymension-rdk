package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/utils"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

func (s State) Validate() error {
	for _, a := range s.GetGenesisAccounts() {
		if err := a.GetAmount().Validate(); err != nil {
			return errorsmod.Wrap(err, "amount")
		}
		if utils.IsIBCDenom(a.Amount.GetDenom()) {
			return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "ibc denoms not allowed in genesis accounts: %s", a.Amount)
		}
		_, err := sdk.AccAddressFromBech32(a.GetAddress())
		if err != nil {
			return errorsmod.Wrap(err, "address from bech 32")
		}
	}
	return nil
}
