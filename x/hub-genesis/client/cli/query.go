package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/dymensionxyz/dymension-rdk/x/hub-genesis/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

// GetQueryCmd returns the cli query commands for the hub-genesis module.
func GetQueryCmd() *cobra.Command {
	hubGenQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the hub-genesis module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	hubGenQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryState(),
	)

	return hubGenQueryCmd
}

// GetCmdQueryParams implements a command to return the current hub-genesis
// parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current hub-genesis parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryParamsRequest{}
			res, err := queryClient.Params(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryState implements a command to return the current hub-genesis
// state.
func GetCmdQueryState() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Query the current hub-genesis state",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			state := &types.QueryStateRequest{}
			res, err := queryClient.State(context.Background(), state)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.State)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
