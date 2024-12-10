package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	hubgentypes "github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for the module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        hubgentypes.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", hubgentypes.ModuleName),
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewSendTransferCmd(),
	)

	return cmd
}

func NewSendTransferCmd() *cobra.Command {
	short := "Send genesis transfer"
	long := "Send genesis transfer - intended for debugging, since only whitelisted relayer enabled and relayer uses RPC"
	cmd := &cobra.Command{
		Use:   "send-transfer [channel id]",
		Args:  cobra.ExactArgs(1),
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &hubgentypes.MsgSendTransfer{
				Relayer:   ctx.GetFromAddress().String(),
				ChannelId: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
