package cli

import (
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

// GetTxCmd returns the transaction commands for the module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "sequencer",
		Short:                      fmt.Sprintf("%s transactions subcommands", "sequencer"),
		DisableFlagParsing:         true, // TODO:?
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewCreateCmd(),
	)

	return cmd
}

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-sequencer [operator addr]",
		Args:  cobra.ExactArgs(5),
		Short: "Create a sequencer object, to claim rewards etc.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			acc, err := clientCtx.AccountRetriever.GetAccount(clientCtx, clientCtx.GetFromAddress())
			if err != nil {
				return fmt.Errorf("get account: %w", err)
			}

			var operatorAddr string

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags())

			var keyUID string

			msg, err := types.BuildMsgCreateSequencer(types.SigningData{
				Account: acc,
				ChainID: clientCtx.ChainID,
				Signer: func(msg []byte) ([]byte, cryptotypes.PubKey, error) {
					return txf.Keybase().Sign(keyUID, msg)
				},
			},
				&types.CreateSequencerPayload{OperatorAddr: operatorAddr},
			)

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
