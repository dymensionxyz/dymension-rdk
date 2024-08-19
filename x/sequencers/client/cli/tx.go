package cli

import (
	"fmt"
	"strings"

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
		UnsafeImportConsensusKeyCmd(),
		NewCreateCmd(),
		NewUpdateCmd(),
	)

	return cmd
}

func NewCreateCmd() *cobra.Command {
	short := "Create a sequencer object, to claim rewards etc."
	long := strings.TrimSpace(short +
		`Requires signature from consensus address public key. Specify consensus key in keyring uid.
Operator addr should be bech32 encoded.`)

	cmd := &cobra.Command{
		Use:     "create-sequencer [keyring uid] [operator addr]",
		Example: "create-sequencer fookey ethmvaloper1jkhslh0k3jtdxfjxrtp0z07a06w3uk8w5yyw9u --from foouser",
		Args:    cobra.ExactArgs(2),
		Short:   short,
		Long:    long,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			acc, err := clientCtx.AccountRetriever.GetAccount(clientCtx, clientCtx.GetFromAddress())
			if err != nil {
				return fmt.Errorf("get account: make sure it has funds: %s: %w", clientCtx.GetFromAddress(), err)
			}

			var keyUID string
			keyUID = args[0]
			var operatorAddr string
			operatorAddr = args[1]

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags())

			if _, err := txf.Keybase().Key(keyUID); err != nil {
				return fmt.Errorf("check key is available: key name: %s: %w", keyUID, err)
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

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUpdateCmd() *cobra.Command {
	short := "Update a sequencer object, to set a new reward addr etc."
	long := strings.TrimSpace(short +
		`Requires signature from consensus address public key. Specify consensus key in keyring uid.
`)

	cmd := &cobra.Command{
		Use:     "update-sequencer [keyring uid] [reward addr]",
		Example: "update-sequencer fookey ethm1lhk5cnfrhgh26w5r6qft36qerg4dclfev9nprc --from foouser",
		Args:    cobra.ExactArgs(2),
		Short:   short,
		Long:    long,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			acc, err := clientCtx.AccountRetriever.GetAccount(clientCtx, clientCtx.GetFromAddress())
			if err != nil {
				return fmt.Errorf("get account: make sure it has funds: %s: %w", clientCtx.GetFromAddress(), err)
			}

			var keyUID string
			keyUID = args[0]
			var operatorAddr string
			operatorAddr = args[1]

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags())

			if _, err := txf.Keybase().Key(keyUID); err != nil {
				return fmt.Errorf("check key is available: key name: %s: %w", keyUID, err)
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

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
