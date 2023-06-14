//go:build evm
// +build evm

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"

	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	tmcfg "github.com/tendermint/tendermint/config"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	tmlog "github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	dymintserver "github.com/dymensionxyz/dymint/server"
	"github.com/dymensionxyz/rollapp/app"
	"github.com/dymensionxyz/rollapp/app/params"
	"github.com/dymensionxyz/rollapp/utils"

	evmflags "github.com/dymensionxyz/rollapp/app/evm/flags"
	ethermintclient "github.com/evmos/ethermint/client"
	"github.com/evmos/ethermint/crypto/hd"
	etherencoding "github.com/evmos/ethermint/encoding"
	evmserver "github.com/evmos/ethermint/server"
	evmconfig "github.com/evmos/ethermint/server/config"
)

const rollappAscii = `
███████ ██    ██ ███    ███     ██████   ██████  ██      ██       █████  ██████  ██████  
██      ██    ██ ████  ████     ██   ██ ██    ██ ██      ██      ██   ██ ██   ██ ██   ██ 
█████   ██    ██ ██ ████ ██     ██████  ██    ██ ██      ██      ███████ ██████  ██████  
██       ██  ██  ██  ██  ██     ██   ██ ██    ██ ██      ██      ██   ██ ██      ██      
███████   ████   ██      ██     ██   ██  ██████  ███████ ███████ ██   ██ ██      ██                                                                                                                                                            
`

func initStartCommandFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(evmflags.JSONRPCEnable, true, "Define if the JSON-RPC server should be enabled")
	cmd.Flags().StringSlice(evmflags.JSONRPCAPI, evmconfig.GetDefaultAPINamespaces(), "Defines a list of JSON-RPC namespaces that should be enabled")
	cmd.Flags().String(evmflags.JSONRPCAddress, evmconfig.DefaultJSONRPCAddress, "the JSON-RPC server address to listen on")
	cmd.Flags().String(evmflags.JSONWsAddress, evmconfig.DefaultJSONRPCWsAddress, "the JSON-RPC WS server address to listen on")
	cmd.Flags().Uint64(evmflags.JSONRPCGasCap, evmconfig.DefaultGasCap, "Sets a cap on gas that can be used in eth_call/estimateGas unit is aphoton (0=infinite)")     //nolint:lll
	cmd.Flags().Float64(evmflags.JSONRPCTxFeeCap, evmconfig.DefaultTxFeeCap, "Sets a cap on transaction fee that can be sent via the RPC APIs (1 = default 1 photon)") //nolint:lll
	cmd.Flags().Int32(evmflags.JSONRPCFilterCap, evmconfig.DefaultFilterCap, "Sets the global cap for total number of filters that can be created")
	cmd.Flags().Duration(evmflags.JSONRPCEVMTimeout, evmconfig.DefaultEVMTimeout, "Sets a timeout used for eth_call (0=infinite)")
	cmd.Flags().Duration(evmflags.JSONRPCHTTPTimeout, evmconfig.DefaultHTTPTimeout, "Sets a read/write timeout for json-rpc http server (0=infinite)")
	cmd.Flags().Duration(evmflags.JSONRPCHTTPIdleTimeout, evmconfig.DefaultHTTPIdleTimeout, "Sets a idle timeout for json-rpc http server (0=infinite)")
	cmd.Flags().Bool(evmflags.JSONRPCAllowUnprotectedTxs, evmconfig.DefaultAllowUnprotectedTxs, "Allow for unprotected (non EIP155 signed) transactions to be submitted via the node's RPC when the global parameter is disabled") //nolint:lll
	cmd.Flags().Int32(evmflags.JSONRPCLogsCap, evmconfig.DefaultLogsCap, "Sets the max number of results can be returned from single `eth_getLogs` query")
	cmd.Flags().Int32(evmflags.JSONRPCBlockRangeCap, evmconfig.DefaultBlockRangeCap, "Sets the max block range allowed for `eth_getLogs` query")
	cmd.Flags().Int(evmflags.JSONRPCMaxOpenConnections, evmconfig.DefaultMaxOpenConnections, "Sets the maximum number of simultaneous connections for the server listener") //nolint:lll
	cmd.Flags().Bool(evmflags.JSONRPCEnableIndexer, false, "Enable the custom tx indexer for json-rpc")
	cmd.Flags().Bool(evmflags.JSONRPCEnableMetrics, false, "Define if EVM rpc metrics server should be enabled")

	cmd.Flags().String(evmflags.EVMTracer, evmconfig.DefaultEVMTracer, "the EVM tracer type to collect execution traces from the EVM transaction execution (json|struct|access_list|markdown)") //nolint:lll
	cmd.Flags().Uint64(evmflags.EVMMaxTxGasWanted, evmconfig.DefaultMaxTxGasWanted, "the gas wanted for each eth tx returned in ante handler in check tx mode")                                 //nolint:lll
}

