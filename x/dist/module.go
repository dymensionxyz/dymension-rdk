package dist

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/distribution"

	distcli "github.com/cosmos/cosmos-sdk/x/distribution/client/cli"
	"github.com/dymensionxyz/dymension-rdk/x/dist/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/dist/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModule embeds the Cosmos SDK's x/distribution AppModuleBasic.
type AppModuleBasic struct {
	distribution.AppModuleBasic
}

// DefaultGenesis returns default genesis state as raw bytes for the distribution
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the distribution module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return types.ValidateGenesis(&data)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the distribution module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns the root tx command for the distribution module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return distcli.NewTxCmd()
}

// GetQueryCmd returns the root query command for the distribution module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return distcli.GetQueryCmd()
}

// AppModule embeds the Cosmos SDK's x/distribution AppModule where we only override specific methods.
type AppModule struct {
	distribution.AppModule

	keeper        keeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
}

// NewAppModule creates a new AppModule object using the native x/distribution AppModule constructor.
func NewAppModule(
	cdc codec.Codec, keeper keeper.Keeper, ak types.AccountKeeper,
	bk types.BankKeeper, sk types.StakingKeeper,
) AppModule {
	distAppMod := distribution.NewAppModule(cdc, keeper.Keeper, ak, bk, sk)
	return AppModule{
		AppModule:     distAppMod,
		keeper:        keeper,
		accountKeeper: ak,
		bankKeeper:    bk,
		stakingKeeper: sk,
	}
}

// BeginBlock returns the begin blocker for the distribution module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// TODO this is Tendermint-dependent
	// ref https://github.com/cosmos/cosmos-sdk/issues/3095
	if ctx.BlockHeight() > 1 {
		previousProposer := am.keeper.GetPreviousProposerConsAddr(ctx)
		am.keeper.AllocateTokens(ctx, previousProposer)
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	am.keeper.SetPreviousProposerConsAddr(ctx, consAddr)
}

// Name returns the distribution module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// QuerierRoute returns the distribution module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// InitGenesis performs genesis initialization for the distribution module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.Keeper.InitGenesis(ctx, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the distribution
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }
