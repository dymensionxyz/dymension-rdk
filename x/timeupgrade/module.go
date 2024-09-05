package timeupgrade

import (
	"encoding/json"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/keeper"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	abci "github.com/tendermint/tendermint/abci/types"

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
	//TODO implement me
	panic("implement me")
}

func (a AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	//TODO implement me
	panic("implement me")
}

func (a AppModuleBasic) DefaultGenesis(codec codec.JSONCodec) json.RawMessage {
	//TODO implement me
	panic("implement me")
}

func (a AppModuleBasic) ValidateGenesis(codec codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	//TODO implement me
	panic("implement me")
}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, mux *runtime.ServeMux) {
	//TODO implement me
	panic("implement me")
}

func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	//TODO implement me
	panic("implement me")
}

func (a AppModuleBasic) GetQueryCmd() *cobra.Command {
	//TODO implement me
	panic("implement me")
}

type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	upgradeKeeper upgradekeeper.Keeper
}

func (a AppModule) BeginBlock(context sdk.Context, block abci.RequestBeginBlock) {
	BeginBlocker(context, a.keeper, a.upgradeKeeper)
}

func (a AppModule) InitGenesis(context sdk.Context, jsonCodec codec.JSONCodec, message json.RawMessage) []abci.ValidatorUpdate {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) ExportGenesis(context sdk.Context, jsonCodec codec.JSONCodec) json.RawMessage {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) Route() sdk.Route {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) QuerierRoute() string {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) RegisterServices(configurator module.Configurator) {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) ConsensusVersion() uint64 {
	return 1
}
