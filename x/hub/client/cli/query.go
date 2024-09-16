package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/dymensionxyz/dymension-rdk/x/hub/types"
)

// GetQueryCmd returns the cli query commands for the hub module.
func GetQueryCmd() *cobra.Command {
	hubQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the hub module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	hubQueryCmd.AddCommand(
		GetCmdQueryState(),
	)

	return hubQueryCmd
}

// GetCmdQueryState implements a command to return the current hub state.
func GetCmdQueryState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Query the current hub state",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			state := &types.QueryStateRequest{}
			res, err := queryClient.State(cmd.Context(), state)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.State)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
