package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/dymensionxyz/dymension-rdk/server/commands"
	"github.com/dymensionxyz/dymint/conv"
	"github.com/libp2p/go-libp2p"
	"github.com/spf13/cobra"
	tmcmd "github.com/tendermint/tendermint/cmd/cometbft/commands"
	cmtos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/p2p"
)

// add Rollapp commands
func AddRollappCommands(rootCmd *cobra.Command, defaultNodeHome string, appCreator types.AppCreator, appExport types.AppExporter, addStartFlags types.ModuleInitFlags) {
	dymintCmd := &cobra.Command{
		Use:   "dymint",
		Short: "dymint subcommands",
	}

	dymintCmd.AddCommand(
		ShowSequencer(),
		ShowNodeIDCmd(),
		commands.InspectStateCmd(),
		ResetAll(),
		server.VersionCmd(),
		ResetState(),
	)

	rootCmd.AddCommand(
		dymintCmd,
		server.ExportCmd(appExport, defaultNodeHome),
		version.NewVersionCommand(),
		server.NewRollbackCmd(appCreator, defaultNodeHome),
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

func ResetState() *cobra.Command {
	// ResetStateCmd removes the database of the specified CometBFT core instance.
	return &cobra.Command{
		Use:     "reset-state",
		Aliases: []string{"reset_state"},
		Short:   "Remove all the data and WAL",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			config := server.GetServerContextFromCmd(cmd).Config
			appdb := filepath.Join(config.DBDir(), "application.db")
			dymintdb := filepath.Join(config.DBDir(), "dymint")
			settlementdb := filepath.Join(config.DBDir(), "settlement")
			snapshotsdb := filepath.Join(config.DBDir(), "snapshots")
			if cmtos.FileExists(appdb) {
				if err := os.RemoveAll(appdb); err == nil {
					fmt.Println("Removed application.db", "file", appdb)
				} else {
					fmt.Println("error removing application.db", "file", appdb, "err", err)
				}
			}
			if cmtos.FileExists(dymintdb) {
				if err := os.RemoveAll(dymintdb); err == nil {
					fmt.Println("Removed all dymint data", "dir", dymintdb)
				} else {
					fmt.Println("error removing dymint data", "dir", dymintdb, "err", err)
				}
			}
			if cmtos.FileExists(settlementdb) {
				if err := os.RemoveAll(settlementdb); err == nil {
					fmt.Println("Removed all settlement data", "dir", settlementdb)
				} else {
					fmt.Println("error removing settlement data", "dir", settlementdb, "err", err)
				}
			}
			if cmtos.FileExists(snapshotsdb) {
				if err := os.RemoveAll(snapshotsdb); err == nil {
					fmt.Println("Removed all snapshots data", "dir", snapshotsdb)
				} else {
					fmt.Println("error removing snapshots data", "dir", snapshotsdb, "err", err)
				}
			}
			return nil
		},
	}
}
