package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"

	"github.com/dymensionxyz/gerr-cosmos/gerrc"
)

func (s *State) Validate() error {
	for _, a := range s.GetGenesisAccounts() {
		if !a.Amount.IsPositive() {
			return errorsmod.Wrapf(gerrc.ErrInvalidArgument, "invalid amount: %s %s", a.Address, a.Amount)
		}

		if err := ValidateHubBech32(a.Address); err != nil {
			return errorsmod.Wrapf(err, "invalid address: %s", a.Address)
		}
	}

	if s.HubPortAndChannel != nil {
		if err := host.PortIdentifierValidator(s.HubPortAndChannel.Port); err != nil {
			return errorsmod.Wrapf(err, "invalid port Id: %s", s.HubPortAndChannel.Port)
		}
		if err := host.ChannelIdentifierValidator(s.HubPortAndChannel.Channel); err != nil {
			return errorsmod.Wrapf(err, "invalid channel Id: %s", s.HubPortAndChannel.Channel)
		}
	}

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
