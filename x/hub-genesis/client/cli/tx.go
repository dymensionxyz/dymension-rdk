package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdGenesisEvent())

	return cmd
}

func CmdGenesisEvent() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "genesis-event [hub-id] [channel-id] [flags]",
		Short:   "Trigger a genesis event from the hub",
		Example: "rollapp-evm tx hub genesis-event [hub-id] [channel-id] ",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get arguments
			hubId := args[0]
			channelId := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgHubGenesisEvent(
				clientCtx.GetFromAddress().String(),
				channelId,
				hubId,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
