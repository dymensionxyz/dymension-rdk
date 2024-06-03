package denommetadata

import (
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
	keeper Keeper
}

type Keeper interface {
	HasDenomMetaData(ctx sdk.Context, denom string) bool
	CreateDenomMetadata(ctx sdk.Context, metadatas ...types.DenomMetadata) error
}

// NewIBCMiddleware creates a new IBCMiddleware given the keeper and underlying application
func NewIBCMiddleware(keeper Keeper, app porttypes.IBCModule) IBCMiddleware {
	return IBCMiddleware{
		IBCModule: app,
		keeper:    keeper,
	}
}

// OnRecvPacket registers the denom metadata if it does not exist.
// It will intercept an incoming packet and check if the denom metadata exists.
// If it does not, it will register the denom metadata.
// The handler will expect a 'transferinject' object in the memo field of the packet.
// If the memo is not an object, or does not contain the metadata, it moves on to the next handler.
func (im IBCMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	packetData := new(transfertypes.FungibleTokenPacketData)
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), packetData); err != nil {
		err = errorsmod.Wrapf(errortypes.ErrJSONUnmarshal, "unmarshal ICS-20 transfer packet data")
		return channeltypes.NewErrorAcknowledgement(err)
	}

	transferInject, err := types.ParsePacketMetadata(packetData.Memo)
	if err != nil {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	dm := transferInject.DenomMetadata
	if dm == nil {
		return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
	}

	if err = dm.Validate(); err != nil {
		return channeltypes.NewErrorAcknowledgement(errortypes.ErrInvalidType)
	}

	denomTrace := transfertypes.ParseDenomTrace(packetData.Denom)
	// if denom trace path is empty (sending chain's native coin, e.g. 'adym'),
	// construct it from the packet destination port and channel, so that the ibc denom can be derived
	if denomTrace.Path == "" {
		denomTrace.Path = fmt.Sprintf("%s/%s", packet.GetDestPort(), packet.GetDestChannel())
	}

	if !im.keeper.HasDenomMetaData(ctx, denomTrace.IBCDenom()) {
		if err = im.createNewDenom(ctx, dm, denomTrace); err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}
	}

	return im.IBCModule.OnRecvPacket(ctx, packet, relayer)
}

func (im IBCMiddleware) createNewDenom(ctx sdk.Context, denomMetadata *banktypes.Metadata, denomTrace transfertypes.DenomTrace) error {
	denomUnits := make([]*banktypes.DenomUnit, 0, len(denomMetadata.DenomUnits))
	for _, du := range denomMetadata.DenomUnits {
		ndu := &banktypes.DenomUnit{
			Denom:    du.Denom,
			Exponent: du.Exponent,
			Aliases:  du.Aliases,
		}
		denomUnits = append(denomUnits, ndu)
	}

	newDenomMetadata := types.DenomMetadata{
		TokenMetadata: banktypes.Metadata{
			Description: denomMetadata.Description,
			DenomUnits:  denomUnits,
			Base:        denomTrace.IBCDenom(),
			Display:     denomMetadata.Display,
			Name:        denomMetadata.Name,
			Symbol:      denomMetadata.Symbol,
			URI:         denomMetadata.URI,
			URIHash:     denomMetadata.URIHash,
		},
		DenomTrace: denomTrace.GetFullDenomPath(),
	}

	if err := im.keeper.CreateDenomMetadata(ctx, newDenomMetadata); err != nil {
		return fmt.Errorf("create denom metadata: %w", err)
	}

	return nil
}
