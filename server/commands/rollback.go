package commands

import (
	"fmt"
	"strconv"

	dymintconf "github.com/dymensionxyz/dymint/config"

	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"

	"github.com/dymensionxyz/dymension-rdk/utils"
)

const skipStorePruningFlag = "skip-store-pruning"

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

			genDocProvider := node.DefaultGenesisDocProviderFunc(cfg)
			genesis, err := genDocProvider()
			if err != nil {
				return err
			}

			blockManager, err := liteBlockManager(cfg, nodeConfig, genesis, nil, proxyApp, ctx.Logger)
			if err != nil {
				return fmt.Errorf("start lite block manager: %w", err)
			}

			// rollback the app multistore
			if err := app.CommitMultiStore().RollbackToVersion(heightInt); err != nil {
				return fmt.Errorf("app rollback to specific height: %w", err)
			}

			fmt.Printf("RollApp state moved back to height %d successfully.\n", heightInt)

			block, err := blockManager.Store.LoadBlock(uint64(heightInt))
			if err != nil {
				return fmt.Errorf("load block header: %w", err)
			}

			currentHeight := blockManager.State.Height()

			// rollback dymint state according to the app
			if err := blockManager.UpdateStateFromApp(block); err != nil {
				return fmt.Errorf("updating dymint from app state: %w", err)
			}

			_, err = blockManager.Store.SaveState(blockManager.State, nil)
			if err != nil {
				return fmt.Errorf("save state: %w", err)
			}

			skipStorePruning := ctx.Viper.GetBool(skipStorePruningFlag)

			if skipStorePruning {
				return nil
			}

			fmt.Printf("Pruning store from height %d \n", heightInt+1)

			// we try to prune height + 2, to prune all blocks in case a block its been already produced but not applied.
			pruned, err := blockManager.Store.PruneHeights(uint64(heightInt+1), currentHeight+2, ctx.Logger)
			if err != nil {
				ctx.Logger.Error("Error pruning block store.", "Error", err)
			}

			fmt.Println("Pruned blocks:", pruned)

			baseHeight, err := blockManager.Store.LoadBaseHeight()
			if err != nil {
				return nil
			}
			if baseHeight <= uint64(heightInt) {
				return nil
			}

			err = blockManager.Store.SaveBaseHeight(uint64(heightInt))
			if err != nil {
				ctx.Logger.Error("saving base height", "error", err)
			}

			return err
		},
	}
	cmd.Flags().Bool(skipStorePruningFlag, false, "rollback only app without pruning dymint store blocks")
	dymintconf.AddNodeFlags(cmd)
	return cmd
}
