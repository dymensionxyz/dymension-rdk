package cli

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

var (
	GrantPermissionsHandler  = govclient.NewProposalHandler(NewCmdGrantPermissionsProposal)
	RevokePermissionsHandler = govclient.NewProposalHandler(NewCmdRevokePermissionsProposal)
)

// NewCmdGrantPermissionsProposal broadcasts a GrantPermissionsProposal message.
func NewCmdGrantPermissionsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "grant-permissions-proposal address permissions [flags]",
		Short:   "proposal to grant permissions for a specific address",
		Example: `dymd tx gov submit-legacy-proposal grant-permissions-proposal address permission_1,permission_2,...`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(govcli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			address := args[0]
			permissions := strings.Split(args[1], ",")

			addrPermission := []types.AddressPermissions{
				{
					Address: address,
					PermissionList: types.PermissionList{
						Permissions: permissions,
					},
				},
			}

			content := types.NewGrantPermissionsProposal(title, description, addrPermission)
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(govcli.FlagDescription, "", "The proposal description")
	cmd.Flags().String(govcli.FlagDeposit, "", "The proposal deposit")

	return cmd
}

// NewCmdRevokePermissionsProposal broadcasts a RevokePermissionsProposal message.
func NewCmdRevokePermissionsProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "revoke-permissions-proposal address permissions [flags]",
		Short:   "proposal to revoke permissions for a specific address",
		Example: `dymd tx gov submit-legacy-proposal revoke-permissions-proposal address permission_1,permission_2,...`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(govcli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			address := args[0]
			permissions := strings.Split(args[1], ",")

			addrPermission := []types.AddressPermissions{
				{
					Address: address,
					PermissionList: types.PermissionList{
						Permissions: permissions,
					},
				},
			}

			content := types.NewRevokePermissionsProposal(title, description, addrPermission)
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().String(govcli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(govcli.FlagDescription, "", "The proposal description")
	cmd.Flags().String(govcli.FlagDeposit, "", "The proposal deposit")

	return cmd
}
