package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	gerrc "github.com/dymensionxyz/gerr-cosmos/gerrc"
)

// ChannelState represents the state of a channel in the genesis bridge process
type ChannelState uint64

const (
	Undefined ChannelState = 0
	// WaitingForAck indicates the channel is waiting for acknowledgment
	WaitingForAck ChannelState = 1
	// Failed indicates the channel has failed and can be retried
	Failed ChannelState = 2
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

// PortAndChannel
func (p *PortAndChannel) Key() string {
	return p.Port + "/" + p.Channel
}

// FromKey
func FromPortAndChannelKey(key string) (PortAndChannel, error) {
	port, channel, found := strings.Cut(key, "/")
	if !found {
		return PortAndChannel{}, errorsmod.Wrapf(gerrc.ErrInvalidArgument, "invalid port/channel key in onGoingChannels: %s", key)
	}

	return PortAndChannel{
		Port:    port,
		Channel: channel,
	}, nil
}
