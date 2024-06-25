package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

func (s *State) Validate() error {
	denomMap := make(map[string]struct{})
	for _, d := range s.Hub.RegisteredDenoms {
		if _, ok := denomMap[d.Base]; ok {
			return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "duplicate denom in registered denoms: %s", d)
		}
		denomMap[d.Base] = struct{}{}
	}
	return nil
}
