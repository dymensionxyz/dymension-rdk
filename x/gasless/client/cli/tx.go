package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

// GetTxCmd returns the transaction commands for the module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "gasless",
		Short:                      fmt.Sprintf("%s transactions subcommands", "gasless"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewCreateGasTankCmd(),
		NewUpdateGasTankStatusCmd(),
		NewUpdateGasTankConfigsCmd(),
		NewBlockConsumerCmd(),
		NewUnblockConsumerCmd(),
		NewUpdateGasConsumerLimitCmd(),
	)

	return cmd
}

func NewCreateGasTankCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-gas-tank [fee-denom] [max-fee-usage-per-tx] [max-fee-usage-per-consumer] [usage-identifiers] [gas-deposit]",
		Args:  cobra.ExactArgs(5),
		Short: "Create a gas tank",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a gas tank.
Example:
$ %s tx %s create-gas-tank stake 25000 5000000 "/cosmos.bank.v1beta1.MsgSend,stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3" 10000000000stake --from mykey
$ %s tx %s create-gas-tank stake 25000 5000000 /cosmos.bank.v1beta1.MsgSend 10000000000stake --from mykey
$ %s tx %s create-gas-tank stake 25000 5000000 stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3 10000000000stake --from mykey
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			feeDenom := args[0]
			if err := sdk.ValidateDenom(feeDenom); err != nil {
				return fmt.Errorf("invalid fee denom: %w", err)
			}

			maxFeeUsagePerTx, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid max-fee-usage-per-tx: %s", args[1])
			}

			maxFeeUsagePerConsumer, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid max-fee-usage-per-consumer: %s", args[2])
			}

			usageIdentifiers, err := ParseStringSliceFromString(args[3], ",")
			if err != nil {
				return err
			}

			gasDeposit, err := sdk.ParseCoinNormalized(args[4])
			if err != nil {
				return fmt.Errorf("invalid gas-deposit: %w", err)
			}

			msg := types.NewMsgCreateGasTank(
				clientCtx.GetFromAddress(),
				feeDenom,
				maxFeeUsagePerTx,
				maxFeeUsagePerConsumer,
				usageIdentifiers,
				gasDeposit,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUpdateGasTankStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-gas-tank-status [gas-tank-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Update status of the gas tank",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update status of the gas tank.
Example:
$ %s tx %s update-gas-tank-status 32 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasTankID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-tank-id: %w", err)
			}

			msg := types.NewMsgUpdateGasTankStatus(
				gasTankID,
				clientCtx.GetFromAddress(),
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUpdateGasTankConfigsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-gas-tank-config [gas-tank-id] [max-fee-usage-per-tx] [max-fee-usage-per-consumer] [usage-identifier]",
		Args:  cobra.ExactArgs(4),
		Short: "Update configs of the gas tank",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update configs of the gas tank.
Example:
$ %s tx %s update-gas-tank-config 1 25000 5000000 /cosmos.bank.v1beta1.MsgSend,stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3 --from mykey
$ %s tx %s update-gas-tank-config 1 25000 5000000 stake1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3 --from mykey
$ %s tx %s update-gas-tank-config 1 25000 5000000 /cosmos.bank.v1beta1.MsgSend --from mykey
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasTankID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-tank-id: %w", err)
			}

			maxFeeUsagePerTx, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid max-fee-usage-per-tx: %s", args[1])
			}

			maxFeeUsagePerConsumer, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid max-fee-usage-per-consumer: %s", args[2])
			}

			usageIdentifier, err := ParseStringSliceFromString(args[3], ",")
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateGasTankConfig(
				gasTankID,
				clientCtx.GetFromAddress(),
				maxFeeUsagePerTx,
				maxFeeUsagePerConsumer,
				usageIdentifier,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewBlockConsumerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block-consumer [gas-tank-id] [consumer]",
		Args:  cobra.ExactArgs(2),
		Short: "Block consumer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Block consumer.
Example:
$ %s tx %s block-consumer 1 stake1.. --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasTankID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-tank-id: %w", err)
			}

			sanitizedConsumer, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgBlockConsumer(
				gasTankID,
				clientCtx.GetFromAddress(),
				sanitizedConsumer,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUnblockConsumerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unblock-consumer [gas-tank-id] [consumer]",
		Args:  cobra.ExactArgs(2),
		Short: "Unblock consumer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unblock consumer.
Example:
$ %s tx %s unblock-consumer 1 stake1.. --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasTankID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-tank-id: %w", err)
			}

			sanitizedConsumer, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnblockConsumer(
				gasTankID,
				clientCtx.GetFromAddress(),
				sanitizedConsumer,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUpdateGasConsumerLimitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-consumer-limit [gas-tank-id] [consumer] [total-fee-consumption-allowed]",
		Args:  cobra.ExactArgs(3),
		Short: "Update consumer consumption limit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update consumer consumption limit.
Example:
$ %s tx %s update-consumer-limit 1 stake1.. 5000000 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasTankID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse gas-tank-id: %w", err)
			}

			sanitizedConsumer, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			totalFeeConsumptionAllowed, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid total-fee-consumption-allowed: %s", args[2])
			}

			msg := types.NewMsgUpdateGasConsumerLimit(
				gasTankID,
				clientCtx.GetFromAddress(),
				sanitizedConsumer,
				totalFeeConsumptionAllowed,
			)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
