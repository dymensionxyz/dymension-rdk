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

func (im IBCMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	var altPacket types.WrappedFungibleTokenPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &altPacket); err != nil {
		// try if it's a FungibleTokenPacketData packet
		altPacket.FungibleTokenPacketData = new(transfertypes.FungibleTokenPacketData)
		if nerr := types.ModuleCdc.UnmarshalJSON(packet.GetData(), altPacket.FungibleTokenPacketData); nerr == nil {
			return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
		}
		err = errorsmod.Wrapf(errortypes.ErrInvalidType, "cannot unmarshal ICS-20 transfer packet data")
		return channeltypes.NewErrorAcknowledgement(err)
	}

	rawDenom := altPacket.FungibleTokenPacketData.Denom

	denomTrace := transfertypes.ParseDenomTrace(rawDenom)
	if denomTrace.Path == "" {
		denomTrace.Path = fmt.Sprintf("%s/%s", packet.GetDestPort(), packet.GetDestChannel())
	}

	_, err := im.keeper.GetDenomMetadata(ctx, denomTrace)
	if err != nil {
		if errors.Is(err, banktypes.ErrDenomMetadataNotFound) {
			if err = im.createNewDenom(ctx, altPacket.DenomMetadata, denomTrace); err != nil {
				return channeltypes.NewErrorAcknowledgement(err)
			}
		} else {
			return channeltypes.NewErrorAcknowledgement(err)
		}
	}

	packetData := altPacket.FungibleTokenPacketData

	packet.Data, err = types.ModuleCdc.MarshalJSON(packetData)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
}

func (im IBCMiddleware) createNewDenom(ctx sdk.Context, dm *banktypes.Metadata, denomTrace transfertypes.DenomTrace) error {
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
