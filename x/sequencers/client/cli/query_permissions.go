package cli

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
	"github.com/spf13/cobra"
)

func CmdQueryPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "address-permissions [address]",
		Short: "shows the address's permission",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			address := args[0]
			if _, err := sdk.AccAddressFromBech32(address); err != nil {
				return errors.Wrapf(err, "address format error")
			}

			res, err := queryClient.Permissions(cmd.Context(), &types.QueryPermissionsRequest{
				Address: address,
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
