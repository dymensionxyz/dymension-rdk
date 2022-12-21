package cli

import (
	"time"

	"github.com/dymensionxyz/rollapp/x/sequencers/types"
	"github.com/spf13/cobra"

	stakingcli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	// cmd := &cobra.Command{
	// 	Use:                        types.ModuleName,
	// 	Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
	// 	DisableFlagParsing:         true,
	// 	SuggestionsMinimumDistance: 2,
	// 	RunE:                       client.ValidateCmd,
	// }
	// return cmd

	// this line is used by starport scaffolding # 1

	cmd := stakingcli.NewTxCmd()
	cmd.Use = types.ModuleName
	return cmd
}
