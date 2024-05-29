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

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

var _ porttypes.IBCModule = &IBCMiddleware{}

// IBCMiddleware implements the ICS26 callbacks for the transfer middleware
type IBCMiddleware struct {
	porttypes.IBCModule
	keeper DenomMetadataKeeper
}

type DenomMetadataKeeper interface {
	HasDenomMetaData(ctx sdk.Context, denom string) bool
	CreateDenomMetadata(ctx sdk.Context, metadatas ...types.DenomMetadata) error
}

// NewIBCMiddleware creates a new IBCMiddleware given the keeper and underlying application
func NewIBCMiddleware(keeper DenomMetadataKeeper, app porttypes.IBCModule) IBCMiddleware {
	return IBCMiddleware{
		IBCModule: app,
		keeper:    keeper,
	}
}

// OnRecvPacket registers the denom metadata if it does not exist
func (im IBCMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	packetData := new(transfertypes.FungibleTokenPacketData)
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), packetData); err != nil {
		err = errorsmod.Wrapf(errortypes.ErrInvalidType, "cannot unmarshal ICS-20 transfer packet data")
		return channeltypes.NewErrorAcknowledgement(err)
	}

	if len(packetData.Memo) == 0 {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	packetMetaData, err := types.ParsePacketMetadata(packetData.Memo)
	// if the memo is not an object, or does not contain the metadata, we can skip
	if errors.Is(err, types.ErrMemoUnmarshal) || errors.Is(err, types.ErrMemoDenomMetadataEmpty) {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	if err = packetMetaData.DenomMetadata.Validate(); err != nil {
		return channeltypes.NewErrorAcknowledgement(errortypes.ErrInvalidType)
	}

	denomTrace := transfertypes.ParseDenomTrace(packetData.Denom)
	if denomTrace.Path == "" {
		denomTrace.Path = fmt.Sprintf("%s/%s", packet.GetDestPort(), packet.GetDestChannel())
	}

	if im.keeper.HasDenomMetaData(ctx, denomTrace.IBCDenom()) {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	if err := im.createNewDenom(ctx, packetMetaData.DenomMetadata, denomTrace); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
}

func (im IBCMiddleware) createNewDenom(ctx sdk.Context, denonmMetadata banktypes.Metadata, denomTrace transfertypes.DenomTrace) error {
	denomUnits := make([]*banktypes.DenomUnit, 0, len(denonmMetadata.DenomUnits))
	for _, du := range denonmMetadata.DenomUnits {
		// we can skip the exp 0, it's not very useful
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
			Description: denonmMetadata.Description,
			DenomUnits:  denomUnits,
			Base:        denomTrace.IBCDenom(),
			Display:     denonmMetadata.Display,
			Name:        denonmMetadata.Name,
			Symbol:      denonmMetadata.Symbol,
			URI:         denonmMetadata.URI,
			URIHash:     denonmMetadata.URIHash,
		},
		DenomTrace: denomTrace.GetFullDenomPath(),
	}

	if err := im.keeper.CreateDenomMetadata(ctx, newDenomMetadata); err != nil {
		return fmt.Errorf("create denom metadata: %w", err)
	}

	return nil
}
