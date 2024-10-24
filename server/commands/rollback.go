package commands

import (
	"fmt"
	"strconv"

	"github.com/dymensionxyz/dymint/block"
	dymintconf "github.com/dymensionxyz/dymint/config"
	dymintconv "github.com/dymensionxyz/dymint/conv"

	"github.com/dymensionxyz/dymint/store"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/proxy"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/dymensionxyz/dymension-rdk/utils"
	"github.com/spf13/cobra"
)

// RollbackCmd rollbacks the app multistore to specific height and updates dymint state according to it
func RollbackCmd(appCreator types.AppCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback [height]",
		Short: "rollback command used to move a full-node back to the state at the specified height.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := server.GetServerContextFromCmd(cmd)
			cfg := ctx.Config
			home := cfg.RootDir

			var heightInt int64
			if len(args) > 0 {
				height, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid height: %w", err)
				}
				heightInt = height
			} else {
				return fmt.Errorf("rollback height not specified")
			}

			db, err := utils.OpenDB(home)
			if err != nil {
				return err
			}
			defer func() {
				err = db.Close()
			}()

			nodeConfig := dymintconf.DefaultConfig("")
			err = nodeConfig.GetViperConfig(cmd, ctx.Viper.GetString(flags.FlagHome))
			if err != nil {
				return err
			}

			app := appCreator(ctx.Logger, db, nil, ctx.Viper)

			proxyApp := proxy.NewLocalClientCreator(app)
			ctx.Logger.Info("starting block manager with ABCI in-process")
			blockManager, err := liteBlockManager(cfg, nodeConfig, proxyApp)
			if err != nil {
				return fmt.Errorf("start lite block manager: %w", err)
			}

			// rollback the app multistore
			if err := app.CommitMultiStore().RollbackToVersion(heightInt); err != nil {
				return fmt.Errorf("app rollback to specific height: %w", err)
			}

			block, err := blockManager.Store.LoadBlock(uint64(heightInt))
			if err != nil {
				return fmt.Errorf("load block header: %w", err)
			}
			// rollback dymint state according to the app
			if err := blockManager.UpdateStateFromApp(block.Header.Hash()); err != nil {
				return fmt.Errorf("updating dymint from app state: %w", err)
			}
			fmt.Printf("RollApp state moved back to height %d successfully.\n", heightInt)
			return err
		},
	}

	dymintconf.AddNodeFlags(cmd)
	return cmd
}

func liteBlockManager(cfg *config.Config, dymintConf *dymintconf.NodeConfig, clientCreator proxy.ClientCreator) (*block.Manager, error) {

	genDocProvider := node.DefaultGenesisDocProviderFunc(cfg)

	privValKey, err := p2p.LoadOrGenNodeKey(cfg.PrivValidatorKeyFile())
	if err != nil {
		return nil, err
	}
	signingKey, err := dymintconv.GetNodeKey(privValKey)
	if err != nil {
		return nil, err
	}
	genesis, err := genDocProvider()
	if err != nil {
		return nil, err
	}

	err = dymintconv.GetNodeConfig(dymintConf, cfg)
	if err != nil {
		return nil, err
	}

	proxyApp := proxy.NewAppConns(clientCreator)
	if err := proxyApp.Start(); err != nil {
		return nil, fmt.Errorf("starting proxy app connections: %w", err)
	}

	var baseKV store.KV
	if dymintConf.RootDir == "" && dymintConf.DBPath == "" { // this is used for testing
		baseKV = store.NewDefaultInMemoryKVStore()
	} else {
		// TODO(omritoptx): Move dymint to const
		baseKV = store.NewDefaultKVStore(dymintConf.RootDir, dymintConf.DBPath, "dymint")
	}
	mainKV := store.NewPrefixKV(baseKV, []byte{0})
	s := store.New(mainKV)

	blockManager, err := block.NewManager(
		signingKey,
		*dymintConf,
		genesis,
		s,
		nil,
		proxyApp,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("BlockManager initialization error: %w", err)
	}

	return blockManager, nil
}
