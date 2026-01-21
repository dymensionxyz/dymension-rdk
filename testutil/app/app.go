package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	// Cosmos SDK & Tendermint
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	// Dymension RDK
	"github.com/dymensionxyz/dymension-rdk/x/converter"
	converterkeeper "github.com/dymensionxyz/dymension-rdk/x/converter/keeper"
	distr "github.com/dymensionxyz/dymension-rdk/x/dist"
	distrkeeper "github.com/dymensionxyz/dymension-rdk/x/dist/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/dividends"
	dividendskeeper "github.com/dymensionxyz/dymension-rdk/x/dividends/keeper"
	dividendstypes "github.com/dymensionxyz/dymension-rdk/x/dividends/types"
	"github.com/dymensionxyz/dymension-rdk/x/epochs"
	epochskeeper "github.com/dymensionxyz/dymension-rdk/x/epochs/keeper"
	epochstypes "github.com/dymensionxyz/dymension-rdk/x/epochs/types"
	"github.com/dymensionxyz/dymension-rdk/x/gasless"
	gaslessclient "github.com/dymensionxyz/dymension-rdk/x/gasless/client"
	gaslesskeeper "github.com/dymensionxyz/dymension-rdk/x/gasless/keeper"
	gaslesstypes "github.com/dymensionxyz/dymension-rdk/x/gasless/types"
	hubgenesis "github.com/dymensionxyz/dymension-rdk/x/hub-genesis"
	hubgenkeeper "github.com/dymensionxyz/dymension-rdk/x/hub-genesis/keeper"
	hubgentypes "github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	hubkeeper "github.com/dymensionxyz/dymension-rdk/x/hub/keeper"
	hubtypes "github.com/dymensionxyz/dymension-rdk/x/hub/types"
	"github.com/dymensionxyz/dymension-rdk/x/mint"
	mintkeeper "github.com/dymensionxyz/dymension-rdk/x/mint/keeper"
	minttypes "github.com/dymensionxyz/dymension-rdk/x/mint/types"
	"github.com/dymensionxyz/dymension-rdk/x/rollappparams"
	rollappparamskeeper "github.com/dymensionxyz/dymension-rdk/x/rollappparams/keeper"
	rollappparamstypes "github.com/dymensionxyz/dymension-rdk/x/rollappparams/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers"
	seqkeeper "github.com/dymensionxyz/dymension-rdk/x/sequencers/keeper"
	seqtypes "github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/dymensionxyz/dymension-rdk/x/staking"
	stakingkeeper "github.com/dymensionxyz/dymension-rdk/x/staking/keeper"
	"github.com/dymensionxyz/dymension-rdk/x/timeupgrade/keeper"
	timeupgradetypes "github.com/dymensionxyz/dymension-rdk/x/timeupgrade/types"
	"github.com/dymensionxyz/dymension-rdk/x/tokenfactory"
	tokenfactorykeeper "github.com/dymensionxyz/dymension-rdk/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/dymensionxyz/dymension-rdk/x/tokenfactory/types"

	// IBC
	ibctransfer "github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v6/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v6/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v6/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	ibctestingtypes "github.com/cosmos/ibc-go/v6/testing/types"

	// Wasm
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	// Ethermint / EVM
	cryptocodec "github.com/evmos/evmos/v12/crypto/codec"
	ethermint "github.com/evmos/evmos/v12/types"
	"github.com/evmos/evmos/v12/x/erc20"
	erc20keeper "github.com/evmos/evmos/v12/x/erc20/keeper"
	erc20types "github.com/evmos/evmos/v12/x/erc20/types"
	"github.com/evmos/evmos/v12/x/evm"
	evmkeeper "github.com/evmos/evmos/v12/x/evm/keeper"
	evmtypes "github.com/evmos/evmos/v12/x/evm/types"
	"github.com/evmos/evmos/v12/x/feemarket"
	feemarketkeeper "github.com/evmos/evmos/v12/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/evmos/v12/x/feemarket/types"
	evmostransferkeeper "github.com/evmos/evmos/v12/x/ibc/transfer/keeper"

	// Utils
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	_ "github.com/cosmos/cosmos-sdk/client/docs/statik"
)

const (
	AccountAddressPrefix = "rol"
	Name                 = "rollapp"
)

// KV Store Keys - Centralized management
var kvstorekeys = []string{
	authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey, seqtypes.StoreKey,
	minttypes.StoreKey, distrtypes.StoreKey, govtypes.StoreKey, paramstypes.StoreKey,
	ibchost.StoreKey, upgradetypes.StoreKey, evmtypes.StoreKey, erc20types.StoreKey,
	feemarkettypes.StoreKey, epochstypes.StoreKey, hubgentypes.StoreKey, hubtypes.StoreKey,
	ibctransfertypes.StoreKey, capabilitytypes.StoreKey, gaslesstypes.StoreKey, wasmtypes.StoreKey,
	tokenfactorytypes.StoreKey, rollappparamstypes.StoreKey, timeupgradetypes.StoreKey,
	dividendstypes.StoreKey,
}

// Module permissions
func GetMaccPerms() map[string][]string {
	return map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		hubgentypes.ModuleName:         {authtypes.Minter},
		evmtypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
		erc20types.ModuleName:          {authtypes.Minter, authtypes.Burner},
		wasmtypes.ModuleName:           {authtypes.Burner},
		tokenfactorytypes.ModuleName:   {authtypes.Minter, authtypes.Burner},
	}
}

