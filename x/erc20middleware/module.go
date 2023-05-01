package erc20middleware

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	ibctransfer "github.com/cosmos/ibc-go/v5/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v5/modules/apps/transfer/keeper"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v5/modules/core/05-port/types"

	keeper "github.com/dymensionxyz/rollapp/x/erc20middleware/keeper"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
	_ porttypes.IBCModule   = AppModule{}
)

// AppModuleBasic embeds the IBC Transfer AppModuleBasic
type AppModuleBasic struct {
	*ibctransfer.AppModuleBasic
}

// AppModule represents the AppModule for this module
type AppModule struct {
	*ibctransfer.AppModule
	*ibctransfer.IBCModule
	keeper keeper.Keeper
}

// ICS 30 callbacks
// OnChanOpenInit implements the IBCModule interface
func (am AppModule) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, chanCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) error {
	// call underlying app's (transfer) callback
	return am.IBCModule.OnChanOpenInit(ctx, order, connectionHops, portID, channelID,
		chanCap, counterparty, version)
}

// OnChanOpenTry implements the IBCModule interface
func (am AppModule) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, chanCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, counterpartyVersion string,
) (version string, err error) {
	// call underlying app's (transfer) callback
	return am.IBCModule.OnChanOpenTry(ctx, order, connectionHops, portID, channelID,
		chanCap, counterparty, counterpartyVersion)
}

// OnChanOpenAck implements the IBCModule interface
func (am AppModule) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return am.IBCModule.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

// OnChanOpenConfirm implements the IBCModule interface
func (am AppModule) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	// call underlying app's OnChanOpenConfirm callback.
	return am.IBCModule.OnChanOpenConfirm(ctx, portID, channelID)
}

// OnChanCloseInit implements the IBCModule interface
func (am AppModule) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	// TODO: Unescrow all remaining funds for unprocessed packets
	return am.IBCModule.OnChanCloseInit(ctx, portID, channelID)
}

// OnChanCloseConfirm implements the IBCModule interface
func (am AppModule) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	// TODO: Unescrow all remaining funds for unprocessed packets
	return am.IBCModule.OnChanCloseConfirm(ctx, portID, channelID)
}

// NewAppModule creates a new 20-transfer module
func NewAppModule(k keeper.Keeper, ibckeeper ibctransferkeeper.Keeper) AppModule {
	ibcm := ibctransfer.NewIBCModule(ibckeeper)
	am := ibctransfer.NewAppModule(*k.Keeper)
	return AppModule{
		AppModule: &am,
		IBCModule: &ibcm,
		keeper:    k,
	}
}
