package cli

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/dymensionxyz/dymension-rdk/x/denommetadata/types"
)

// NewTxCmd returns a root CLI command handler for certain modules transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Denom Metadata subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewCmdCreateDenomMetadata(),
		NewCmdUpdateDenomMetadata(),
	)
	return txCmd
}

// NewCmdCreateDenomMetadata broadcasts a CreateMetadata message.
func NewCmdCreateDenomMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-denom-metadata denommetadata.json [flags]",
		Short: "create new denom metadata for a specific token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := clientCtx.GetFromAddress()

			path := args[0]

			//nolint:gosec
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			metadata := banktypes.Metadata{}
			err = json.Unmarshal([]byte(fileContent), &metadata)
			if err != nil {
				return err
			}

			err = metadata.Validate()
			if err != nil {
				return err
			}

			msg := &types.MsgCreateDenomMetadata{
				SenderAddress: sender.String(),
				TokenMetadata: metadata,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewCmdUpdateDenomMetadata broadcasts a UpdateMetadata message.
func NewCmdUpdateDenomMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-denom-metadata denommetadata.json [flags]",
		Short: "update new denom metadata for a specific token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := clientCtx.GetFromAddress()

			path := args[0]

			//nolint:gosec
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			metadata := banktypes.Metadata{}
			err = json.Unmarshal([]byte(fileContent), &metadata)
			if err != nil {
				return err
			}

			err = metadata.Validate()
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateDenomMetadata{
				SenderAddress: sender.String(),
				TokenMetadata: metadata,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
