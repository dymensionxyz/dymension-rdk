package commands

import (
	"fmt"

	dymintconf "github.com/dymensionxyz/dymint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/node"

	slregistry "github.com/dymensionxyz/dymint/settlement/registry"
	"github.com/tendermint/tendermint/proxy"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/dymensionxyz/dymension-rdk/utils"
)

// ValidateGenesisBridgeCmd runs init chain and genesis bridge validation against the hub
func ValidateGenesisBridgeCmd(appCreator types.AppCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate-genesis-bridge",
		Short: "validate init chain and genesis bridge against the hub.",
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

			genDocProvider := node.DefaultGenesisDocProviderFunc(cfg)
			genesis, err := genDocProvider()
			if err != nil {
				return err
			}

			slclient := slregistry.GetClient(slregistry.Client(nodeConfig.SettlementLayer))
			if slclient == nil {
				return fmt.Errorf("get settlement client: named: %s", nodeConfig.SettlementLayer)
			}
			err = slclient.Init(nodeConfig.SettlementConfig, genesis.ChainID, nil, ctx.Logger.With("module", "settlement_client"))
			if err != nil {
				return fmt.Errorf("settlement layer client initialization: %w", err)
			}

			blockManager, err := liteBlockManager(cfg, nodeConfig, genesis, slclient, proxyApp, ctx.Logger)
			if err != nil {
				return fmt.Errorf("start lite block manager: %w", err)
			}

			valset := []*tmtypes.Validator{tmtypes.NewValidator(ed25519.GenPrivKey().PubKey(), 1)}
			res, err := blockManager.Executor.InitChain(blockManager.Genesis, blockManager.GenesisChecksum, valset)
			if err != nil {
				return fmt.Errorf("Cannot validate genesis bridge data: %w.", err)
			}

			// validate the resulting genesis bridge data against the hub
			err = blockManager.ValidateGenesisBridgeData(res.GenesisBridgeDataBytes)
			if err != nil {
				return fmt.Errorf("Cannot validate genesis bridge data: %w.", err)
			}

			fmt.Println("Genesis bridge validated successfully.")

			return nil
		},
	}

	dymintconf.AddNodeFlags(cmd)
	return cmd
}
