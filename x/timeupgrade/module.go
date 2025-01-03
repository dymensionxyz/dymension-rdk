package timeupgrade

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.BeginBlockAppModule = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
)

type AppModuleBasic struct {
}

func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

func (a AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	types.RegisterCodec(amino)
}

func (a AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (a AppModuleBasic) DefaultGenesis(codec codec.JSONCodec) json.RawMessage {
	return nil
}

func (a AppModuleBasic) ValidateGenesis(codec codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, mux *runtime.ServeMux) {}

func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

func (a AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	upgradeKeeper upgradekeeper.Keeper
}

func NewAppModule(keeper keeper.Keeper, upgradeKeeper upgradekeeper.Keeper) *AppModule {
	return &AppModule{keeper: keeper, upgradeKeeper: upgradeKeeper}
}

func (a AppModule) BeginBlock(context sdk.Context, block abci.RequestBeginBlock) {
	BeginBlocker(context, a.keeper, a.upgradeKeeper)
}

func (a AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, message json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(message, &genesisState)

	a.keeper.InitGenesis(ctx, &genesisState)
	return []abci.ValidatorUpdate{}
}

func (a AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := a.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

func (a AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {
}

func (a AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, nil)
}

func (a AppModule) QuerierRoute() string {
	return types.RouterKey
}

func (a AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (a AppModule) RegisterServices(configurator module.Configurator) {
	types.RegisterMsgServer(configurator.MsgServer(), keeper.NewMsgServerImpl(a.keeper))
}

func (a AppModule) ConsensusVersion() uint64 {
	return 1
}
