package commands

import (
	"fmt"

	dymintconf "github.com/dymensionxyz/dymint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/tendermint/tendermint/proxy"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/dymensionxyz/dymension-rdk/utils"
)

// ValidateInitCmd runs init chain and validation bridge against the hub
func ValidateInitCmd(appCreator types.AppCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validateinitbridge",
		Short: " init chain and validation bridge against the hub.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := server.GetServerContextFromCmd(cmd)
			cfg := ctx.Config
			home := cfg.RootDir

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
			blockManager, err := liteBlockManager(cfg, nodeConfig, proxyApp, ctx.Logger)
			if err != nil {
				return fmt.Errorf("start lite block manager: %w", err)
			}

			valset := []*tmtypes.Validator{tmtypes.NewValidator(ed25519.GenPrivKey().PubKey(), 1)}
			res, err := blockManager.Executor.InitChain(blockManager.Genesis, blockManager.GenesisChecksum, valset)
			if err != nil {
				return err
			}

			// validate the resulting genesis bridge data against the hub
			err = blockManager.ValidateGenesisBridgeData(res.GenesisBridgeDataBytes)
			if err != nil {
				return fmt.Errorf("Cannot validate genesis bridge data: %w. Please call `$EXECUTABLE dymint unsafe-reset-all` before the next launch to reset this node to genesis state.", err)
			}

			fmt.Printf("Genesis validated successfully")

			return nil
		},
	}

	dymintconf.AddNodeFlags(cmd)
	return cmd
}
