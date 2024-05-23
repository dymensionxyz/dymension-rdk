// Copyright 2022 Evmos Foundation
// This file is part of the Evmos Network packages.
//
// Evmos is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Evmos packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Evmos packages. If not, see https://github.com/evmos/evmos/blob/main/LICENSE

package denommetadata

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

var _ porttypes.IBCModule = &IBCMiddleware{}

// IBCMiddleware implements the ICS26 callbacks for the transfer middleware
type IBCMiddleware struct {
	porttypes.IBCModule
	keeper        keeper.Keeper
	channelKeeper types.ChannelKeeper
}

// NewIBCMiddleware creates a new IBCMiddleware given the keeper and underlying application
func NewIBCMiddleware(k keeper.Keeper, channelKeeper types.ChannelKeeper, app porttypes.IBCModule) IBCMiddleware {
	return IBCMiddleware{
		channelKeeper: channelKeeper,
		IBCModule:     app,
		keeper:        k,
	}
}

// OnRecvPacket registers the denom metadata if it does not exist
func (im IBCMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	/*_, clientState, err := im.channelKeeper.GetChannelClientState(ctx, packet.DestinationPort, packet.DestinationChannel)
	if err != nil {
		err = errorsmod.Wrapf(errortypes.ErrInvalidRequest, "client state not found")
		return channeltypes.NewErrorAcknowledgement(err)
	}

	// Extract the chain ID from the client state
	tmClientState, ok := clientState.(*tenderminttypes.ClientState)
	if !ok {
		return channeltypes.NewErrorAcknowledgement(errors.New("expected tendermint client state"))
	}

	sourceChainID := tmClientState.GetChainID()
	_ = sourceChainID*/
	// TODO: check source chain against a whitelist or something

	packetData := new(transfertypes.FungibleTokenPacketData)
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), packetData); err != nil {
		err = errorsmod.Wrapf(errortypes.ErrInvalidType, "cannot unmarshal ICS-20 transfer packet data")
		return channeltypes.NewErrorAcknowledgement(err)
	}

	if len(packetData.Memo) == 0 {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	rawDenom := packetData.Denom

	denomTrace := transfertypes.ParseDenomTrace(rawDenom)
	if denomTrace.Path == "" {
		denomTrace.Path = fmt.Sprintf("%s/%s", packet.GetDestPort(), packet.GetDestChannel())
	}

	if im.keeper.HasDenomMetadata(ctx, denomTrace) {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	if err := im.createNewDenom(ctx, packetData, denomTrace); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
}

func (im IBCMiddleware) createNewDenom(ctx sdk.Context, packetData *transfertypes.FungibleTokenPacketData, denomTrace transfertypes.DenomTrace) error {
	packetMetaData, err := types.ParsePacketMetadata(packetData.Memo)
	if errors.Is(err, types.ErrMemoUnmarshal) || errors.Is(err, types.ErrMemoDMEmpty) {
		return nil
	}
	if err != nil {
		return err
	}

	dm := packetMetaData.DenomMetadata

	denomUnits := make([]*banktypes.DenomUnit, 0, len(dm.DenomUnits))
	for _, du := range dm.DenomUnits {
		if du.Exponent == 0 {
			continue
		}
		ndu := &banktypes.DenomUnit{
			Denom:    du.Denom,
			Exponent: du.Exponent,
			Aliases:  du.Aliases,
		}
		denomUnits = append(denomUnits, ndu)
	}

	newDenomMetadata := types.DenomMetadata{
		TokenMetadata: banktypes.Metadata{
			Description: dm.Description,
			DenomUnits:  denomUnits,
			Base:        denomTrace.IBCDenom(),
			Display:     dm.Display,
			Name:        dm.Name,
			Symbol:      dm.Symbol,
			URI:         dm.URI,
			URIHash:     dm.URIHash,
		},
		DenomTrace: denomTrace.GetFullDenomPath(),
	}

	if err := im.keeper.CreateDenomMetadata(ctx, newDenomMetadata); err != nil {
		return fmt.Errorf("failed to create denom metadata: %w", err)
	}

	return nil
}
