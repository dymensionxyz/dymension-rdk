package cli

import (

	// "strings"

	"github.com/dymensionxyz/rollapp/x/sequencers/types"
	"github.com/spf13/cobra"

	stakingcli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group sequencers queries under a subcommand
	// cmd := &cobra.Command{
	// 	Use:                        types.ModuleName,
	// 	Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
	// 	DisableFlagParsing:         true,
	// 	SuggestionsMinimumDistance: 2,
	// 	RunE:                       client.ValidateCmd,
	// }

	// cmd.AddCommand(CmdQueryParams())
	// return cmd
	// this line is used by starport scaffolding # 1

	cmd := stakingcli.GetQueryCmd()
	cmd.Use = types.ModuleName
	return cmd
}
