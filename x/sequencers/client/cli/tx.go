package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
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
		UnsafeImportKeyCommand(),
	)

	return cmd
}

func NewCreateCmd() *cobra.Command {
	short := "Create a sequencer object, to claim rewards etc."
	long := strings.TrimSpace(short +
		`Requires signature from consensus address public key. Specify consensus key in keyring uid.
Operator addr should be bech32 encoded.`)

	cmd := &cobra.Command{
		Use:   "create-sequencer [keyring uid] [operator addr] [priv key path]",
		Args:  cobra.ExactArgs(3),
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			clientCtx = clientCtx.WithKeyringOptions(func(options *keyring.Options) {
				// options.SupportedAlgos = append(options.SupportedAlgos,)
				_ = options.SupportedAlgos
			})

			acc, err := clientCtx.AccountRetriever.GetAccount(clientCtx, clientCtx.GetFromAddress())
			if err != nil {
				return fmt.Errorf("get account: %w", err)
			}

			var operatorAddr string
			operatorAddr = args[0]
			var keyUID string
			keyUID = args[1]

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags())

			if _, err := txf.Keybase().Key(keyUID); err != nil {
				return fmt.Errorf("check key is available: %w", err)
			}

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
