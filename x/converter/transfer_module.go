package converter

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	ibctransfer "github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/dymensionxyz/dymension-rdk/x/converter/keeper"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic embeds the IBC Transfer AppModuleBasic
type AppModuleBasic struct {
	*ibctransfer.AppModuleBasic
}

// AppModule represents the AppModule for this module
type AppModule struct {
	*ibctransfer.AppModule
	keeper keeper.Keeper
}

// NewAppModule creates a new transfer app module with the wrapped keeper
func NewAppModule(wrappedKeeper keeper.Keeper) AppModule {
	// Create the base Evmos module with the embedded keeper
	baseModule := ibctransfer.NewAppModule(wrappedKeeper.Keeper)

	return AppModule{
		AppModule: &baseModule,
		keeper:    wrappedKeeper,
	}
}

// RegisterServices overrides the Evmos module's RegisterServices to use our wrapped keeper
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// Register our wrapped keeper as the MsgServer instead of the base keeper
	// This ensures our Transfer override is used
	transfertypes.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	transfertypes.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}
