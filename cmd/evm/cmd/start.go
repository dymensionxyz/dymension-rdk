//go:build evm
// +build evm

package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	pruningtypes "github.com/cosmos/cosmos-sdk/pruning/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	servergrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	"github.com/cosmos/cosmos-sdk/server/rosetta"
	crgserver "github.com/cosmos/cosmos-sdk/server/rosetta/lib/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	dymintconf "github.com/dymensionxyz/dymint/config"
	dymintconv "github.com/dymensionxyz/dymint/conv"
	dymintnode "github.com/dymensionxyz/dymint/node"
	dymintrpc "github.com/dymensionxyz/dymint/rpc"
	"github.com/dymensionxyz/rollapp/app"
	"github.com/dymensionxyz/rollapp/utils"

	"github.com/evmos/ethermint/indexer"
	ethserver "github.com/evmos/ethermint/server"
	"github.com/evmos/ethermint/server/config"
	ethconfig "github.com/evmos/ethermint/server/config"
	srvflags "github.com/evmos/ethermint/server/flags"
	ethermint "github.com/evmos/ethermint/types"
)

const (
	// Tendermint full-node start flags
	flagWithTendermint     = "with-tendermint"
	flagAddress            = "address"
	flagTransport          = "transport"
	flagTraceStore         = "trace-store"
	flagCPUProfile         = "cpu-profile"
	FlagMinGasPrices       = "minimum-gas-prices"
	FlagHaltHeight         = "halt-height"
	FlagHaltTime           = "halt-time"
	FlagInterBlockCache    = "inter-block-cache"
	FlagUnsafeSkipUpgrades = "unsafe-skip-upgrades"
	FlagTrace              = "trace"
	FlagInvCheckPeriod     = "inv-check-period"

	FlagPruning             = "pruning"
	FlagPruningKeepRecent   = "pruning-keep-recent"
	FlagPruningInterval     = "pruning-interval"
	FlagIndexEvents         = "index-events"
	FlagMinRetainBlocks     = "min-retain-blocks"
	FlagIAVLCacheSize       = "iavl-cache-size"
	FlagDisableIAVLFastNode = "iavl-disable-fastnode"
	FlagIAVLLazyLoading     = "iavl-lazy-loading"

	// state sync-related flags
	FlagStateSyncSnapshotInterval   = "state-sync.snapshot-interval"
	FlagStateSyncSnapshotKeepRecent = "state-sync.snapshot-keep-recent"

	// api-related flags
	FlagAPIEnable             = "api.enable"
	FlagAPISwagger            = "api.swagger"
	FlagAPIAddress            = "api.address"
	FlagAPIMaxOpenConnections = "api.max-open-connections"
	FlagRPCReadTimeout        = "api.rpc-read-timeout"
	FlagRPCWriteTimeout       = "api.rpc-write-timeout"
	FlagRPCMaxBodyBytes       = "api.rpc-max-body-bytes"
	FlagAPIEnableUnsafeCORS   = "api.enabled-unsafe-cors"

	// gRPC-related flags
	flagGRPCOnly       = "grpc-only"
	flagGRPCEnable     = "grpc.enable"
	flagGRPCAddress    = "grpc.address"
	flagGRPCWebEnable  = "grpc-web.enable"
	flagGRPCWebAddress = "grpc-web.address"

	// logging flags
	flagLogFile                = "log-file"
	flagLogLevel               = "log_level"
	flagMaxLogSize             = "max-log-size"
	flagModuleLogLevelOverride = "module-log-level-override"
)

