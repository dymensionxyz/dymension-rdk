package dist

import (
	"encoding/json"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/dymensionxyz/dymension-rdk/x/dist/keeper"
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
		proposerConsAddr := am.keeper.GetPreviousProposerConsAddr(ctx)
		am.keeper.AllocateTokens(ctx, proposerConsAddr)
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	am.keeper.SetPreviousProposerConsAddr(ctx, consAddr)
}

func (am AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	defGenesis := types.DefaultGenesisState()

	// by default, all rewards goes to the governers
	defGenesis.Params.CommunityTax = sdk.ZeroDec()
	defGenesis.Params.BaseProposerReward = sdk.ZeroDec()
	defGenesis.Params.BonusProposerReward = sdk.ZeroDec()

	return cdc.MustMarshalJSON(defGenesis)
}
