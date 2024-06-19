package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s State) Validate() error {
	for _, a := range s.GetGenesisAccounts() {
		if err := a.GetAmount().Validate(); err != nil {
			return errorsmod.Wrap(err, "amount")
		}
		_, err := sdk.AccAddressFromBech32(a.GetAddress())
		if err != nil {
			return errorsmod.Wrap(err, "address from bech 32")
		}
	}
	return nil
}
