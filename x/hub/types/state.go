package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

func (s *State) Validate() error {
	denomMap := make(map[string]struct{})
	for base := range s.Hub.RegisteredDenoms {
		if _, ok := denomMap[base]; ok {
			return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "duplicate denom in registered denoms: %s", base)
		}
		denomMap[base] = struct{}{}
	}
	return nil
}