type App struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry
	invCheckPeriod    uint

	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// Keepers (Organized by domain)
	AccountKeeper      authkeeper.AccountKeeper
	BankKeeper         bankkeeper.Keeper
	CapabilityKeeper   *capabilitykeeper.Keeper
	StakingKeeper      stakingkeeper.Keeper
	SequencersKeeper   seqkeeper.Keeper
	MintKeeper         mintkeeper.Keeper
	EpochsKeeper       epochskeeper.Keeper
	DistrKeeper        distrkeeper.Keeper
	GovKeeper          govkeeper.Keeper
	UpgradeKeeper      upgradekeeper.Keeper
	ParamsKeeper       paramskeeper.Keeper
	IBCKeeper          *ibckeeper.Keeper
	TransferKeeper     converterkeeper.Keeper
	WasmKeeper         wasmkeeper.Keeper
	TokenFactoryKeeper tokenfactorykeeper.Keeper
	EvmKeeper          *evmkeeper.Keeper
	FeeMarketKeeper    feemarketkeeper.Keeper
	Erc20Keeper        erc20keeper.Keeper

	// RDK Keepers
	HubKeeper           hubkeeper.Keeper
	HubGenesisKeeper    hubgenkeeper.Keeper
	RollappParamsKeeper rollappparamskeeper.Keeper
	DividendsKeeper     dividendskeeper.Keeper
	GaslessKeeper       gaslesskeeper.Keeper
	TimeUpgradeKeeper   keeper.Keeper

	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper     capabilitykeeper.ScopedKeeper

	mm           *module.Manager
	sm           *module.SimulationManager
	configurator module.Configurator
}

func NewRollapp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	skipUpgradeHeights map[int64]bool, homePath string, invCheckPeriod uint,
	encodingConfig EncodingConfig, enabledProposals []wasm.ProposalType,
	appOpts servertypes.AppOptions, wasmOpts []wasmkeeper.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {

	appCodec := encodingConfig.Codec
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(kvstorekeys...)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, keys); err != nil {
		panic(fmt.Errorf("failed to load state streaming services: %w", err))
	}

	app := &App{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()))

	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	app.ScopedIBCKeeper = app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	app.ScopedTransferKeeper = app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	app.ScopedWasmKeeper = app.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)
	app.CapabilityKeeper.Seal()

	// Initialize Keepers with Dependency Injection Pattern
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName),
		ethermint.ProtoAccount, GetMaccPerms(), sdk.Bech32PrefixAccAddr,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName), app.BlockedModuleAccountAddrs(),
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper,
		&app.Erc20Keeper, app.GetSubspace(stakingtypes.ModuleName),
	)

	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], app.GetSubspace(minttypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, app.EpochsKeeper, authtypes.FeeCollectorName,
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, &stakingKeeper, &app.SequencersKeeper,
		&app.Erc20Keeper, authtypes.FeeCollectorName,
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// IBC & EVM setup
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec, authtypes.NewModuleAddress(govtypes.ModuleName), keys[feemarkettypes.StoreKey],
		tkeys[feemarkettypes.TransientKey], app.GetSubspace(feemarkettypes.ModuleName),
	)

	app.EvmKeeper = evmkeeper.NewKeeper(
		appCodec, keys[evmtypes.StoreKey], tkeys[evmtypes.TransientKey],
		authtypes.NewModuleAddress(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		app.SequencersKeeper, app.FeeMarketKeeper, "", app.GetSubspace(evmtypes.ModuleName),
	)

	app.Erc20Keeper = erc20keeper.NewKeeper(
		keys[erc20types.StoreKey], appCodec, authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, app.EvmKeeper,
	)

	// IBC Transfer Middleware Stack
	var ics4Wrapper ibcporttypes.ICS4Wrapper
	ics4Wrapper = denommetadata.NewICS4Wrapper(app.IBCKeeper.ChannelKeeper, app.HubKeeper, app.BankKeeper, app.HubGenesisKeeper.GetState)
	ics4Wrapper = hubgenkeeper.NewICS4Wrapper(ics4Wrapper, app.HubGenesisKeeper)

	ibctransferKeeper := ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		ics4Wrapper, app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, app.ScopedTransferKeeper,
	)

	evmosTransferKeeper := evmostransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		ics4Wrapper, app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, app.ScopedTransferKeeper, &app.Erc20Keeper,
	)

	app.TransferKeeper = converterkeeper.NewTransferKeeper(ibctransferKeeper, evmosTransferKeeper, app.HubKeeper, app.BankKeeper)

	// Finalizing Module Manager
	app.mm = module.NewManager(app.getAppModules(encodingConfig)...)
	app.mm.SetOrderBeginBlockers(beginBlockersList...)
	app.mm.SetOrderEndBlockers(endBlockersList...)
	app.mm.SetOrderInitGenesis(initGenesisList...)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.setPostHandler()

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// Internal helper for module list to keep NewRollapp clean
func (app *App) getAppModules(encodingConfig EncodingConfig) []module.AppModule {
	return []module.AppModule{
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
		auth.NewAppModule(app.appCodec, app.AccountKeeper, nil),
		bank.NewAppModule(app.appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(app.appCodec, *app.CapabilityKeeper),
		gov.NewAppModule(app.appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper, app.GetSubspace(evmtypes.ModuleName)),
		feemarket.NewAppModule(app.FeeMarketKeeper, app.GetSubspace(feemarkettypes.ModuleName)),
		erc20.NewAppModule(app.Erc20Keeper, app.AccountKeeper, app.GetSubspace(erc20types.ModuleName)),
		// Diğer modüller...
	}
}

func (app *App) setPostHandler() {
	postHandler, err := posthandler.NewPostHandler(posthandler.HandlerOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to create post handler: %w", err))
	}
	app.SetPostHandler(postHandler)
}

// Boilerplate SDK methods (Name, BeginBlocker, EndBlocker etc.) keep as optimized versions...

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	pk := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)
	// Subspace registrations...
	return pk
}

//
