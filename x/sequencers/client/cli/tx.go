package cli

import (
	"fmt"
	"strings"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
Operator addr should be bech32 encoded. You may supply a reward addr optionally.`)

	cmd := &cobra.Command{
		Use:     "create-sequencer [keyring uid] [operator addr] {reward addr}",
		Example: "create-sequencer fookey ethmvaloper1jkhslh0k3jtdxfjxrtp0z07a06w3uk8w5yyw9u --from foouser --reward-addr",
		Args:    cobra.ExactArgs(2),
		Short:   short,
		Long:    long,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, txf, signingData, err := signingData(cmd, args[0])
			if err != nil {
				return err
			}

			msgs := make([]sdk.Msg, 1)

			msg, err := types.BuildMsgCreateSequencer(signingData, &types.CreateSequencerPayload{OperatorAddr: args[1]})
			if err != nil {
				return fmt.Errorf("build create seq msg: %w", err)
			}

			msgs[0] = msg

			rewardAddr, _ := cmd.Flags().GetString(FlagRewardAddr)
			if rewardAddr != "" {
				msgU, err := types.BuildMsgUpdateSequencer(signingData, &types.UpdateSequencerPayload{RewardAddr: rewardAddr})
				if err != nil {
					return fmt.Errorf("build update seq msg: %w", err)
				}

				msgs = append(msgs, msgU)
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msgs...)
		},
	}

	cmd.Flags().String(FlagRewardAddr, "", "Address to receive rewards for each block proposed.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUpdateCmd() *cobra.Command {
	short := "Create a sequencer object, to claim rewards etc."
	long := strings.TrimSpace(short +
		`Requires signature from consensus address public key. Specify consensus key in keyring uid.
Operator addr should be bech32 encoded.`)

	cmd := &cobra.Command{
		Use:     "update-sequencer [keyring uid] [reward addr]",
		Example: "update-sequencer fookey ethm1lhk5cnfrhgh26w5r6qft36qerg4dclfev9nprc --from foouser",
		Args:    cobra.ExactArgs(2),
		Short:   short,
		Long:    long,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, txf, signingData, err := signingData(cmd, args[0])
			if err != nil {
				return err
			}

			msg, err := types.BuildMsgUpdateSequencer(signingData, &types.UpdateSequencerPayload{RewardAddr: args[1]})
			if err != nil {
				return fmt.Errorf("build update seq msg: %w", err)
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func signingData(cmd *cobra.Command, keyUID string) (client.Context, tx.Factory, types.SigningData, error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return client.Context{}, tx.Factory{}, types.SigningData{}, err
	}

	acc, err := clientCtx.AccountRetriever.GetAccount(clientCtx, clientCtx.GetFromAddress())
	if err != nil {
		return client.Context{}, tx.Factory{}, types.SigningData{}, fmt.Errorf("get account: make sure it has funds: %s: %w", clientCtx.GetFromAddress(), err)
	}

	txf := tx.NewFactoryCLI(clientCtx, cmd.Flags())

	if _, err := txf.Keybase().Key(keyUID); err != nil {
		return client.Context{}, tx.Factory{}, types.SigningData{}, fmt.Errorf("check key is available: key name: %s: %w", keyUID, err)
	}

	return clientCtx, txf, types.SigningData{
		Operator:
		Account: acc,
		ChainID: clientCtx.ChainID,
		Signer: func(msg []byte) ([]byte, cryptotypes.PubKey, error) {
			return txf.Keybase().Sign(keyUID, msg)
		},
	}, nil
}
