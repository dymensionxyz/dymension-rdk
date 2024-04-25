package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	denommetadataQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	denommetadataQueryCmd.AddCommand(
		GetCmdIBCDenomBaseOnDenomTrace(),
	)

	return denommetadataQueryCmd
}

// GetCmdIBCDenomBaseOnDenomTrace implements a command to return the IBC denom base on a denom trace.
func GetCmdIBCDenomBaseOnDenomTrace() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ibc-denom [port-id-1]/[channel-id-1]/.../[port-id-n]/[channel-id-n]/[denom]",
		Short: "Get IBC denom base on a denom trace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			denomTrace := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.IBCDenomByDenomTrace(context.Background(), &types.QueryGetIBCDenomByDenomTraceRequest{
				DenomTrace: denomTrace,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
