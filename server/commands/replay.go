package commands

import (
	"context"
	"fmt"
	"strconv"

	dymintconf "github.com/dymensionxyz/dymint/config"
	dymintconv "github.com/dymensionxyz/dymint/conv"
	dymintmemp "github.com/dymensionxyz/dymint/mempool"
	dymintnode "github.com/dymensionxyz/dymint/node"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/proxy"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/dymensionxyz/dymension-rdk/utils"
	rdklogger "github.com/dymensionxyz/dymension-rdk/utils/logger"
	"github.com/spf13/cobra"
)

// StartCmd runs the service passed in, either stand-alone or in-process with Dymint.
func ReplayCmd(appCreator types.AppCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replay",
		Short: "replay",
		Long:  `replay command.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := server.GetServerContextFromCmd(cmd)
			cfg := ctx.Config
			home := cfg.RootDir

			var heightInt uint64
			if len(args) > 0 {
				heightInt, _ = strconv.ParseUint(args[0], 10, 64)
			} else {
				return fmt.Errorf("replay height not specified")
			}

			db, err := utils.OpenDB(home)
			if err != nil {
				return err
			}
			defer db.Close()

			nodeConfig := dymintconf.DefaultConfig("", "")
			err = nodeConfig.GetViperConfig(cmd, ctx.Viper.GetString(flags.FlagHome))
			if err != nil {
				return err
			}

			app := appCreator(ctx.Logger, db, nil, ctx.Viper)

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

			err = dymintconv.GetNodeConfig(nodeConfig, cfg)
			if err != nil {
				return err
			}
			proxyApp := proxy.NewLocalClientCreator(app)
			ctx.Logger.Info("starting node with ABCI dymint in-process")
			node, err := dymintnode.NewNode(
				context.Background(),
				*nodeConfig,
				p2pKey,
				signingKey,
				proxyApp,
				genesis,
				ctx.Logger,
				dymintmemp.PrometheusMetrics("dymint"),
			)
			if err != nil {
				return err
			}

			// rollback the app multistore
			if err := app.CommitMultiStore().RollbackToVersion(int64(heightInt)); err != nil {
				return fmt.Errorf("failed to rollback to version: %w", err)
			}

			// rollback dymint state according to the app
			if err := node.BlockManager.UpdateStateFromApp(); err != nil {
				return fmt.Errorf("failed to rollback tendermint state: %w", err)
			}
			fmt.Printf("RollApp state moved back to height %d successfuly.\n", heightInt)
			return nil
		},
	}

	dymintconf.AddNodeFlags(cmd)
	rdklogger.AddLogFlags(cmd)
	return cmd
}
