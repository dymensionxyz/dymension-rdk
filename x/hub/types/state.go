package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

func (s *State) Validate() error {
	denomMap := make(map[string]struct{})
	for _, d := range s.Hub.RegisteredDenoms {
		if _, ok := denomMap[d.Base]; ok {
			return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "duplicate denom in registered denoms: %s", d)
		}
		if err := sdk.ValidateDenom(d.Base); err != nil {
			return errorsmod.Wrapf(err, "invalid denom: %s", d.Base)
		}
		denomMap[d.Base] = struct{}{}
	}

	// Validate decimal conversion pair if it exists
	if s.Hub.DecimalConversionPair != nil {
		if err := s.Hub.DecimalConversionPair.Validate(); err != nil {
			return errorsmod.Wrapf(err, "invalid decimal conversion pair")
		}
	}

	return nil
}

// Validate validates a DecimalConversionPair
func (p *DecimalConversionPair) Validate() error {
	if !strings.HasPrefix(p.FromToken, "ibc/") {
		return errorsmod.Wrap(gerrc.ErrInvalidArgument, "from_token must be an IBC denom")
	}
	if !(p.FromDecimals > 0 && p.FromDecimals < 18) {
		return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "from_decimals must be < 18, got %d", p.FromDecimals)
	}
	return nil
}
