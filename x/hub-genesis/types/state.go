package types

import (
	"reflect"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *State) Validate() error {
	for _, a := range s.GetGenesisAccounts() {
		if err := a.GetAmount().Validate(); err != nil {
			return errorsmod.Wrap(err, "amount")
		}
		// TODO: validate to make sure no ibc denom attack possible https://github.com/dymensionxyz/dymension-rdk/issues/462
		_, err := sdk.AccAddressFromBech32(a.GetAddress())
		if err != nil {
			return errorsmod.Wrap(err, "address from bech 32")
		}
	}
	return nil
}

func (s *State) IsCanonicalHubTransferChannel(port, channel string) bool {
	return s.CanonicalHubTransferChannelHasBeenSet() && s.HubPortAndChannel.Port == port && s.HubPortAndChannel.Channel == channel
}

func (s *State) CanonicalHubTransferChannelHasBeenSet() bool {
	return !reflect.ValueOf(s.HubPortAndChannel).IsZero()
}

func (s *State) SetCanonicalTransferChannel(port, channel string) {
	s.HubPortAndChannel = &PortAndChannel{
		Port:    port,
		Channel: channel,
	}
}