// NewRootCmd creates a new root rollappd command. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	ethEncodingConfig := etherencoding.MakeConfig(app.ModuleBasics)
	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: ethEncodingConfig.InterfaceRegistry,
		Codec:             ethEncodingConfig.Codec,
		TxConfig:          ethEncodingConfig.TxConfig,
		Amino:             ethEncodingConfig.Amino,
	}

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithKeyringOptions(hd.EthSecp256k1Option()).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   version.AppName,
		Short: rollappAscii,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customTMConfig := initTendermintConfig()
			customAppTemplate, customAppConfig := initAppConfig()
			err = server.InterceptConfigsPreRunHandler(
				cmd, customAppTemplate, customAppConfig, customTMConfig,
			)
			if err != nil {
				return err
			}

			//We initilaze dyming config after tendermint initialize, so we could read from it's configuration
			err = dymintserver.DymintConfigPreRunHandler(cmd)
			if err != nil {
				return err
			}

			return nil
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

// initTendermintConfig helps to override default Tendermint Config values.
// return tmcfg.DefaultConfig if no custom configuration is required for the application.
func initTendermintConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	// these values put a higher strain on node memory
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	customAppTemplate, customAppConfig := evmconfig.AppConfig("")

	srvCfg, ok := customAppConfig.(evmconfig.Config)
	if !ok {
		panic(fmt.Errorf("unknown app config type %T", customAppConfig))
	}

	// srvCfg.StateSync.SnapshotInterval = 5000
	// srvCfg.StateSync.SnapshotKeepRecent = 2
	// srvCfg.IAVLDisableFastNode = false

	return customAppTemplate, srvCfg
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {

	sdkconfig := sdk.GetConfig()
	utils.SetPrefixes(sdkconfig, app.AccountAddressPrefix)
	utils.SetBip44CoinType(sdkconfig)
	sdkconfig.Seal()

	ac := appCreator{
		encCfg: encodingConfig,
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.MigrateGenesisCmd(),
		genutilcli.GenTxCmd(
			app.ModuleBasics,
			encodingConfig.TxConfig,
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
		),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		// testnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),
		debug.Cmd(),
		config.Cmd(),
	)

	dymintserver.AddRollappCommands(rootCmd, app.DefaultNodeHome, ac.newApp, ac.appExport, addModuleInitFlags)
	rootCmd.AddCommand(StartCmd(ac.newApp, app.DefaultNodeHome))

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		ethermintclient.KeyCommands(app.DefaultNodeHome),
	)

	rootCmd.AddCommand(evmserver.NewIndexTxCmd())
	//FIXME: validate and uncomment or remove
	// rootCmd, err := srvflags.AddTxFlags(rootCmd)
	// if err != nil {
	// 	panic(err)
	// }

	// add rosetta
	// rootCmd.AddCommand(sdkserver.RosettaCommand(encodingConfig.InterfaceRegistry, encodingConfig.Codec))
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

// queryCommand returns the sub-command to send queries to the app
func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// txCommand returns the sub-command to send transactions to the app
func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetAuxToFeeCommand(),
	)

	app.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg params.EncodingConfig
}

func (ac appCreator) newApp(
	logger tmlog.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	return app.NewRollapp(
		logger,
		db,
		traceStore,
		true,
		skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		ac.encCfg,
		appOpts,
		baseappOptions...)
}

func (ac appCreator) appExport(
	logger tmlog.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
) (servertypes.ExportedApp, error) {
	var rollapp *app.App
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	loadLatest := height == -1
	rollapp = app.NewRollapp(
		logger,
		db,
		traceStore,
		loadLatest,
		map[int64]bool{},
		homePath,
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		ac.encCfg,
		appOpts,
	)

	if height != -1 {
		if err := rollapp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return rollapp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}
