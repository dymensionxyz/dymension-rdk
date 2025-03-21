package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/libp2p/go-libp2p"
	"github.com/spf13/cobra"
	tmcmd "github.com/tendermint/tendermint/cmd/cometbft/commands"
	cmtos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/p2p"

	"github.com/dymensionxyz/dymension-rdk/server/commands"
	dymintcmds "github.com/dymensionxyz/dymint/cmd/dymint/commands"
	"github.com/dymensionxyz/dymint/conv"
)

// add Rollapp commands
func AddRollappCommands(rootCmd *cobra.Command, defaultNodeHome string, appCreator types.AppCreator, appExport types.AppExporter) {
	dymintCmd := &cobra.Command{
		Use:   "dymint",
		Short: "dymint subcommands",
	}

	dymintCmd.AddCommand(
		ShowSequencer(),
		ShowNodeIDCmd(),
		commands.InspectStateCmd(),
		commands.ShowP2PInfoCmd(),

		ResetAll(),
		server.VersionCmd(),
		ResetState(),
		dymintcmds.Run3dMigrationCmd(),
	)

	rootCmd.AddCommand(
		dymintCmd,
		server.ExportCmd(appExport, defaultNodeHome),
		version.NewVersionCommand(),
		commands.RollbackCmd(appCreator),
		commands.ValidateGenesisBridgeCmd(appCreator),
		pruning.Cmd(appCreator, defaultNodeHome),
	)
}

// ShowNodeIDCmd - ported from Tendermint, dump node ID to stdout
func ShowNodeIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-node-id",
		Short: "Show this node's ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			cfg := serverCtx.Config

			nodeKey, err := p2p.LoadNodeKey(cfg.NodeKeyFile())
			if err != nil {
				return err
			}
			signingKey, err := conv.GetNodeKey(nodeKey)
			if err != nil {
				return err
			}
			// convert nodeKey to libp2p key
			// nolint: typecheck
			host, err := libp2p.New(libp2p.Identity(signingKey))
			if err != nil {
				return err
			}

			fmt.Println(host.ID())
			return nil
		},
	}
}

func ShowSequencer() *cobra.Command {
	showSequencer := server.ShowValidatorCmd()
	showSequencer.Use = "show-sequencer"
	showSequencer.Short = "Show the current sequencer address"

	return showSequencer
}

func ResetAll() *cobra.Command {
	resetAll := tmcmd.ResetAllCmd
	resetAll.Short = "(unsafe) Remove all the data and WAL, reset this node's sequencer to genesis state"

	return resetAll
}

// ResetState removes all the block state stored in dymint
func ResetState() *cobra.Command {
	return &cobra.Command{
		Use:     "reset-state",
		Aliases: []string{"reset_state"},
		Short:   "Remove all the data and WAL",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			config := server.GetServerContextFromCmd(cmd).Config
			var paths []string
			appdb := filepath.Join(config.DBDir(), "application.db")
			dymintdb := filepath.Join(config.DBDir(), "dymint")
			settlementdb := filepath.Join(config.DBDir(), "settlement")
			paths = append(paths, appdb, dymintdb, settlementdb)
			for _, path := range paths {
				err := removePath(path)
				if err != nil {
					fmt.Printf("error removing %s with error %s\n", path, err)
				}
			}
			return nil
		},
	}
}

func removePath(path string) error {
	if cmtos.FileExists(path) {
		if err := os.RemoveAll(path); err == nil {
			fmt.Println("Path removed", "path", path)
		} else {
			return err
		}
	}
	return nil
}
