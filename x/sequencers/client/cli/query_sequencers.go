package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/dymensionxyz/dymension-rdk/x/sequencers/types"
)

func CmdQuerySequencers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sequencers",
		Short: "Query for all sequencers",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all sequencers on a network.

Example:
$ %s query %s sequencers
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QuerySequencersRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Sequencers(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQuerySequencer() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "sequencer [sequencer-address]",
		Short: "Query a sequencer",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about an individual sequencer along with its reward address and whitelisted relayers.
Reward address and whitelisted relayers are not returned if not found.

Example:
$ %s query %s sequencer %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, types.ModuleName, bech32PrefixValAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &types.QuerySequencerRequest{
				SequencerAddr: addr.String(),
			}

			res, err := queryClient.Sequencer(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryRewardAddress() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "reward-address [sequencer-address]",
		Short: "Query sequencer reward address",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(fmt.Sprintf(`Query the reward address of a sequencer.

Example:
$ %s query %s reward-address %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`, version.AppName, types.ModuleName, bech32PrefixValAddr)),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryRewardAddressRequest{
				SequencerAddr: addr.String(),
			}

			res, err := queryClient.RewardAddress(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryWhitelistedRelayers() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "whitelisted-relayers [sequencer-address]",
		Short: "Query sequencer whitelisted relayers",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(fmt.Sprintf(`Query whitelisted relayers of a sequencer.

Example:
$ %s query %s whitelisted-relayers %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`, version.AppName, types.ModuleName, bech32PrefixValAddr)),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			addr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := &types.QueryWhitelistedRelayersRequest{
				SequencerAddr: addr.String(),
			}

			res, err := queryClient.WhitelistedRelayers(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
