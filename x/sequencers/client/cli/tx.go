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
		Use:     "create-sequencer [keyring uid for cons key] {reward addr}",
		Example: "create-sequencer fooCons --from fooOper --reward-addr ethm1cv7qcksr7cyxv9wgjn3tpxd74n2pffryq7ujw4",
		Args:    cobra.ExactArgs(1),
		Short:   short,
		Long:    long,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.GetClientContextFromCmd(cmd)

			txf, signingData, err := signingData(ctx, cmd, args[0])
			if err != nil {
				return err
			}

			msgs := make([]sdk.Msg, 1)

			msg, err := types.BuildMsgCreateSequencer(signingData, &types.CreateSequencerPayload{OperatorAddr: sdk.ValAddress(ctx.GetFromAddress()).String()})
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

			return tx.GenerateOrBroadcastTxWithFactory(ctx, txf, msgs...)
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
			ctx := client.GetClientContextFromCmd(cmd)

			txf, signingData, err := signingData(ctx, cmd, args[0])
			if err != nil {
				return err
			}

			msg, err := types.BuildMsgUpdateSequencer(signingData, &types.UpdateSequencerPayload{RewardAddr: args[1]})
			if err != nil {
				return fmt.Errorf("build update seq msg: %w", err)
			}

			return tx.GenerateOrBroadcastTxWithFactory(ctx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func signingData(ctx client.Context, cmd *cobra.Command, keyUID string) (tx.Factory, types.SigningData, error) {
	addr := ctx.GetFromAddress()

	acc, err := ctx.AccountRetriever.GetAccount(ctx, addr)
	if err != nil {
		return tx.Factory{}, types.SigningData{}, fmt.Errorf("get account: make sure it has funds: %s: %w", addr, err)
	}

	txf := tx.NewFactoryCLI(ctx, cmd.Flags())

	if _, err := txf.Keybase().Key(keyUID); err != nil {
		return tx.Factory{}, types.SigningData{}, fmt.Errorf("check key is available: key name: %s: %w", keyUID, err)
	}

	return txf, types.SigningData{
		Operator: sdk.ValAddress(addr),
		Account:  acc,
		ChainID:  ctx.ChainID,
		Signer: func(msg []byte) ([]byte, cryptotypes.PubKey, error) {
			return txf.Keybase().Sign(keyUID, msg)
		},
	}, nil
}
