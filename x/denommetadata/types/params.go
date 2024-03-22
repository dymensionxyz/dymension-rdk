package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	KeyAllowedAddresses     = []byte("AllowedAddresses")
	DefaultAllowedAddresses = []string(nil) // no one allowed
)

// ParamKeyTable for minting module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	allowedAddresses []string,
) Params {
	return Params{
		AllowedAddresses: allowedAddresses,
	}
}

// DefaultParams returns default x/denommetadata module parameters.
func DefaultParams() Params {
	return Params{
		AllowedAddresses: DefaultAllowedAddresses,
	}
}

func (p Params) Validate() error {
	return assertValidAddresses(p.AllowedAddresses)
}

func assertValidAddresses(i any) error {
	addrs, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	idx := make(map[string]struct{}, len(addrs))
	for _, a := range addrs {

		// this also checks for empty addresses
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

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAllowedAddresses, &p.AllowedAddresses, assertValidAddresses),
	}
}
