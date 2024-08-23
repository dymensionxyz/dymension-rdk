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
Operator addr should be bech32 encoded. You may supply a different reward addr optionally.`)

	cmd := &cobra.Command{
		Use:     "create-sequencer [key name] {reward addr}",
		Example: "create-sequencer fooCons --reward-addr ethm1cv7qcksr7cyxv9wgjn3tpxd74n2pffryq7ujw4  --from foo",
		Args:    cobra.ExactArgs(1),
		Short:   short,
		Long:    long,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			keyID := args[0]

			addr := ctx.GetFromAddress()

			txf := tx.NewFactoryCLI(ctx, cmd.Flags())

			if _, err := txf.Keybase().Key(keyID); err != nil {
				return fmt.Errorf("keybase key: %w", err)
			}

			msgs := make([]sdk.Msg, 1)

			msg, err := types.BuildMsgCreateSequencer(func(msg []byte) ([]byte, cryptotypes.PubKey, error) {
				return txf.Keybase().Sign(keyID, msg)
			}, sdk.ValAddress(addr))
			if err != nil {
				return fmt.Errorf("build create seq msg: %w", err)
			}

			msgs[0] = msg

			rewardAddr, _ := cmd.Flags().GetString(FlagRewardAddr)
			if rewardAddr == "" {
				rewardAddr = ctx.GetFromAddress().String()
			}
			msgU := &types.MsgUpdateSequencer{
				Operator:   sdk.ValAddress(ctx.GetFromAddress()).String(),
				RewardAddr: rewardAddr,
			}

			msgs = append(msgs, msgU)

			return tx.GenerateOrBroadcastTxWithFactory(ctx, txf, msgs...)
		},
	}

	cmd.Flags().String(FlagRewardAddr, "", "Address to receive rewards for each block proposed.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUpdateCmd() *cobra.Command {
	short := "Update a sequencer object, to claim rewards etc."
	long := strings.TrimSpace(short +
		`Requires signature from consensus address public key. Specify consensus key in keyring uid.
Operator addr should be bech32 encoded.`)

	cmd := &cobra.Command{
		Use:     "update-sequencer [reward addr]",
		Example: "update-sequencer ethm1lhk5cnfrhgh26w5r6qft36qerg4dclfev9nprc --from foouser",
		Args:    cobra.ExactArgs(1),
		Short:   short,
		Long:    long,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateSequencer{
				Operator:   sdk.ValAddress(ctx.GetFromAddress()).String(),
				RewardAddr: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
