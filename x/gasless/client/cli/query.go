package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension-rdk/x/gasless/types"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "gasless",
		Short:                      fmt.Sprintf("Querying commands for the %s module", "gasless"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewQueryParamsCmd(),
		NewQueryUsageIdentifiersCmd(),
		NewQueryGasTankCmd(),
		NewQueryGasTanksCmd(),
		NewQueryGasTanksByProviderCmd(),
		NewQueryGasConsumerCmd(),
		NewQueryGasConsumersCmd(),
		NewQueryGasConsumersByGasTankIDCmd(),
		NewQueryUsageIdentifierToTankIdsCmd(),
	)

	return cmd
}

// NewQueryParamsCmd implements the params query command.
func NewQueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current gasless module's parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as gasless module's parameters.
Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			resp, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&resp.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryUsageIdentifiersCmd implements the AvailableUsageIdentifiers query command.
func NewQueryUsageIdentifiersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "usage-identifiers",
		Args:  cobra.NoArgs,
		Short: "Query all the available usage identifiers",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all the available usage identifiers.
Example:
$ %s query %s usage-identifiers
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			resp, err := queryClient.AvailableUsageIdentifiers(
				cmd.Context(),
				&types.QueryAvailableUsageIdentifiersRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryGasTankCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gastank [gas-tank-id]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Query details of the gas tank",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the gas tank
Example:
$ %s query %s gastank 1
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
				return fmt.Errorf("parse gas_tank_id: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GasTank(
				cmd.Context(),
				&types.QueryGasTankRequest{
					GasTankId: gasTankID,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryGasTanksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gastanks ",
		Args:  cobra.NoArgs,
		Short: "Query details of all the gas tanks",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of all the gas tanks
Example:
$ %s query %s gastanks
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GasTanks(
				cmd.Context(),
				&types.QueryGasTanksRequest{
					Pagination: pageReq,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryGasTanksByProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gas-tanks-by-provider [provider]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Query details of all the gas tanks for the given provider",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of all the gas tanks for the given provider
Example:
$ %s query %s gas-tanks-by-provider aib1y755txyzr5n5yy956ydkjttmj8jhwdljawwve8
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sanitizedProvider, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GasTanksByProvider(
				cmd.Context(),
				&types.QueryGasTanksByProviderRequest{
					Provider: sanitizedProvider.String(),
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryGasConsumerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gasconsumer [consumer]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Query details of the gas consumer",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of the gas consumer
Example:
$ %s query %s gasconsumer aib1y755txyzr5n5yy956ydkjttmj8jhwdljawwve8
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sanitizedConsumer, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GasConsumer(
				cmd.Context(),
				&types.QueryGasConsumerRequest{
					Consumer: sanitizedConsumer.String(),
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryGasConsumersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gasconsumers",
		Args:  cobra.NoArgs,
		Short: "Query details of all the gas consumers",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details of all the gas consumers
Example:
$ %s query %s gasconsumers
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GasConsumers(
				cmd.Context(),
				&types.QueryGasConsumersRequest{
					Pagination: pageReq,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryGasConsumersByGasTankIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gas-consumers-by-tank-id [gas-tank-id]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Query all gas consumers for given gas tank id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all gas consumers for given gas tank id
Example:
$ %s query %s gas-consumers-by-tank-id 1
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
				return fmt.Errorf("parse gas_tank_id: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.GasConsumersByGasTankID(
				cmd.Context(),
				&types.QueryGasConsumersByGasTankIDRequest{
					GasTankId: gasTankID,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryUsageIdentifierToTankIdsCmd implements the GasTankIdsForAllUsageIdentifiers query command.
func NewQueryUsageIdentifierToTankIdsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "usage-identifier-tank-ids",
		Args:  cobra.NoArgs,
		Short: "Query all the usage identifiers along with their associated gas tank ids",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all the usage identifiers along with their associated gas tank ids
Example:
$ %s query %s usage-identifier-tank-ids
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			resp, err := queryClient.GasTankIdsForAllUsageIdentifiers(cmd.Context(), &types.QueryGasTankIdsForAllUsageIdentifiersRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