// StartCmd runs the service passed in, either stand-alone or in-process with Dymint.
func StartCmd(appCreator types.AppCreator, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		Long: `Run the full node application with Tendermint in or out of process. By
default, the application will run with Tendermint in process.

Pruning options can be provided via the '--pruning' flag or alternatively with '--pruning-keep-recent', and
'pruning-interval' together.

For '--pruning' the options are as follows:

default: the last 362880 states are kept, pruning at 10 block intervals
nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
everything: 2 latest states will be kept; pruning at 10 block intervals.
custom: allow pruning options to be manually specified through 'pruning-keep-recent', and 'pruning-interval'

Node halting configurations exist in the form of two flags: '--halt-height' and '--halt-time'. During
the ABCI Commit phase, the node will check if the current block height is greater than or equal to
the halt-height or if the current block time is greater than or equal to the halt-time. If so, the
node will attempt to gracefully shutdown and the block will not be committed. In addition, the node
will not be able to commit subsequent blocks.

For profiling and benchmarking purposes, CPU profiling can be enabled via the '--cpu-profile' flag
which accepts a path for the resulting pprof file.
`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)

			// Bind flags to the Context's Viper so the app construction can set
			// options accordingly.
			err := serverCtx.Viper.BindPFlags(cmd.Flags())
			if err != nil {
				return err
			}

			_, err = server.GetPruningOptionsFromFlags(serverCtx.Viper)
			return err
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)

			// setup logging
			moduleOverrides := utils.ConvertStringToStringMap(serverCtx.Viper.GetString(flagModuleLogLevelOverride), ",", ":")

			log_path := serverCtx.Viper.GetString(flagLogFile)
			maxLogSize, err := strconv.Atoi(serverCtx.Viper.GetString(flagMaxLogSize))
			if err != nil {
				return err
			}
			if maxLogSize <= 0 {
				return fmt.Errorf("max log size <=0 not supported")
			}

			serverCtx.Logger = app.NewLogger(log_path, maxLogSize, serverCtx.Viper.GetString(flagLogLevel), moduleOverrides)

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			withTM, _ := cmd.Flags().GetBool(flagWithTendermint)
			if !withTM {
				serverCtx.Logger.Error("starting ABCI without Dymint not supported")
				return fmt.Errorf("starting ABCI without Dymint not supported")
			}
			serverCtx.Logger.Info("Unlocking keyring")

			// fire unlock precess for keyring
			keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
			if keyringBackend == keyring.BackendFile {
				_, err = clientCtx.Keyring.List()
				if err != nil {
					return err
				}
			}

			serverCtx.Logger.Info("starting ABCI with Dymint")

			dymconfig := dymintconf.DefaultConfig()
			err = dymconfig.GetViperConfig(cmd, serverCtx.Viper.GetString(flags.FlagHome))
			if err != nil {
				return err
			}

			// amino is needed here for backwards compatibility of REST routes
			err = wrapCPUProfile(serverCtx, func() error {
				return startInProcess(serverCtx, clientCtx, dymconfig, appCreator)
			})

			errCode, ok := err.(server.ErrorCode)
			if !ok {
				return err
			}

			serverCtx.Logger.Debug(fmt.Sprintf("received quit signal: %d", errCode.Code))
			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().Bool(flagWithTendermint, true, "Run abci app embedded in-process with tendermint")
	cmd.Flags().String(flagAddress, "tcp://0.0.0.0:26658", "Listen address")
	cmd.Flags().String(flagTransport, "socket", "Transport protocol: socket, grpc")
	cmd.Flags().String(flagTraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().String(FlagMinGasPrices, "", "Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum (e.g. 0.01photino;0.0001stake)")
	cmd.Flags().IntSlice(FlagUnsafeSkipUpgrades, []int{}, "Skip a set of upgrade heights to continue the old binary")
	cmd.Flags().Uint64(FlagHaltHeight, 0, "Block height at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Uint64(FlagHaltTime, 0, "Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Bool(FlagInterBlockCache, true, "Enable inter-block caching")
	cmd.Flags().String(flagCPUProfile, "", "Enable CPU profiling and write to the provided file")
	cmd.Flags().Bool(FlagTrace, false, "Provide full stack traces for errors in ABCI Log")
	cmd.Flags().String(FlagPruning, pruningtypes.PruningOptionDefault, "Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().Uint64(FlagPruningKeepRecent, 0, "Number of recent heights to keep on disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(FlagPruningInterval, 0, "Height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint(FlagInvCheckPeriod, 0, "Assert registered invariants every N blocks")
	cmd.Flags().Uint64(FlagMinRetainBlocks, 0, "Minimum block height offset during ABCI commit to prune Tendermint blocks")

	cmd.Flags().Bool(FlagAPIEnable, false, "Define if the API server should be enabled")
	cmd.Flags().Bool(FlagAPISwagger, false, "Define if swagger documentation should automatically be registered (Note: the API must also be enabled)")
	cmd.Flags().String(FlagAPIAddress, serverconfig.DefaultAPIAddress, "the API server address to listen on")
	cmd.Flags().Uint(FlagAPIMaxOpenConnections, 1000, "Define the number of maximum open connections")
	cmd.Flags().Uint(FlagRPCReadTimeout, 10, "Define the Tendermint RPC read timeout (in seconds)")
	cmd.Flags().Uint(FlagRPCWriteTimeout, 0, "Define the Tendermint RPC write timeout (in seconds)")
	cmd.Flags().Uint(FlagRPCMaxBodyBytes, 1000000, "Define the Tendermint maximum response body (in bytes)")
	cmd.Flags().Bool(FlagAPIEnableUnsafeCORS, false, "Define if CORS should be enabled (unsafe - use it at your own risk)")

	cmd.Flags().Bool(flagGRPCOnly, false, "Start the node in gRPC query only mode (no Tendermint process is started)")
	cmd.Flags().Bool(flagGRPCEnable, true, "Define if the gRPC server should be enabled")
	cmd.Flags().String(flagGRPCAddress, serverconfig.DefaultGRPCAddress, "the gRPC server address to listen on")

	cmd.Flags().Bool(flagGRPCWebEnable, true, "Define if the gRPC-Web server should be enabled. (Note: gRPC must also be enabled)")
	cmd.Flags().String(flagGRPCWebAddress, serverconfig.DefaultGRPCWebAddress, "The gRPC-Web server address to listen on")

	cmd.Flags().Uint64(FlagStateSyncSnapshotInterval, 0, "State sync snapshot interval")
	cmd.Flags().Uint32(FlagStateSyncSnapshotKeepRecent, 2, "State sync snapshot to keep")

	cmd.Flags().Bool(FlagDisableIAVLFastNode, false, "Disable fast node for IAVL tree")

	cmd.Flags().Bool(srvflags.JSONRPCEnable, true, "Define if the JSON-RPC server should be enabled")
	cmd.Flags().StringSlice(srvflags.JSONRPCAPI, config.GetDefaultAPINamespaces(), "Defines a list of JSON-RPC namespaces that should be enabled")
	cmd.Flags().String(srvflags.JSONRPCAddress, config.DefaultJSONRPCAddress, "the JSON-RPC server address to listen on")
	cmd.Flags().String(srvflags.JSONWsAddress, config.DefaultJSONRPCWsAddress, "the JSON-RPC WS server address to listen on")
	cmd.Flags().Uint64(srvflags.JSONRPCGasCap, config.DefaultGasCap, "Sets a cap on gas that can be used in eth_call/estimateGas unit is aphoton (0=infinite)")     //nolint:lll
	cmd.Flags().Float64(srvflags.JSONRPCTxFeeCap, config.DefaultTxFeeCap, "Sets a cap on transaction fee that can be sent via the RPC APIs (1 = default 1 photon)") //nolint:lll
	cmd.Flags().Int32(srvflags.JSONRPCFilterCap, config.DefaultFilterCap, "Sets the global cap for total number of filters that can be created")
	cmd.Flags().Duration(srvflags.JSONRPCEVMTimeout, config.DefaultEVMTimeout, "Sets a timeout used for eth_call (0=infinite)")
	cmd.Flags().Duration(srvflags.JSONRPCHTTPTimeout, config.DefaultHTTPTimeout, "Sets a read/write timeout for json-rpc http server (0=infinite)")
	cmd.Flags().Duration(srvflags.JSONRPCHTTPIdleTimeout, config.DefaultHTTPIdleTimeout, "Sets a idle timeout for json-rpc http server (0=infinite)")
	cmd.Flags().Bool(srvflags.JSONRPCAllowUnprotectedTxs, config.DefaultAllowUnprotectedTxs, "Allow for unprotected (non EIP155 signed) transactions to be submitted via the node's RPC when the global parameter is disabled") //nolint:lll
	cmd.Flags().Int32(srvflags.JSONRPCLogsCap, config.DefaultLogsCap, "Sets the max number of results can be returned from single `eth_getLogs` query")
	cmd.Flags().Int32(srvflags.JSONRPCBlockRangeCap, config.DefaultBlockRangeCap, "Sets the max block range allowed for `eth_getLogs` query")
	cmd.Flags().Int(srvflags.JSONRPCMaxOpenConnections, config.DefaultMaxOpenConnections, "Sets the maximum number of simultaneous connections for the server listener") //nolint:lll
	cmd.Flags().Bool(srvflags.JSONRPCEnableIndexer, false, "Enable the custom tx indexer for json-rpc")
	cmd.Flags().Bool(srvflags.JSONRPCEnableMetrics, false, "Define if EVM rpc metrics server should be enabled")

	cmd.Flags().String(flagLogLevel, "debug", "Log leve. one of [\"debug\", \"info\", \"warn\", \"error\", \"dpanic\", \"panic\", \"fatal\"]")
	cmd.Flags().String(flagLogFile, "", "log file full path. If not set, logs to stdout")
	cmd.Flags().String(flagMaxLogSize, "1000", "Max log size in MB")

	//dev option
	cmd.Flags().String(flagModuleLogLevelOverride, "", "Override module log level for customizable logging. For example \"module1:info,module2:error\"")
	cmd.Flags().MarkHidden(flagModuleLogLevelOverride)

	// add support for all Tendermint-specific command line options
	dymintconf.AddNodeFlags(cmd)
	return cmd
}

func startInProcess(ctx *server.Context, clientCtx client.Context, nodeConfig *dymintconf.NodeConfig, appCreator types.AppCreator) error {
	cfg := ctx.Config
	home := cfg.RootDir

	db, err := utils.OpenDB(home)
	if err != nil {
		return err
	}

	traceWriterFile := ctx.Viper.GetString(flagTraceStore)
	traceWriter, err := utils.OpenTraceWriter(traceWriterFile)
	if err != nil {
		return err
	}

	config, err := ethconfig.GetConfig(ctx.Viper)
	if err != nil {
		return err
	}
	gRPCOnly := ctx.Viper.GetBool(flagGRPCOnly)

	if err := config.ValidateBasic(); err != nil {
		return err
	}

	app := appCreator(ctx.Logger, db, traceWriter, ctx.Viper)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return err
	}

	privValKey, err := p2p.LoadOrGenNodeKey(cfg.PrivValidatorKeyFile())
	if err != nil {
		return err
	}

	genDocProvider := node.DefaultGenesisDocProviderFunc(cfg)

	// keys in dymint format
	p2pKey, err := dymintconv.GetNodeKey(nodeKey)
	if err != nil {
		return err
	}
	signingKey, err := dymintconv.GetNodeKey(privValKey)
	if err != nil {
		return err
	}
	genesis, err := genDocProvider()
	if err != nil {
		return err
	}

	ctx.Logger.Info("starting node with ABCI dymint in-process")
	tmNode, err := dymintnode.NewNode(
		context.Background(),
		*nodeConfig,
		p2pKey,
		signingKey,
		proxy.NewLocalClientCreator(app),
		genesis,
		ctx.Logger,
	)
	if err != nil {
		return err
	}

	dymserver := dymintrpc.NewServer(tmNode, cfg.RPC, ctx.Logger)
	err = dymserver.Start()
	if err != nil {
		return err
	}

	ctx.Logger.Debug("initialization: tmNode created")
	if err := tmNode.Start(); err != nil {
		return err
	}

	// Add the tx service to the gRPC router. We only need to register this
	// service if API or gRPC is enabled, and avoid doing so in the general
	// case, because it spawns a new local tendermint RPC client.
	if config.API.Enable || config.GRPC.Enable || config.JSONRPC.Enable || config.JSONRPC.EnableIndexer {
		clientCtx := clientCtx.WithClient(dymserver.Client())

		app.RegisterTxService(clientCtx)
		app.RegisterTendermintService(clientCtx)

		if a, ok := app.(types.ApplicationQueryService); ok {
			a.RegisterNodeService(clientCtx)
		}
	}

	metrics, err := startTelemetry(config)
	if err != nil {
		return err
	}

	/* --------------------------- Adding EVM indexer --------------------------- */
	var idxer ethermint.EVMTxIndexer
	if config.JSONRPC.EnableIndexer {
		idxDB, err := ethserver.OpenIndexerDB(home, server.GetAppDBBackend(ctx.Viper))
		if err != nil {
			ctx.Logger.Error("failed to open evm indexer DB", "error", err.Error())
			return err
		}
		idxLogger := ctx.Logger.With("module", "evmindex")
		idxer = indexer.NewKVIndexer(idxDB, idxLogger, clientCtx)
		indexerService := ethserver.NewEVMIndexerService(idxer, clientCtx.Client)
		indexerService.SetLogger(idxLogger)

		errCh := make(chan error)
		go func() {
			if err := indexerService.Start(); err != nil {
				errCh <- err
			}
		}()

		select {
		case err := <-errCh:
			return err
		case <-time.After(types.ServerStartTime): // assume server started successfully
		}
	}

	if config.API.Enable || config.JSONRPC.Enable {
		genDoc, err := genDocProvider()
		if err != nil {
			return err
		}

		clientCtx := clientCtx.WithHomeDir(home).WithChainID(genDoc.ChainID)

		// Set `GRPCClient` to `clientCtx` to enjoy concurrent grpc query.
		// only use it if gRPC server is enabled.
		if config.GRPC.Enable {
			_, port, err := net.SplitHostPort(config.GRPC.Address)
			if err != nil {
				return errors.Wrapf(err, "invalid grpc address %s", config.GRPC.Address)
			}

			maxSendMsgSize := config.GRPC.MaxSendMsgSize
			if maxSendMsgSize == 0 {
				maxSendMsgSize = serverconfig.DefaultGRPCMaxSendMsgSize
			}

			maxRecvMsgSize := config.GRPC.MaxRecvMsgSize
			if maxRecvMsgSize == 0 {
				maxRecvMsgSize = serverconfig.DefaultGRPCMaxRecvMsgSize
			}

			grpcAddress := fmt.Sprintf("127.0.0.1:%s", port)

			// If grpc is enabled, configure grpc client for grpc gateway and json-rpc.
			grpcClient, err := grpc.Dial(
				grpcAddress,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithDefaultCallOptions(
					grpc.ForceCodec(codec.NewProtoCodec(clientCtx.InterfaceRegistry).GRPCCodec()),
					grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
					grpc.MaxCallSendMsgSize(maxSendMsgSize),
				),
			)
			if err != nil {
				return err
			}

			clientCtx = clientCtx.WithGRPCClient(grpcClient)
			ctx.Logger.Debug("grpc client assigned to client context", "target", grpcAddress)
		}
	}

	var apiSrv *api.Server
	if config.API.Enable {
		apiSrv = api.New(clientCtx, ctx.Logger.With("server", "api"))
		app.RegisterAPIRoutes(apiSrv, config.API)
		if config.Telemetry.Enabled {
			apiSrv.SetTelemetry(metrics)
		}
		errCh := make(chan error)
		go func() {
			if err := apiSrv.Start(config.Config); err != nil {
				errCh <- err
			}
		}()

		select {
		case err := <-errCh:
			return err
		case <-time.After(types.ServerStartTime): // assume server started successfully
		}

		defer apiSrv.Close()
	}

	var (
		grpcSrv    *grpc.Server
		grpcWebSrv *http.Server
	)

	if config.GRPC.Enable {
		grpcSrv, err = servergrpc.StartGRPCServer(clientCtx, app, config.GRPC)
		if err != nil {
			return err
		}
		defer grpcSrv.Stop()
		if config.GRPCWeb.Enable {
			grpcWebSrv, err = servergrpc.StartGRPCWeb(grpcSrv, config.Config)
			if err != nil {
				ctx.Logger.Error("failed to start grpc-web http server", "error", err)
				return err
			}
			defer func() {
				if err := grpcWebSrv.Close(); err != nil {
					ctx.Logger.Error("failed to close grpc-web http server: ", err)
				}
			}()
		}
	}

	var (
		httpSrv     *http.Server
		httpSrvDone chan struct{}
	)

	if config.JSONRPC.Enable {
		genDoc, err := genDocProvider()
		if err != nil {
			return err
		}

		clientCtx := clientCtx.WithChainID(genDoc.ChainID)

		tmEndpoint := "/websocket"
		tmRPCAddr := cfg.RPC.ListenAddress
		httpSrv, httpSrvDone, err = ethserver.StartJSONRPC(ctx, clientCtx, tmRPCAddr, tmEndpoint, &config, idxer)
		if err != nil {
			return err
		}
		defer func() {
			shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFn()
			if err := httpSrv.Shutdown(shutdownCtx); err != nil {
				ctx.Logger.Error("HTTP server shutdown produced a warning", "error", err.Error())
			} else {
				ctx.Logger.Info("HTTP server shut down, waiting 5 sec")
				select {
				case <-time.Tick(5 * time.Second):
				case <-httpSrvDone:
				}
			}
		}()
	}

	if gRPCOnly {
		// wait for signal capture and gracefully return
		return utils.WaitForQuitSignals()
	}

	var rosettaSrv crgserver.Server
	if config.Rosetta.Enable {
		offlineMode := config.Rosetta.Offline

		// If GRPC is not enabled rosetta cannot work in online mode, so it works in
		// offline mode.
		if !config.GRPC.Enable {
			offlineMode = true
		}

		minGasPrices, err := sdk.ParseDecCoins(config.MinGasPrices)
		if err != nil {
			ctx.Logger.Error("failed to parse minimum-gas-prices: ", err)
			return err
		}

		conf := &rosetta.Config{
			Blockchain:          config.Rosetta.Blockchain,
			Network:             config.Rosetta.Network,
			TendermintRPC:       ctx.Config.RPC.ListenAddress,
			GRPCEndpoint:        config.GRPC.Address,
			Addr:                config.Rosetta.Address,
			Retries:             config.Rosetta.Retries,
			Offline:             offlineMode,
			GasToSuggest:        config.Rosetta.GasToSuggest,
			EnableFeeSuggestion: config.Rosetta.EnableFeeSuggestion,
			GasPrices:           minGasPrices.Sort(),
			Codec:               clientCtx.Codec.(*codec.ProtoCodec),
			InterfaceRegistry:   clientCtx.InterfaceRegistry,
		}

		rosettaSrv, err = rosetta.ServerFromConfig(conf)
		if err != nil {
			return err
		}

		errCh := make(chan error)
		go func() {
			if err := rosettaSrv.Start(); err != nil {
				errCh <- err
			}
		}()

		select {
		case err := <-errCh:
			return err
		case <-time.After(types.ServerStartTime): // assume server started successfully
		}
	}

	defer func() {
		if tmNode != nil && tmNode.IsRunning() {
			_ = tmNode.Stop()
		}

		if apiSrv != nil {
			_ = apiSrv.Close()
		}

		ctx.Logger.Info("exiting...")
	}()

	// wait for signal capture and gracefully return
	return utils.WaitForQuitSignals()
}

func startTelemetry(cfg ethconfig.Config) (*telemetry.Metrics, error) {
	if !cfg.Telemetry.Enabled {
		return nil, nil
	}
	return telemetry.New(cfg.Telemetry)
}

// wrapCPUProfile runs callback in a goroutine, then wait for quit signals.
func wrapCPUProfile(ctx *server.Context, callback func() error) error {
	if cpuProfile := ctx.Viper.GetString(flagCPUProfile); cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			return err
		}

		ctx.Logger.Info("starting CPU profiler", "profile", cpuProfile)
		if err := pprof.StartCPUProfile(f); err != nil {
			return err
		}

		defer func() {
			ctx.Logger.Info("stopping CPU profiler", "profile", cpuProfile)
			pprof.StopCPUProfile()
			if err := f.Close(); err != nil {
				ctx.Logger.Info("failed to close cpu-profile file", "profile", cpuProfile, "err", err.Error())
			}
		}()
	}

	errCh := make(chan error)
	go func() {
		errCh <- callback()
	}()

	select {
	case err := <-errCh:
		return err

	case <-time.After(types.ServerStartTime):
	}

	return <-errCh
}
