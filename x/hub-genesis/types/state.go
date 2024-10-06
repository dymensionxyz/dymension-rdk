package types

import (
	errorsmod "cosmossdk.io/errors"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

func (s *State) Validate() error {
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
