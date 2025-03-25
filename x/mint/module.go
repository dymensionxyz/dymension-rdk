package mint

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/dymensionxyz/dymension-rdk/x/mint/client/cli"
	"github.com/dymensionxyz/dymension-rdk/x/mint/keeper"

	//"github.com/dymensionxyz/dymension-rdk/x/mint/simulation"
	"github.com/dymensionxyz/dymension-rdk/x/mint/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
	//_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the mint module.
type AppModuleBasic struct {
	cdc codec.Codec
}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the mint module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the mint module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// RegisterInterfaces registers the module's interface types.
func (b AppModuleBasic) RegisterInterfaces(_ cdctypes.InterfaceRegistry) {}

// DefaultGenesis returns default genesis state as raw bytes for the mint
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the mint module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return types.ValidateGenesis(data)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the mint module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		return
	}
}

// GetTxCmd returns no root tx command for the mint module.
func (AppModuleBasic) GetTxCmd() *cobra.Command { return nil }

// GetQueryCmd returns the root query command for the mint module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// ___________________________________________________________________________

// AppModule implements an application module for the mint module.
type AppModule struct {
	AppModuleBasic

	keeper     keeper.Keeper
	authKeeper types.AccountKeeper
	bankKeeper types.BankKeeper

	// this is a workaround used for migration v1 to v2 params
	emptySubspace paramstypes.Subspace
}

// NewAppModule creates a new AppModule object.
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, emptySubspace paramstypes.Subspace) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		authKeeper:     ak,
		bankKeeper:     bk,
		emptySubspace:  emptySubspace,
	}
}

// Name returns the mint module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the mint module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the mint module.
func (AppModule) Route() sdk.Route { return sdk.Route{} }

// QuerierRoute returns the mint module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// LegacyQuerierHandler returns the x/mint module's sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(sdk.Context, []string, abci.RequestQuery) ([]byte, error) {
		return nil, fmt.Errorf("legacy querier not supported for the x/%s module", types.ModuleName)
	}
}

// RegisterServices registers a gRPC query service to respond to the
// module-specific gRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQuerier(am.keeper))

	migrator := keeper.NewMigrator(am.keeper, am.emptySubspace)

	// register v1 -> v2 migration
	if err := cfg.RegisterMigration(types.ModuleName, 1, migrator.Migrate1to2); err != nil {
		panic(fmt.Errorf("failed to migrate %s to v2: %w", types.ModuleName, err))
	}

}

// InitGenesis performs genesis initialization for the mint module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	am.keeper.InitGenesis(ctx, &genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the mint
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the mint module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
}

// EndBlock returns the end blocker for the mint module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ConsensusVersion implements AppModule/ConsensusVersion.
// 2 - moved mint denom from params to minter
func (AppModule) ConsensusVersion() uint64 { return 2 }

// ___________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the mint module.
// func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
// 	simulation.RandomizedGenState(simState)
// }

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized mint param changes for the simulator.
// func (AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
// 	return simulation.ParamChanges(r)
// }

// RegisterStoreDecoder registers a decoder for mint module's types.
// func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
// 	sdr[types.StoreKey] = simulation.NewDecodeStore(am.cdc)
// }

// WeightedOperations doesn't return any mint module operation.
func (AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}
