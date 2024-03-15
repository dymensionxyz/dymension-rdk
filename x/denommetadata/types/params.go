package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultEnableDrip       = true
	DefaultAllowedAddresses = []string(nil) // no one allowed
)

// NewParams creates a new Params object
func NewParams(
	allowedAddresses []string,
) Params {
	return Params{
		AllowedAddresses: allowedAddresses,
	}
}

// DefaultParams returns default x/drip module parameters.
func DefaultParams() Params {
	return Params{
		AllowedAddresses: DefaultAllowedAddresses,
	}
}

func validateArray(i interface{}) error {
	_, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func (p Params) Validate() error {
	if err := validateArray(p.AllowedAddresses); err != nil {
		return err
	}

	return assertValidAddresses(p.AllowedAddresses)
}

func assertValidAddresses(addrs []string) error {
	idx := make(map[string]struct{}, len(addrs))
	for _, a := range addrs {
		if a == "" {
			return ErrBlank.Wrapf("address: %s", a)
		}
		if _, err := sdk.AccAddressFromBech32(a); err != nil {
			return errorsmod.Wrapf(err, "address: %s", a)
		}
		if _, exists := idx[a]; exists {
			return ErrDuplicate.Wrapf("address: %s", a)
		}
		idx[a] = struct{}{}
	}
	return nil
}
