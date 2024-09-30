package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

func (s *State) Validate() error {
	for _, a := range s.GetGenesisAccounts() {
		if !a.Amount.IsPositive() {
			return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "invalid amount: %s %s", a.Address, a.Amount)
		}

		if err := ValidateHubBech32(a.Address); err != nil {
			return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "invalid address: %s", a.Address)
		}
	}

	// TODO: validate port and channel?

	return nil
}

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func ValidateHubBech32(address string) (err error) {
	bech32PrefixAccAddr := "dym"
	bz, err := sdk.GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		return err
	}
	return sdk.VerifyAddressFormat(bz)
}

func (s *State) IsCanonicalHubTransferChannel(port, channel string) bool {
	return s.CanonicalHubTransferChannelHasBeenSet() && s.HubPortAndChannel.Port == port && s.HubPortAndChannel.Channel == channel
}

func (s *State) CanonicalHubTransferChannelHasBeenSet() bool {
	return s.HubPortAndChannel != nil
}

func (s *State) SetCanonicalTransferChannel(port, channel string) {
	s.HubPortAndChannel = &PortAndChannel{
		Port:    port,
		Channel: channel,
	}
}

func (g GenesisInfo) BaseDenom() string {
	return g.NativeDenom.Base
}
