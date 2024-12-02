package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func ValidateHubBech32(address string) (err error) {
	bech32PrefixAccAddr := "dym"
	bz, err := sdk.GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		return err
	}
	return sdk.VerifyAddressFormat(bz)
}

func (a GenesisAccount) ValidateBasic() error {
	if !a.Amount.IsPositive() {
		return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "invalid amount: %s %s", a.Address, a.Amount)
	}

	if err := ValidateHubBech32(a.Address); err != nil {
		return errorsmod.Wrapf(err, "invalid address: %s", a.Address)
	}

	return nil
}

func (g GenesisInfo) BaseDenom() string {
	return g.NativeDenom.Base
}

// BaseCoinSupply returns the total supply of the base denom: the sum of all the genesis account amounts.
func (g GenesisInfo) BaseCoinSupply() sdk.Coin {
	amount := sdk.ZeroInt()
	for _, acc := range g.GenesisAccounts {
		amount = amount.Add(acc.Amount)
	}
	return sdk.NewCoin(g.BaseDenom(), amount)
}
