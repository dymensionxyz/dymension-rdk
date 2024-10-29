package cli

import (
	"fmt"
	"strings"

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
		NewUpdateRewardAddressCmd(),
		NewUpdateWhitelistedRelayersCmd(),
	)

	return cmd
}

func NewUpdateRewardAddressCmd() *cobra.Command {
	short := "Update a sequencer reward address."
	cmd := &cobra.Command{
		Use:     "update-reward-address [addr]",
		Example: "update-reward-address ethm1lhk5cnfrhgh26w5r6qft36qerg4dclfev9nprc --from foouser",
		Args:    cobra.ExactArgs(1),
		Short:   short,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateRewardAddress{
				Operator:   sdk.ValAddress(ctx.GetFromAddress()).String(),
				RewardAddr: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUpdateWhitelistedRelayersCmd() *cobra.Command {
	short := "Update a sequencer whitelisted relayer list."
	cmd := &cobra.Command{
		Use:     "update-whitelisted-relayers [relayers]",
		Example: "update-whitelisted-relayers ethm1lhk5cnfrhgh26w5r6qft36qerg4dclfev9nprc,ethm1lhasdf8969asdfgj2g3j4,ethmasdfkjhgjkhg123j4hgasv7ghi4v --from foouser",
		Args:    cobra.ExactArgs(1),
		Short:   short,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateWhitelistedRelayers{
				Operator: sdk.ValAddress(ctx.GetFromAddress()).String(),
				Relayers: strings.Split(args[0], ","),
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
